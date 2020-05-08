package service

import (
	"net/http"

	"git.01.alem.school/qjawko/forum/dao"
	"git.01.alem.school/qjawko/forum/http_errors"
	"git.01.alem.school/qjawko/forum/model"
	uuid "github.com/satori/go.uuid"
)

type PostService struct {
	postDao *dao.PostStore
}

func NewPostService(postDao *dao.PostStore) *PostService {
	return &PostService{postDao: postDao}
}

func (this *PostService) CreatePost(post *model.Post) (*model.Post, error) {
	post.ID = uuid.NewV4()

	if err := this.postDao.CreatePost(post); err != nil {
		return nil, &http_errors.HttpError{Err: err, Code: http.StatusInternalServerError}
	}
	return post, nil
}

func (this *PostService) GetPostById(id uuid.UUID) (*model.Post, error) {
	post, err := this.postDao.GetPostById(id)
	if err != nil {
		return nil, &http_errors.HttpError{Err: err, Code: http.StatusNotFound}
	}
	return post, nil
}

func (this *PostService) GetAllPosts() ([]model.Post, error) {
	posts, err := this.postDao.GetAllPosts()
	if err != nil {
		return nil, &http_errors.HttpError{Err: err, Code: http.StatusInternalServerError}
	}
	return posts, nil
}

func (this *PostService) GetUserPosts(userID uuid.UUID) ([]model.Post, error) {
	posts, err := this.postDao.GetAllPostsByUserId(userID)
	if err != nil {
		return nil, &http_errors.HttpError{Err: err, Code: http.StatusInternalServerError}
	}
	return posts, nil
}

func (this *PostService) GetAllPostsBySubforumId(subforumID uuid.UUID) ([]model.Post, error) {
	posts, err := this.postDao.GetAllPostsBySubforumId(subforumID)
	if err != nil {
		return nil, &http_errors.HttpError{Err: err, Code: http.StatusInternalServerError}
	}
	return posts, nil
}

func (this *PostService) UpdatePost(post *model.Post) (*model.Post, error) {
	if err := this.postDao.UpdatePost(post); err != nil {
		return nil, &http_errors.HttpError{Err: err, Code: http.StatusInternalServerError}
	}
	return this.GetPostById(post.ID)
}

func (this *PostService) DeletePost(id uuid.UUID) error {
	if err := this.postDao.DeletePost(id); err != nil {
		return &http_errors.HttpError{Err: err, Code: http.StatusInternalServerError}
	}
	return nil
}
