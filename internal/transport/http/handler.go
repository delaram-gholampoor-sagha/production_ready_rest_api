package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"net/http"

	"github.com/Delaram-Gholampoor-Sagha/production_ready_rest_api/internal/comment"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// Handler - stores pointer to our comments service
type Handler struct {
	Router  *mux.Router
	Serivce *comment.Service
}

// Response - an object to store response from our API
type Response struct {
	Message string
	Error   string
}

// NewHandler returns a pointer to a Handler
func NewHandler(service *comment.Service) *Handler {
	return &Handler{
		Serivce: service,
	}
}

func BasicAuth(original func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if user == "admin" && pass == "password" && ok {
			original(w, r)
		} else {
			
			SendErrorResponse(w, "not authorized", errors.New("not authorized"))
			return
		}
	}
}

// JWTAuth - a handy middleware function that will provide basic auth around specific endpoints
func JWTAuth(original func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("jwt auth endpoint hit")
		authHeader := r.Header["Authorization"]
		if authHeader == nil {

			SendErrorResponse(w, "not authorized", errors.New("not authorized"))
			return
		}

		//Bearer jwt-token
		authHeaderParts := strings.Split(authHeader[0], " ")
		if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {

			SendErrorResponse(w, "not authorized", errors.New("not authorized"))
			return
		}

		if validateToken(authHeaderParts[1]) {
			original(w, r)
		} else {

			SendErrorResponse(w, "not authorized", errors.New("not authorized"))
			return
		}
	}
}

// validateToken - validates an incoming jwt token
func validateToken(accessToken string) bool {
	// replace this by loading in a private RSA cert for more security
	var mySigningKey = []byte("missionimpossible")
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		return mySigningKey, nil
	})

	if err != nil {
		return false
	}

	return token.Valid
}

// LoggingMiddleware - a handy middleware function that logs out incoming requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(
			log.Fields{
				"Method": r.Method,
				"Path":   r.URL.Path,
			}).
			Info("handled request")
		next.ServeHTTP(w, r)
	})
}

// SetupRoutes - sets up all the routes for our application
func (h *Handler) SetupRoutes() {
	log.Info("Setting up Routes")
	h.Router = mux.NewRouter()
	h.Router.Use(LoggingMiddleware)
	h.Router.HandleFunc("/api/comment", h.GetAllComment).Methods("GET")
	h.Router.HandleFunc("/api/comment", JWTAuth(h.PostComment)).Methods("POST")
	h.Router.HandleFunc("/api/comment/{id}", h.GetComment).Methods("GET")
	h.Router.HandleFunc("/api/comment/{id}", JWTAuth(h.UpdateComment)).Methods("PUT")
	h.Router.HandleFunc("/api/comment/{id}", JWTAuth(h.DeleteComment)).Methods("DELETE")
	h.Router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {

		if err := SendOkResponse(w, Response{Message: "I am Alive"}); err != nil {
			panic(err)
		}
	})
}

func SendOkResponse(w http.ResponseWriter, resp interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(resp)
}

func SendErrorResponse(w http.ResponseWriter, message string, err error) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusInternalServerError)
	if err := json.NewEncoder(w).Encode(Response{Message: message, Error: err.Error()}); err != nil {
		log.Error(err)
	}
}
