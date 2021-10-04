package main

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/picatz/roku"
)

const ABCAppId = "73376"

type Config struct {
	DeviceName string `envconfig:"DEVICE_NAME" required:"true"`
}

func main() {
	conf := &Config{}
	if err := envconfig.Process("", conf); err != nil {
		log.Fatalf("invalid config: %s", err.Error())
	}

	devices, err := roku.Find(roku.DefaultWaitTime)
	if err != nil {
		log.Fatalf("SSDP request looking for Roku devices failed")
	}

	for _, device := range devices {
		if device != nil {
			info, err := device.DeviceInfo()
			if err == nil {
				if info.FriendlyDeviceName == conf.DeviceName {
					if err := device.LaunchApp(ABCAppId, nil); err != nil {
						panic(err)
					}

					// app takes a long ass time to load
					time.Sleep(10 * time.Second)

					// keys needed to navigate to ABC Live TV
					device.Keypress(roku.RightKey)
					time.Sleep(2 * time.Second)
					device.Keypress(roku.RightKey)
					time.Sleep(2 * time.Second)
					device.Keypress(roku.SelectKey)

					// live TV takes an even longer time to load
					time.Sleep(12 * time.Second)

					// navigate to the "Watch" button and click it
					device.Keypress(roku.DownKey)
					device.Keypress(roku.SelectKey)
				}
			}
		}
	}
}
