package logging

import (
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	TargetLocation string `envconfig:"LOG_TARGET_LOCATION" default:"file"`
}

func NewLogger() *log.Logger {
	conf := Config{}
	if err := envconfig.Process("", &conf); err != nil {
		panic("failed to process env for logger: " + err.Error())
	}

	switch conf.TargetLocation {
	default:
		fallthrough
	case "file":
		f, err := os.OpenFile("wofed.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
		if err != nil {
			panic("failed to open file for logger: " + err.Error())
		}
		log.SetOutput(f)

		return log.New(f, "", log.LstdFlags)
	case "stdout":
		return log.Default()
	}
}
