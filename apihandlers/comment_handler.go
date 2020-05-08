package apihandlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"git.01.alem.school/qjawko/forum/http_errors"
	"git.01.alem.school/qjawko/forum/model"
	"git.01.alem.school/qjawko/forum/service"
	"git.01.alem.school/qjawko/forum/utils"
	uuid "github.com/satori/go.uuid"
)

type CommentHandler struct {
	CommentService *service.CommentService
	Endpoint       string
}

func NewCommentHandler(endpoint string, commentService *service.CommentService) *CommentHandler {
	return &CommentHandler{
		CommentService: commentService,
		Endpoint:       endpoint,
	}
}

func (ch *CommentHandler) Route(w http.ResponseWriter, r *http.Request) {
	var funcToCall func(http.ResponseWriter, *http.Request)

	switch r.Method {
	case http.MethodGet:
		endpoint := r.URL.Path[len(ch.Endpoint):]
		if len(endpoint) == 0 || r.FormValue("user_id") != "" || r.FormValue("search_param") != "" {
			funcToCall = ch.GetAll
		} else {
			funcToCall = ch.GetCommentByID
		}

	case http.MethodPost:
		funcToCall = ch.CreateComment
	case http.MethodPut:
		funcToCall = ch.UpdateComment
	case http.MethodDelete:
		funcToCall = ch.DeleteComment
	default:
		http.Error(w, "Route Not found", http.StatusNotFound)
		return
	}

	funcToCall(w, r)
}

func (ch *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	var comment model.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
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

func (ch *CommentHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	comments, err := ch.CommentService.GetAllComments()
	if err != nil {
		httpErr := err.(*http_errors.HttpError)
		http.Error(w, httpErr.Error(), httpErr.Code)
		return
	}

	userID := r.FormValue("user_id")
	searchParam := r.FormValue("search_param")

	if userID != "" {
		comments = utils.CommentsFilter(comments, func(comment model.Comment) bool {
			return comment.UserID == uuid.FromStringOrNil(userID)
		})
	}

	if searchParam != "" {
		comments = utils.CommentsFilter(comments, func(comment model.Comment) bool {
			return strings.Contains(comment.Content, searchParam)
		})
	}

	json.NewEncoder(w).Encode(comments)
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

func (ch *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len(ch.Endpoint):]

	if err := ch.CommentService.DeleteComment(uuid.FromStringOrNil(id)); err != nil {
		httpErr := err.(*http_errors.HttpError)
		http.Error(w, httpErr.Error(), httpErr.Code)
		return
	}
}

func (ch *CommentHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	var comment model.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "Bad Body", http.StatusBadRequest)
		return
	}

	id := r.URL.Path[len(ch.Endpoint):]
	comment.ID = uuid.FromStringOrNil(id)

	updated, err := ch.CommentService.UpdateComment(&comment)
	if err != nil {
		httpErr := err.(*http_errors.HttpError)
		http.Error(w, httpErr.Error(), httpErr.Code)
		return
	}

	json.NewEncoder(w).Encode(updated)
}
