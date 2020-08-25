package routes

import (
	"encoding/json"
	"go-jwt/database/models"
	"go-jwt/internal"
	"go-jwt/logging"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type routeServiceHandler struct {
	ctlr internal.Handler
	Log  *logging.Handler
}

func newRouteServiceHandler(ctlr internal.Handler, dLog *logging.Handler) *routeServiceHandler {
	return &routeServiceHandler{ctlr: ctlr, Log: dLog}
}

// Index - Welcome Route
func (rsh *routeServiceHandler) Index(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, r, StatusSuccess, "Welcome to go-jwt", nil)
	return
}

// Register - Create New User
func (rsh *routeServiceHandler) Register(w http.ResponseWriter, r *http.Request) {

	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		sendResponse(w, r, StatusError, err.Error(), nil)
		return
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		sendResponse(w, r, StatusError, err.Error(), nil)
		return
	}

	user.Password = string(pass)

	createdUser, err := rsh.ctlr.Register(*user)
	if err != nil {
		sendResponse(w, r, StatusError, err.Error(), createdUser)
		return
	}

	sendResponse(w, r, StatusSuccess, "", createdUser)
	return
}

// Login - user login
func (rsh *routeServiceHandler) Login(w http.ResponseWriter, r *http.Request) {

	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		sendResponse(w, r, StatusError, err.Error(), nil)
		return
	}

	resp, err := rsh.ctlr.Login(user.Email, user.Password)
	if err != nil {
		sendResponse(w, r, StatusError, err.Error(), nil)
		return
	}

	sendResponse(w, r, StatusSuccess, "", resp)
	return
}

// GetUsers -
func (rsh *routeServiceHandler) GetUsers(w http.ResponseWriter, r *http.Request) {

	resp, err := rsh.ctlr.GetUsers()
	if err != nil {
		sendResponse(w, r, StatusError, err.Error(), nil)
		return
	}

	sendResponse(w, r, StatusSuccess, "", resp)
	return
}
