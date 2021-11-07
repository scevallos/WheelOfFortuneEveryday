package web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	rk "github.com/picatz/roku"
	"github.com/pkg/errors"
	"github.com/scevallos/WheelOfFortuneEveryday/pkg/roku"
)

type Service struct {
	Device *rk.Endpoint
	Client *roku.Client
	Router chi.Router
}

// Job is the request input, which describes the job that will be run
// (e.g. starting an app on the TV)
// swagger:model
type Job struct {
	AppName string `json:"appName"`
}

func NewService(client *roku.Client) (*Service, error) {
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
	// we need a single device to work with

	svc := &Service{
		Device: devices[0],
		Client: client,
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
	r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`OK`))
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
			fmt.Println("failed to start app: " + err.Error())
			return
		}
	}()
}

func writeErrorString(w http.ResponseWriter, err string, code int) {
	w.WriteHeader(code)
	w.Write([]byte(err))
}
