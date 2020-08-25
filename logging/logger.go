package logging

import "github.com/sirupsen/logrus"

// Logger - hold logick for logging
type Logger interface {
	Info(fields logrus.Fields, msg string)
	Error(fields logrus.Fields, msg string)
	Fatal(fields logrus.Fields, msg string)
}
