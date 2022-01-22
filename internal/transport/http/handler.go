package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Delaram-Gholampoor-Sagha/production_ready_rest_api/internal/comment"
	"github.com/gorilla/mux"
)

// Handler - stores pointer to our comments service
type Handler struct {
	Router  *mux.Router
	Serivce *comment.Service
}

// NewHandler returns a pointer to a Handler
func NewHandler(service *comment.Service) *Handler {
	return &Handler{
		Serivce: service,
	}
}

// SetupRoutes - sets up all the routes for our application
func (h *Handler) SetupRoutes() {
	fmt.Println("Setting up Routes")
	h.Router = mux.NewRouter()
	h.Router.HandleFunc("/api/comment", h.GetAllComment).Methods("GET")
	h.Router.HandleFunc("/api/comment", h.PostComment).Methods("POST")
	h.Router.HandleFunc("/api/comment/{id}", h.GetComment).Methods("GET")
	h.Router.HandleFunc("/api/comment/{id}", h.UpdateComment).Methods("PUT")
	h.Router.HandleFunc("/api/comment/{id}", h.DeleteComment).Methods("DELETE")
	h.Router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "I am alive!")
	})
}

// GetComment - retrieve a comment by ID
func (h *Handler) GetComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	i, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		fmt.Fprintf(w, "Unable to parse UINT from ID")
	}

	comment, err := h.Serivce.GetComment(uint(i))
	if err != nil {
		fmt.Fprintf(w, "Error Retrieving Comment by ID")
	}
	fmt.Fprintf(w, "%v", comment)

}


// GetAllComment - retrieve all comments from the comment service 
func (h *Handler) GetAllComment(w http.ResponseWriter, r *http.Request) {
   comments , err := h.Serivce.GetAllComment()
   if err != nil {
	   fmt.Fprintf(w , "Failed to retrieve all comments")
   }
	fmt.Fprintf(w, "%v", comments)

}


// PostComment - adds a new comment
func (h *Handler) PostComment(w http.ResponseWriter, r *http.Request) {
	comments , err := h.Serivce.PostComment(comment.Comment{
		Slug: "/",
	})
	if err != nil {
		fmt.Fprintf(w , "Failed to Post a new comment")
	}
	 fmt.Fprintf(w, "%v", comments)
 
 }


 // UpdateComment - update a comment
func (h *Handler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	comments , err := h.Serivce.UpdateComment( 1 , comment.Comment{
		Slug: "/new",
	})
	if err != nil {
		fmt.Fprintf(w , "Failed to updatet comment")
	}
	 fmt.Fprintf(w, "%v", comments)
 
 }


// DeleteComment - delete a comment by ID
func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	commentID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		fmt.Fprintf(w, "Unable to parse UINT from ID")
	}
	
  err = h.Serivce.DeleteComment(uint(commentID))
	if err != nil {
		fmt.Fprintf(w , "Failed to delete comment by comment ID")
	}
	 fmt.Fprintf(w, "successfully deleted comment")
 
 }





