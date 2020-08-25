package logging

import (
	"fmt"
	"go-jwt/config"
	jwtLogrus "go-jwt/logging/logrus"
	"go-jwt/logging/slack"
	"regexp"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

// Handler represented custom logger
type Handler struct {
	cfg         config.General
	logger      Logger
	notifySlack slack.Notifier
}

// Get create new Handler
func Get(cfg config.General, notifyS slack.Notifier) *Handler {
	logrus := jwtLogrus.New()
	return &Handler{cfg: cfg, logger: logrus, notifySlack: notifyS}
}

// PrintLog Notify Slack Alerts
func (dl *Handler) PrintLog(logStr string) {

	err := dl.notifySlack.ProcessForSlackMessage(logStr)
	if err != nil {
		fmt.Errorf("PrintLog, %v", err)
	}
}

// Info log handler
func (dl *Handler) Info(str string) {
	loc, errr := time.LoadLocation(dl.cfg.Timezone)
	if errr != nil {
		err := dl.notifySlack.ProcessForSlackMessage(fmt.Sprintf("Error loading location : %s", errr))
		if err != nil {
			fmt.Errorf("Can not LoadLocation, %v", err)
		}
		return
	}
	ct := time.Now().In(loc).Format(dl.cfg.TimeFullFormat)

	// Skip this function, and fetch the PC and file for its parent
	pc, _, _, _ := runtime.Caller(1)

	// Retrieve a Function object this functions parent
	functionObject := runtime.FuncForPC(pc)

	// Regex to extract just the function name (and not the module path)
	extractFnName := regexp.MustCompile(`^.*\.(.*)$`)
	method := extractFnName.ReplaceAllString(functionObject.Name(), "$1")

	fields := logrus.Fields{
		"method": method,
	}

	dl.logger.Info(fields, str)

	jsonStr := fmt.Sprintf("{time: %s, method: %s, msg: %s}", ct, method, str)
	dl.PrintLog(jsonStr)

}

// Error - log error
func (dl *Handler) Error(msg string, errr error) {
	// Skip this function, and fetch the PC and file for its parent
	pc, _, _, _ := runtime.Caller(1)

	// Retrieve a Function object this functions parent
	functionObject := runtime.FuncForPC(pc)

	// Regex to extract just the function name (and not the module path)
	extractFnName := regexp.MustCompile(`^.*\.(.*)$`)
	method := extractFnName.ReplaceAllString(functionObject.Name(), "$1")

	fields := logrus.Fields{
		"method": method,
	}

	msg = msg + ": " + errr.Error()
	dl.logger.Error(fields, msg)

	jsonStr := fmt.Sprintf("{time: %s, method: %s, msg: %s}", time.Now().Format(dl.cfg.TimeFullFormat), method, msg)
	dl.PrintLog(jsonStr)
}

// DBError - error in DB, send msg to slack
func (dl *Handler) DBError(msg string, errr error) {
	// Skip this function, and fetch the PC and file for its parent
	pc, _, _, _ := runtime.Caller(1)

	// Retrieve a Function object this functions parent
	functionObject := runtime.FuncForPC(pc)

	// Regex to extract just the function name (and not the module path)
	extractFnName := regexp.MustCompile(`^.*\.(.*)$`)
	method := extractFnName.ReplaceAllString(functionObject.Name(), "$1")

	fields := logrus.Fields{
		"method": method,
	}

	msg = msg + ": " + errr.Error()
	dl.logger.Error(fields, msg)

	jsonStr := fmt.Sprintf("{time: %s, method: %s, msg: %s}", time.Now().Format(dl.cfg.TimeFullFormat), method, msg)
	dl.PrintLog(jsonStr)

}

// Fatal - log fatal
func (dl *Handler) Fatal(method string, msg string, errr error) {
	fields := logrus.Fields{
		"method": method,
	}

	msg = msg + ": " + errr.Error()
	dl.logger.Fatal(fields, msg)
}

// Enter - start time of working of method
func (dl *Handler) Enter() (time.Time, string) {
	start := time.Now()

	// Skip this function, and fetch the PC and file for its parent
	pc, _, _, _ := runtime.Caller(1)

	// Retrieve a Function object this functions parent
	functionObject := runtime.FuncForPC(pc)

	// Regex to extract just the function name (and not the module path)
	extractFnName := regexp.MustCompile(`^.*\.(.*)$`)
	method := extractFnName.ReplaceAllString(functionObject.Name(), "$1")

	fields := logrus.Fields{
		"method": method,
	}

	msg := "Started " + method
	dl.logger.Info(fields, msg)

	return start, method
}

// Exit this is to track how much time function took
func (dl *Handler) Exit(start time.Time, name string) {
	elapsed := time.Since(start)

	fields := logrus.Fields{
		"method": name,
	}

	msg := fmt.Sprintf("%s took %s", name, elapsed)
	dl.logger.Info(fields, msg)
}
