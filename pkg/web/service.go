package web

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kelseyhightower/envconfig"
	rk "github.com/picatz/roku"
	"github.com/pkg/errors"
	"github.com/scevallos/WheelOfFortuneEveryday/pkg/roku"
)

type Config struct {
	TVAddress string `envconfig:"TV_ADDRESS"`
}

type Service struct {
	Device  *rk.Endpoint
	Client  *roku.Client
	Router  chi.Router
	log     *log.Logger
	healthy *int32
}

// Job is the request input, which describes the job that will be run
// (e.g. starting an app on the TV)
// swagger:model
type Job struct {
	AppName string `json:"appName"`
}

func NewService(client *roku.Client, log *log.Logger, healthy *int32) (*Service, error) {
	conf := Config{}
	if err := envconfig.Process("", &conf); err != nil {
		return nil, err
	}
	// we need a single device to work with
	var device *rk.Endpoint
	if conf.TVAddress != "" {
		// TODO: validate given address
		log.Printf("Using configured address: %s\n", conf.TVAddress)
		device = rk.NewEndpoint(conf.TVAddress)
	} else {
		devices, err := client.Search()
		if err != nil {
			return nil, errors.Wrap(err, "failed roku search")
		}

		num := len(devices)
		switch {
		case num == 0:
			// TODO: add retry bc randomly, no devices will be found for some reason
			return nil, errors.New("no devices found with configured name")
		case num > 1:
			return nil, errors.Errorf("unexpectedly found %d devices with that name\n", num)
		}
		device = devices[0]
	}

	svc := &Service{
		Device:  device,
		Client:  client,
		log:     log,
		healthy: healthy,
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	svc.Routes(r)

	svc.Router = r

	return svc, nil
}

func (s *Service) GetRouter() chi.Router {
	return s.Router
}

func (s *Service) Routes(r chi.Router) {
	r.Use(tracing, logging(s.log))
	r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		if atomic.LoadInt32(s.healthy) == 1 {
			w.Write([]byte(`OK`))
			return
		}
		w.WriteHeader(http.StatusServiceUnavailable)
	})

	// POST /api/v1/jobs CreateJob
	//
	// Creates a new job to start an application on the configured Roku
	//
	// ---
	// consumes:
	// - application/json
	// parameters:
	// - in: body
	//   required: true
	//   name: job
	//   schema:
	//     "$ref": "#/definitions/Job"
	r.Post("/api/v1/jobs", s.tvJobHandler)
}

func (s *Service) tvJobHandler(w http.ResponseWriter, r *http.Request) {
	job := &Job{}
	blob, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeErrorString(w, "failed to read request body", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(blob, job); err != nil {
		writeErrorString(w, "failed to unmarshal request body", http.StatusBadRequest)
		return
	}

	if job.AppName == "" {
		writeErrorString(w, "bad or missing app name: "+job.AppName, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)

	go func() {
		if err := s.Client.StartApp(job.AppName, s.Device); err != nil {
			s.log.Println("failed to start app: " + err.Error())
			return
		}
	}()
}

func writeErrorString(w http.ResponseWriter, err string, code int) {
	w.WriteHeader(code)
	w.Write([]byte(err))
}
