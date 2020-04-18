package service

import (
	"net/http"

	"git.01.alem.school/qjawko/forum/http_errors"
	"git.01.alem.school/qjawko/forum/model"
	"git.01.alem.school/qjawko/forum/repo"
	uuid "github.com/satori/go.uuid"
)

type CommentService struct {
	commentRepo *repo.CommentStore
}

func NewCommentService(commentRepo *repo.CommentStore) *CommentService {
	return &CommentService{commentRepo: commentRepo}
}

func (this *CommentService) CreateComment(comment *model.Comment) (*model.Comment, error) {
	comment.ID = uuid.NewV4()

	if err := this.commentRepo.CreateComment(comment); err != nil {
		return nil, &http_errors.HttpError{Err: err, Code: http.StatusInternalServerError}
	}
	return comment, nil
}

func (this *CommentService) GetCommentByID(Id uuid.UUID) (*model.Comment, error) {
	comment, err := this.commentRepo.GetCommentByID(Id)
	if err != nil {
		return nil, &http_errors.HttpError{Err: err, Code: http.StatusNotFound}
	}
	return comment, nil
}

func (this *CommentService) GetAllCommentsByPostID(postId uuid.UUID) ([]model.Comment, error) {
	comments, err := this.commentRepo.GetAllCommentsByUserID(postId)
	if err != nil {
		return nil, &http_errors.HttpError{Err: err, Code: http.StatusInternalServerError}
	}
	return comments, nil
}

func (this *CommentService) GetAllCommentsByUserID(userId uuid.UUID) ([]model.Comment, error) {
	comments, err := this.commentRepo.GetAllCommentsByUserID(userId)
	if err != nil {
		return nil, &http_errors.HttpError{Err: err, Code: http.StatusInternalServerError}
	}
	return comments, nil
}

//по контенту, когда пользователь хочет сделать поиск по комментам
// только не =, а через оператор LIKE
func (this *CommentService) GetAllCommentsByContent(content string) ([]model.Comment, error) {
	comments, err := this.commentRepo.GetAllCommentsByContent(content)
	if err != nil {
		return nil, &http_errors.HttpError{Err: err, Code: http.StatusInternalServerError}
	}
	return comments, nil
}

func (this *CommentService) UpdateComment(comment *model.Comment) (*model.Comment, error) {
	if err := this.commentRepo.UpdateComment(comment); err != nil {
		return nil, &http_errors.HttpError{Err: err, Code: http.StatusInternalServerError}
	}
	return this.GetCommentByID(comment.ID)
}

func (this *CommentService) DeleteComment(id uuid.UUID) error {
	if err := this.commentRepo.DeleteComment(id); err != nil {
		return &http_errors.HttpError{Err: err, Code: http.StatusInternalServerError}
	}
	return nil
}
