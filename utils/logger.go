package utils

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/penguin-statistics/widget-backend/config"
)

// NewLogger creates a new logrus.Entry with name provided
func NewLogger(name string) *logrus.Entry {
	logger := logrus.New()

	l := logger.WithFields(logrus.Fields{
		"name": name,
	})
	if config.C.DevMode {
		fmt.Println("dev mode true")
		logger.SetLevel(logrus.TraceLevel)
	} else {
		logger.SetFormatter(&logrus.JSONFormatter{})
		logger.SetReportCaller(true)
		logger.SetLevel(logrus.InfoLevel)
	}
	return l
}
