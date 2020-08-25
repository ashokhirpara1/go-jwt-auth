package routes

import (
	"go-jwt/internal"
	"go-jwt/logging"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// InitRoutes - init routs and handlers
func InitRoutes(ctlr internal.Handler, dLog *logging.Handler) {

	// Creates a http server
	router := mux.NewRouter()

	handler := newRouteServiceHandler(ctlr, dLog)

	router.HandleFunc("/", handler.Index).Methods(http.MethodGet)
	router.HandleFunc("/register", handler.Register).Methods(http.MethodPost)
	router.HandleFunc("/login", handler.Login).Methods(http.MethodPost)

	subRouter := router.PathPrefix("/api/v1").Subrouter()
	subRouter.Use(routerLogging)
	subRouter.Use(handler.AuthorizeMiddleware)

	subRouter.HandleFunc("/users", http.HandlerFunc(handler.GetUsers)).Methods(http.MethodGet)

	originsOk := handlers.AllowedOrigins([]string{"*"})
	headersOk := handlers.AllowedHeaders([]string{"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization"})
	methodsOk := handlers.AllowedMethods([]string{http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodOptions})

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: handlers.CORS(originsOk, headersOk, methodsOk)(router)}

	log.Println("Running on port: " + port)

	log.Fatal(server.ListenAndServe())

}
