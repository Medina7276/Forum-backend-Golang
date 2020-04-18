package service

import (
	"net/http"

	"git.01.alem.school/qjawko/forum/http_errors"
	"git.01.alem.school/qjawko/forum/model"
	"git.01.alem.school/qjawko/forum/repo"
	uuid "github.com/satori/go.uuid"
)

type PostService struct {
	postRepo *repo.PostStore
}

func NewPostService(postRepo *repo.PostStore) *PostService {
	return &PostService{postRepo: postRepo}
}

func (this *PostService) CreatePost(post *model.Post) (*model.Post, error) {
	post.ID = uuid.NewV4()

	if err := this.postRepo.CreatePost(post); err != nil {
		return nil, &http_errors.HttpError{Err: err, Code: http.StatusInternalServerError}
	}
	return post, nil
}

func (this *PostService) GetPostById(id uuid.UUID) (*model.Post, error) {
	post, err := this.postRepo.GetPostById(id)
	if err != nil {
		return nil, &http_errors.HttpError{Err: err, Code: http.StatusNotFound}
	}
	return post, nil
}

func (this *PostService) GetAllPosts() ([]model.Post, error) {
	posts, err := this.postRepo.GetAllPosts()
	if err != nil {
		return nil, &http_errors.HttpError{Err: err, Code: http.StatusInternalServerError}
	}
	return posts, nil
}

func (this *PostService) GetUserPosts(userID uuid.UUID) ([]model.Post, error) {
	posts, err := this.postRepo.GetAllPostsByUserId(userID)
	if err != nil {
		return nil, &http_errors.HttpError{Err: err, Code: http.StatusInternalServerError}
	}
	return posts, nil
}

func (this *PostService) GetAllPostsBySubforumId(subforumID uuid.UUID) ([]model.Post, error) {
	posts, err := this.postRepo.GetAllPostsBySubforumId(subforumID)
	if err != nil {
		return nil, &http_errors.HttpError{Err: err, Code: http.StatusInternalServerError}
	}
	return posts, nil
}

// func (this *PostService) GetPostByTitle(title string) (*model.Post, error) {
// 	post, err := this.postRepo.GetPostByTitle(title)
// 	if err != nil {
// 		return nil, &http_errors.HttpError{Err: err, Code: http.StatusBadRequest}
// 	}
// 	return post, nil
// }

func (this *PostService) UpdatePost(post *model.Post) (*model.Post, error) {
	if err := this.postRepo.UpdatePost(post); err != nil {
		return nil, &http_errors.HttpError{Err: err, Code: http.StatusInternalServerError}
	}
	return this.GetPostById(post.ID)
}

func (this *PostService) DeletePost(id uuid.UUID) error {
	if err := this.postRepo.DeletePost(id); err != nil {
		return &http_errors.HttpError{Err: err, Code: http.StatusInternalServerError}
	}
	return nil
}
