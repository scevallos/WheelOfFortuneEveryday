package roku

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/picatz/roku"
	"github.com/pkg/errors"
)

var (
	// mapping of app name to app IDs (needed to launch the app)
	appIds = map[string]string{
		"ABC":     roku.ABC,
		"NBC":     "68669",
		"Live TV": roku.LiveTV,
	}

	// instructions maps the apps to the set of steps that need to
	// occur in order to reach the Live section of that app
	instructions = map[string]func(device *roku.Endpoint){
		"ABC": func(device *roku.Endpoint) {
			// app takes a long ass time to load
			time.Sleep(15 * time.Second)

			// keys needed to navigate to ABC Live TV
			device.Keypress(roku.RightKey)
			time.Sleep(2 * time.Second)
			device.Keypress(roku.RightKey)
			time.Sleep(2 * time.Second)
			device.Keypress(roku.SelectKey)

			// live TV takes an even longer time to load
			time.Sleep(16 * time.Second)

			// navigate to the "Watch" button and click it
			// otherwise the overlay won't go away
			device.Keypress(roku.DownKey)
			device.Keypress(roku.SelectKey)
		},
		"NBC": func(device *roku.Endpoint) {
			// wait for app to load
			time.Sleep(16 * time.Second)

			device.Keypress(roku.UpKey)
			time.Sleep(2 * time.Second)
			device.Keypress(roku.RightKey)
			time.Sleep(2 * time.Second)
			device.Keypress(roku.RightKey)
			time.Sleep(2 * time.Second)
			device.Keypress(roku.SelectKey)
		},
		"Live TV": func(device *roku.Endpoint) {
		},
	}

	defaultOptions = &ClientOptions{
		Logger: log.Default(),
		Config: &Config{},
	}
)

type Config struct {
	// Name of device to search for
	DeviceName string `envconfig:"DEVICE_NAME"`

	// Time to wait during search for devices
	WaitTime time.Duration `envconfig:"WAIT_TIME" default:"5s"`
}

type StartAppOptions struct {
	AppName string
	Device  *roku.Endpoint
	Force   bool
}

type Client struct {
	*Config
	log *log.Logger
}

type ClientOptions struct {
	Logger *log.Logger
	Config *Config
}

func NewClient(options *ClientOptions) (*Client, error) {
	if options == nil {
		options = defaultOptions
	}

	conf := options.Config
	if conf == nil {
		if err := envconfig.Process("", conf); err != nil {
			return nil, errors.Wrap(err, "failed to configure roku client")
		}
	}

	return &Client{
		Config: conf,
		log:    options.Logger,
	}, nil
}

// Search returns a list of devices with a name that match the configured name
func (c *Client) Search() (roku.Endpoints, error) {
	c.log.Println("Searching for nearby devices...")
	devices, err := roku.Find(int(c.WaitTime.Seconds()))
	if err != nil {
		return nil, errors.Wrap(err, "SSDP request to find Roku devices failed")
	}

	matches := roku.Endpoints{}

	for _, device := range devices {
		if device == nil {
			continue
		}

		info, err := device.DeviceInfo()
		if err != nil {
			// log.Debug("could not get info for a device, skipping")
			continue
		}

		if info.FriendlyDeviceName == c.DeviceName {
			matches = append(matches, device)
			// apps, err := device.Apps()
			// if err == nil {
			// 	for _, app := range apps {
			// 		fmt.Println(app.ID, app.Name)
			// 	}
			// }
		}
	}

	return matches, nil
}

func (c *Client) StartApp(opts *StartAppOptions) error {
	c.log.Printf("Starting %s...\n", opts.AppName)
	id, ok := appIds[opts.AppName]
	if !ok {
		return errors.New("app '" + opts.AppName + "' does not have registered id mapping")
	}

	steps, ok := instructions[opts.AppName]
	if !ok {
		return errors.New("app '" + opts.AppName + "' does not have registered instructions")
	}

	info, err := opts.Device.DeviceInfo()
	if err != nil {
		return errors.Wrap(err, "failed to get device info")
	}

	// turn on the TV if it's not already on
	if info.PowerMode != "PowerOn" {
		c.log.Println("Turning on TV...")
		opts.Device.Keypress(roku.PowerOffKey)
		time.Sleep(5 * time.Second)
	} else if !opts.Force {
		// unless we are force starting the app, do nothing if the TV is already on
		c.log.Println("Doing nothing since TV is on already")
		return nil
	}

	if err := opts.Device.LaunchApp(id, nil); err != nil {
		return errors.Wrapf(err, "failed to launch %s app", opts.AppName)
	}

	if steps != nil {
		go steps(opts.Device)
	}

	return nil
}
