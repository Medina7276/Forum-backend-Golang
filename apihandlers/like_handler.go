package apihandlers

import (
	"encoding/json"
	"net/http"

	"git.01.alem.school/qjawko/forum/http_errors"
	"git.01.alem.school/qjawko/forum/model"
	"git.01.alem.school/qjawko/forum/service"
	uuid "github.com/satori/go.uuid"
)

type LikeHandler struct {
	LikeService *service.LikeService
	Endpoint    string
}

func NewLikeHandler(endpoint string, likeService *service.LikeService) {
	return &LikeHandler{
		LikeService: likeService,
		Endpoint:    endpoint,
	}
}

func (lh *LikeHandler) CreateLike(w http.ResponseWriter, r *http.Request) {
	var like model.Like
	if err := json.NewDecoder(r.Body).Decode(&like); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	created, err := lh.LikeService.CreateLike(&like)
	if err != nil {
		httpErr := err.(*http_errors.HttpError)
		http.Error(w, httpErr.Error(), httpErr.Code)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (lh *LikeHandler) GetLikeByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len(lh.Endpoint):]
	like, err := lh.LikeService.GetLikeByID(uuid.FromStringOrNil(id))
	if err != nil {
		httpErr := err.(*http_errors.HttpError)
		http.Error(w, httpErr.Error(), httpErr.Code)
		return
	}

	json.NewEncoder(w).Encode(like)
}

func (lh *LikeHandler) Update(w http.ResponseWriter, r *http.Request) {
	var like model.Like
	if err := json.NewDecoder(r.Body).Decode(&like); err != nil {
		http.Error(w, "Bad Body", http.StatusBadRequest)
		return
	}

	idFromURL := r.URL.Path[len(lh.Endpoint):]
	like.ID = uuid.FromStringOrNil(idFromURL)

	updated, err := lh.LikeService.UpdateLike(&like)
	if err != nil {
		httpErr := err.(*http_errors.HttpError)
		http.Error(w, httpErr.Error(), httpErr.Code)
		return
	}

	json.NewEncoder(w).Encode(updated)
}

//Delete delete
func (lh *LikeHandler) DeleteLike(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Path[len(lh.Endpoint):]
	if err := lh.LikeService.DeleteLike(uuid.FromStringOrNil(id)); err != nil {
		httpErr := err.(*http_errors.HttpError)
		http.Error(w, httpErr.Error(), httpErr.Code)
		return
	}
}
