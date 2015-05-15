package logger

import (
	"github.com/Sirupsen/logrus"
	"github.com/jamesclonk-io/stdlib/env"
)

func init() {
	if env.Get("JCIO_ENV", "") == "production" || // manual
		env.Get("VCAP_APPLICATION", "") != "" || // cf / lattice
		env.Get("DYNO", "") != "" { // heroku
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
}

func GetLogger() *logrus.Logger {
	return logrus.StandardLogger()
}
