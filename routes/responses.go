package routes

import (
	"encoding/json"
	"io"
	"net/http"
)

// StatusSuccess
const (
	StatusSuccess      = "success"
	StatusFail         = "fail"
	StatusUnauthorized = "unauthorized"
	StatusForbidden    = "forbidden"
	StatusError        = "error"
)

// Response data format for HTTP
type Response struct {
	Status  string      `json:"status" bson:"status"`                       // Status code (error|fail|success)
	Code    int         `json:"code"  bson:"code"`                          // HTTP status code
	Message string      `json:"message,omitempty" bson:"message,omitempty"` // Error or status message
	Data    interface{} `json:"data,omitempty" bson:"data,omitempty"`       // Data payload
}

// jsonHTTPEncode is a wrapper for json.NewEncoder(w).Encode(v)
func jsonHTTPEncode(w io.Writer, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}

// convert the HTTP status code into JSend status
func getStatusCode(status string) int {
	if status == StatusFail {
		return http.StatusBadRequest
	}

	if status == StatusUnauthorized {
		return http.StatusUnauthorized
	}

	if status == StatusForbidden {
		return http.StatusForbidden
	}

	return http.StatusOK
}

// setHeaders set the default headers
func setHeaders(hw http.ResponseWriter, contentType string, code int) {
	hw.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	hw.Header().Set("Pragma", "no-cache")
	hw.Header().Set("Expires", "0")
	hw.Header().Set("Content-Type", contentType)
	hw.WriteHeader(code)
}

// sendResponse sends the HTTP response in JSON format
func sendResponse(hw http.ResponseWriter, hr *http.Request, status string, message string, data interface{}) {

	code := getStatusCode(status)

	if status != StatusError && status != StatusSuccess {
		status = StatusFail
	}

	response := Response{
		Status:  status,
		Code:    code,
		Message: message,
		Data:    data,
	}

	// send JSON response
	setHeaders(hw, "application/json", code)
	err := jsonHTTPEncode(hw, response)
	if err != nil {
		//log.Error("Unable to send JSON response", err)
	}
}
