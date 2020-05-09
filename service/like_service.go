package service

import (
	"net/http"

	"git.01.alem.school/qjawko/forum/dao"
	"git.01.alem.school/qjawko/forum/http_errors"
	"git.01.alem.school/qjawko/forum/model"
	uuid "github.com/satori/go.uuid"
)

type LikeService struct {
	likeDao *dao.LikeStore
}

func NewLikeService(likeDao *dao.LikeStore) *LikeService {
	return &LikeService{likeDao: likeDao}
}

func (this *LikeService) CreateLike(like *model.Like) (*model.Like, error) {

	like.ID = uuid.NewV4()

	if err := this.likeDao.CreateLike(like); err != nil {
		return nil, &http_errors.HttpError{Err: err, Code: http.StatusInternalServerError}
	}
	return like, nil
}

func (this *LikeService) GetLikeByID(id uuid.UUID) (*model.Like, error) {

	like, err := this.likeDao.GetLikeByID(id)
	if err != nil {
		return nil, &http_errors.HttpError{Err: err, Code: http.StatusNotFound}
	}
	return like, nil
}

func (this *LikeService) UpdateLike(like *model.Like) (*model.Like, error) {

	if err := this.likeDao.UpdateLike(like); err != nil {
		return nil, &http_errors.HttpError{Err: err, Code: http.StatusInternalServerError}
	}
	return like, nil
}

func (this *LikeService) DeleteLike(id uuid.UUID) error {

	if err := this.likeDao.DeleteLike(id); err != nil {
		return &http_errors.HttpError{Err: err, Code: http.StatusInternalServerError}
	}

}
