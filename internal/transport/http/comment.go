package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Delaram-Gholampoor-Sagha/production_ready_rest_api/internal/comment"
	"github.com/gorilla/mux"
)

// GetComment - retrieve a comment by ID
func (h *Handler) GetComment(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	i, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		SendErrorResponse(w, "Unable to parse UINT from ID", err)
		return
	}

	comment, err := h.Serivce.GetComment(uint(i))
	if err != nil {
		SendErrorResponse(w, "Error Retrieving Comment by ID", err)
		return
	}
	if err = SendOkResponse(w, comment); err != nil {
		panic(err)
	}
}

// GetAllComment - retrieve all comments from the comment service
func (h *Handler) GetAllComment(w http.ResponseWriter, r *http.Request) {

	comments, err := h.Serivce.GetAllComment()
	if err != nil {
		SendErrorResponse(w, "Failed to retrieve all comments", err)
		return
	}
	if err = SendOkResponse(w, comments); err != nil {
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
		SendErrorResponse(w, "Failed to Post a new comment", err)
		return
	}
	if err = SendOkResponse(w, comment); err != nil {
		panic(err)
	}

}

// UpdateComment - update a comment
func (h *Handler) UpdateComment(w http.ResponseWriter, r *http.Request) {

	var cmt comment.Comment
	if err := json.NewDecoder(r.Body).Decode(&cmt); err != nil {
		SendErrorResponse(w, "Failed to decodde JSON body", err)
		return
	}
	vars := mux.Vars(r)
	id := vars["id"]
	commentID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		SendErrorResponse(w, "Unable to parse UINT from ID", err)
		return
	}
	comment, err := h.Serivce.UpdateComment(uint(commentID), cmt)
	if err != nil {
		SendErrorResponse(w, "Failed to updatet comment", err)
		return
	}
	if err = SendOkResponse(w, comment); err != nil {
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
		SendErrorResponse(w, "Unable to parse UINT from ID", err)
		return
	}

	err = h.Serivce.DeleteComment(uint(commentID))
	if err != nil {
		SendErrorResponse(w, "Failed to delete comment by comment ID", err)
		return
	}

	if err = SendOkResponse(w, Response{Message: "Successfully Deleted"}); err != nil {
		panic(err)
	}

}


