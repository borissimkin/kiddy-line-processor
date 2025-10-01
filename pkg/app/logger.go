package app

import (
	"os"

	"github.com/sirupsen/logrus"
)

// SetLogger initialize logger.
func SetLogger(level string) {
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logLevel)
	}

	logrus.SetFormatter(&logrus.JSONFormatter{ //nolint:exhaustruct
		TimestampFormat: "2006-01-02 15:04:05",
	})

	logrus.SetOutput(os.Stdout)
}
