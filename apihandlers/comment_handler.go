package apihandlers

import (
	"encoding/json"
	"net/http"

	"git.01.alem.school/qjawko/forum/dto"
	"git.01.alem.school/qjawko/forum/http_errors"
	"git.01.alem.school/qjawko/forum/model"
	"git.01.alem.school/qjawko/forum/service"
	"git.01.alem.school/qjawko/forum/util"
	uuid "github.com/satori/go.uuid"
)

type CommentHandler struct {
	CommentService *service.CommentService
	Endpoint       string
}

func NewCommentHandler(endpoint string, CommentService *service.CommentService) *CommentHandler {
	return &CommentHandler{
		CommentService: commentService,
		Endpoint:       endpoint,
	}
}

func (ch *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	var comment model.Comment
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	created, err := ch.CommentService.CreateComment(&comment)
	if err != nil {
		httpErr := err.(*http_errors.HttpError)
		http.Error(w, httpErr.Error(), httpErr.Code)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (ch *CommentHandler) GetCommentByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len(ch.Endpoint):]

	comment, err := ch.CommentService.GetCommentByID(uuid.FromStringOrNil(id))
	if err != nil {
		httpErr := err.(*http_errors.HttpError)
		http.Error(w, httpErr.Error(), httpErr.Code)
		return
	}
	json.NewEncoder(w).Encode(comment)
}

func (ch *CommentHandler) GetAllCommentsByUserID(w http.ResponseWriter, r *http.Request) {

	user, err := util.GetUser(r)
	if err != nil {
		http.Error(w, "User Not Authorized", http.StatusBadRequest)
		return
	}

	comments, err := ch.CommentService.GetAllCommentsByUserID(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(comments)
}

func (ch *CommentHandler) GetAllCommentsByContent(w http.ResponseWriter, r *http.Request) {

	var searchParams dto.CommentSeacrhParamsDTO

	json.NewDecoder(r.Body).Decode(&searchParams)

	comments, err := ch.CommentService.GetAllCommentsByContent(searchParams.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(comments)
}

func (ch *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len(ch.Endpoint):]

	if err := ch.CommentService.DeleteComment(uuid.FromStringOrNil(id)); err != nil {
		httpErr := err.(*http_errors.HttpError)
		http.Error(w, httpErr.Error(), httpErr.Code)
		return
	}

}

func (ch *CommentHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {

}
