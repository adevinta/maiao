package gh

import (
	"github.com/sirupsen/logrus"
)

// Logger provides the default logger for the whole maiao project
var Logger = logrus.New()

type nopWriter struct{}

func (w nopWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func init() {
	Logger.SetOutput(nopWriter{})
}
