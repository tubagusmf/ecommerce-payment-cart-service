package config

import "github.com/sirupsen/logrus"

func SetupLogger() {
	// TODO: setup logger
	log := logrus.New()

	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
}
