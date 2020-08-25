package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go-jwt/config"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Slack - holds configurations data
type Slack struct {
	messaging config.Messaging
	general   config.General
}

// Notifier Hold logick for notify slack
type Notifier interface {
	ProcessForSlackMessage(logStr string) error
}

// Input - holds data for slack
type Input struct {
	Channel   string `json:"channel"`
	Text      string `json:"text"`
	LinkNames bool   `json:"link_names"`
	Username  string `json:"username"`
}

// Get - create new Notify
func Get(messaging config.Messaging, general config.General) Notifier {
	return &Slack{messaging, general}
}

// ProcessForSlackMessage (logStr contains string with log time, method name and log message)
func (s *Slack) ProcessForSlackMessage(logStr string) error {

	// check for error on string
	if !strings.Contains(strings.ToLower(logStr), "error") {
		return nil
	}

	//send slack notification
	notificationText := logStr

	inputStr := Input{Text: "@here " + notificationText, Channel: s.messaging.SlackChannelID, LinkNames: s.messaging.SlackLinkNames}

	err := s.sendSlackMessage(inputStr)
	if err != nil {
		return err
	}

	return nil
}

func (s *Slack) sendSlackMessage(input interface{}) error {

	// Create a new context and add a timeout to it
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(s.general.HTTPRequestTimeout))
	defer cancel() // The cancel should be deferred so resources are cleaned up

	// Convert structs to JSON.
	jsReq, _ := json.Marshal(input)
	plReq := bytes.NewBuffer(jsReq)
	reqReq, _ := http.NewRequest(http.MethodPost, s.messaging.SlackBaseURL+"/api/chat.postMessage", plReq)

	reqReq.Header.Add("authorization", s.messaging.SlackAuthToken)
	reqReq.Header.Add("content-type", "application/json")

	// Associate the cancellable context we just created to the request
	reqReq = reqReq.WithContext(ctx)

	reqRes, err := http.DefaultClient.Do(reqReq)

	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("HttpRequestTimeout error - %v", ctx.Err())
	}

	if err != nil {
		return fmt.Errorf("POST error - %v", err)
	}

	defer reqRes.Body.Close()

	body, err := ioutil.ReadAll(reqRes.Body)
	if err != nil {
		fmt.Println(fmt.Sprintf("ReadAll error - %v - body : %s", err, string(body)))
	}
	return nil
}
