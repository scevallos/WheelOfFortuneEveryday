package main

import (
	"fmt"
	"log"

	"github.com/kelseyhightower/envconfig"
	"github.com/picatz/roku"
)

type Config struct {
	DeviceName string `envconfig:"DEVICE_NAME" required:"true"`
}

func main() {
	fmt.Println("test")

	conf := &Config{}
	if err := envconfig.Process("", conf); err != nil {
		log.Fatalf("invalid config: %s", err.Error())
	}


	devices, err := roku.Find(roku.DefaultWaitTime)
	if err != nil {
		panic(err)
	}

	fmt.Println("len(devices)", len(devices))

	for _, device := range devices {
		if device != nil {
			info, err := device.DeviceInfo()
			if err == nil {
				if info.FriendlyDeviceName == conf.DeviceName {
					// launch ABC app
					// keys needed to navigate to ABC Live
				}
			}
		}
	}
}
