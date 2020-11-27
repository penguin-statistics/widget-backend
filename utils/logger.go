package utils

import "github.com/sirupsen/logrus"

func NewLogger(name string) *logrus.Entry {
	logger := logrus.New()
	logger.SetLevel(logrus.TraceLevel)

	l := logger.WithFields(logrus.Fields{
		"name": name,
	})
	return l
}
