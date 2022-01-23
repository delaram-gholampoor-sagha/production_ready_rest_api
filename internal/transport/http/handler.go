package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Delaram-Gholampoor-Sagha/production_ready_rest_api/internal/comment"
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
	Error string
}

// NewHandler returns a pointer to a Handler
func NewHandler(service *comment.Service) *Handler {
	return &Handler{
		Serivce: service,
	}
}

// SetupRoutes - sets up all the routes for our application
func (h *Handler) SetupRoutes() {
	log.Info("Setting up Routes")
	h.Router = mux.NewRouter()
	h.Router.HandleFunc("/api/comment", h.GetAllComment).Methods("GET")
	h.Router.HandleFunc("/api/comment", h.PostComment).Methods("POST")
	h.Router.HandleFunc("/api/comment/{id}", h.GetComment).Methods("GET")
	h.Router.HandleFunc("/api/comment/{id}", h.UpdateComment).Methods("PUT")
	h.Router.HandleFunc("/api/comment/{id}", h.DeleteComment).Methods("DELETE")
	h.Router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
	
		if err := sendOkResponse(w ,Response{Message: "I am Alive"}); err != nil {
			panic(err)
		}
	})
}

// GetComment - retrieve a comment by ID
func (h *Handler) GetComment(w http.ResponseWriter, r *http.Request) {
	
	vars := mux.Vars(r)
	id := vars["id"]
	i, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		sendErrorResponse(w, "Unable to parse UINT from ID" , err)
		return
	}

	comment, err := h.Serivce.GetComment(uint(i))
	if err != nil {
		sendErrorResponse(w, "Error Retrieving Comment by ID" , err)
		return
	}
	if err = sendOkResponse(w, comment); err != nil {
		panic(err)
	}
}

// GetAllComment - retrieve all comments from the comment service
func (h *Handler) GetAllComment(w http.ResponseWriter, r *http.Request) {
	
	comments, err := h.Serivce.GetAllComment()
	if err != nil {
		sendErrorResponse(w, "Failed to retrieve all comments"  ,err)
		return
	}
	if err = sendOkResponse(w, comments); err != nil {
		panic(err)
	}

}

// PostComment - adds a new comment
func (h *Handler) PostComment(w http.ResponseWriter, r *http.Request) {

	var cmt comment.Comment
	if err := json.NewDecoder(r.Body).Decode(&cmt); err != nil {
		fmt.Fprintf(w, "Failed to decodde JSON body")
	}
	comment, err := h.Serivce.PostComment(cmt)
	if err != nil {
		sendErrorResponse(w, "Failed to Post a new comment" , err)
		return
	}
	if err = sendOkResponse(w, comment); err != nil {
		panic(err)
	}

}

// UpdateComment - update a comment
func (h *Handler) UpdateComment(w http.ResponseWriter, r *http.Request) {

	var cmt comment.Comment
	if err := json.NewDecoder(r.Body).Decode(&cmt); err != nil {
		sendErrorResponse(w, "Failed to decodde JSON body" , err)
		return
	}
	vars := mux.Vars(r)
	id := vars["id"]
	commentID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		sendErrorResponse(w, "Unable to parse UINT from ID" , err)
		return
	}
	comment, err := h.Serivce.UpdateComment(uint(commentID), cmt)
	if err != nil {
		sendErrorResponse(w, "Failed to updatet comment" , err)
		return
	}
	if err = sendOkResponse(w, comment); err != nil {
		panic(err)
	}

}

// DeleteComment - delete a comment by ID
func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	vars := mux.Vars(r)
	id := vars["id"]
	commentID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		sendErrorResponse(w, "Unable to parse UINT from ID" , err)
		return
	}

	err = h.Serivce.DeleteComment(uint(commentID))
	if err != nil {
		sendErrorResponse(w, "Failed to delete comment by comment ID" , err )
		return
	}
	
	if err = sendOkResponse(w, Response{Message: "Successfully Deleted"}); err != nil {
		panic(err)
	}

}

func sendOkResponse(w http.ResponseWriter, resp interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(resp)
}


func sendErrorResponse(w http.ResponseWriter, message string, err error) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusInternalServerError)
	if err := json.NewEncoder(w).Encode(Response{Message: message, Error: err.Error()}); err != nil {
		panic(err)
	}
}
