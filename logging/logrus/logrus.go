package logrus

import (
	"github.com/sirupsen/logrus"

	"bytes"
	"os"

	stackdriver "github.com/andyfusniak/stackdriver-gae-logrus-plugin"
)

//LogLogrus - hold logrus
type LogLogrus struct {
	*logrus.Logger
}

type logOutputSplitter struct{}

// Splits log output, error and fatal to stderr and the rest to stdout
func (splitter *logOutputSplitter) Write(p []byte) (n int, err error) {
	if bytes.Contains(p, []byte("level=debug")) || bytes.Contains(p, []byte("level=info")) ||
		bytes.Contains(p, []byte("level=trace")) || bytes.Contains(p, []byte("level=warn")) {
		return os.Stdout.Write(p)
	}
	return os.Stderr.Write(p)
}

// New - return LogLogrus
func New() *LogLogrus {
	logger := logrus.New()

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")

	if projectID != "" {
		// Line of code that need to be added in main func
		formatter := stackdriver.GAEStandardFormatter(
			stackdriver.WithProjectID(os.Getenv("GOOGLE_CLOUD_PROJECT")),
		)
		logger.SetFormatter(formatter)
	} else {
		logger.SetFormatter(&logrus.JSONFormatter{})
	}

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logger.SetOutput(os.Stdout)

	// logger the debug severity or above.
	logger.SetLevel(logrus.DebugLevel)

	return &LogLogrus{logger}
}

// Info handler
func (l *LogLogrus) Info(fields logrus.Fields, msg string) {
	l.Logger.WithFields(fields).Info(msg)
}

// Error handler
func (l *LogLogrus) Error(fields logrus.Fields, msg string) {
	l.Logger.WithFields(fields).Error(msg)
}

// Fatal handler
func (l *LogLogrus) Fatal(fields logrus.Fields, msg string) {
	l.Logger.WithFields(fields).Fatal(msg)
}
