package routes

import (
	"fmt"
	"go-jwt/database/models"
	"net/http"
	"os"
	"strings"

	"context"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"

	log "github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

// ContextKey is used for context.Context value. The value requires a key that is not primitive type.
type ContextKey string // can be unexported

// ContextKeyRequestID is the ContextKey for RequestID
const ContextKeyRequestID ContextKey = "requestID" // can be unexported

// AttachRequestID will attach a brand new request ID to a http request
func AssignRequestID(ctx context.Context) context.Context {

	reqID := uuid.New()

	return context.WithValue(ctx, ContextKeyRequestID, reqID.String())
}

// GetRequestID will get reqID from a http request and return it as a string
func GetRequestID(ctx context.Context) string {

	reqID := ctx.Value(ContextKeyRequestID)

	if ret, ok := reqID.(string); ok {
		return ret
	}

	return ""
}

func routerLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		ctx = AssignRequestID(ctx)

		r = r.WithContext(ctx)

		log.WithFields(log.Fields{"request_id": GetRequestID(ctx), "method": r.Method, "url": r.RequestURI, "remote_address": r.RemoteAddr}).Info("Start Request")

		next.ServeHTTP(w, r)

		log.WithFields(log.Fields{"request_id": GetRequestID(ctx)}).Info("End Request")
	})
}

func (c *routeServiceHandler) AuthorizeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("authorization")

		if authorizationHeader == "" {
			sendResponse(w, r, StatusUnauthorized, "An authorization header is required", nil)
			return
		}

		bearerToken := strings.Split(authorizationHeader, " ")
		if len(bearerToken) != 2 {
			sendResponse(w, r, StatusUnauthorized, "Invalid authorization token", nil)
			return
		}

		tk := &models.Token{}
		token, error := jwt.ParseWithClaims(bearerToken[1], tk, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error")
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if error != nil {
			sendResponse(w, r, StatusError, error.Error(), nil)
			return
		}

		if token.Valid {
			ctx := context.WithValue(r.Context(), "user", tk)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			sendResponse(w, r, StatusUnauthorized, "Invalid authorization token", nil)
		}
	})
}
