package repo

import (
	"database/sql"
	"fmt"

	"git.01.alem.school/qjawko/forum/model"
	uuid "github.com/satori/go.uuid"
)

type CommentStore struct {
	*sql.DB
}

func NewCommentStore(db *sql.DB) *CommentStore {
	return &CommentStore{DB: db}
}

func (store *CommentStore) CreateComment(comment *model.Comment) error {
	_, err := store.Exec(`INSERT INTO comment (id, postid, userid, content, creationdate) 
	VALUES (?, ?, ?, ?, ?)`,
		comment.ID, comment.PostID, comment.UserID, comment.Content,
		comment.CreationDate)
	return err
}

func (store *CommentStore) GetCommentByID(id uuid.UUID) (*model.Comment, error) {
	comment := &model.Comment{}

	row := store.QueryRow(`
	SELECT id, postid, userid, content, creationdate
	FROM comments
	WHERE id = ?`, id)

	err := row.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content,
		&comment.CreationDate)

	return comment, err
}

func (store *CommentStore) GetAllCommentsByPostID(postId uuid.UUID) ([]model.Comment, error) {
	comments := []model.Comment{}

	rows, err := store.Query(`
	SELECT id, postid, userid, content, creationdate
	FROM comments
	WHERE postid = ?`, postId)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var comment model.Comment

		err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreationDate)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

// GetAllCommentsByUserID - когда пользователь будет заходить в свой профиль там будет статистика.
// Какие посты создал и комменты, что лайкнул.
func (store *CommentStore) GetAllCommentsByUserID(userId uuid.UUID) ([]model.Comment, error) {
	comments := []model.Comment{}

	rows, err := store.Query(`
	SELECT id, postid, userid, content, creationdate 
	FROM comments
	WHERE userid = ?`, userId)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var comment model.Comment

		err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreationDate)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

//по контенту, когда пользователь хочет сделать поиск по комментам
// только не =, а через оператор LIKE
func (store *CommentStore) GetAllCommentsByContent(content string) ([]model.Comment, error) {

	comments := []model.Comment{}

	rows, err := store.Query(` 
	SELECT id, postid, userid, content, creationdate 
	FROM comments
	WHERE content Like '%?%'`, content)
	//To query data based on partial information, you use the LIKE
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var comment model.Comment

		err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreationDate)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (store *CommentStore) UpdateComment(comment *model.Comment) error {

	res, err := store.Exec(`UPDATE postid = ?, userid = ?, content = ?, creationdate = ?
	WHERE id = ?`, comment.PostID, comment.UserID, comment.Content,
		comment.CreationDate, comment.ID)

	if err != nil {
		return err
	}

	n, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if n == 0 {
		return fmt.Errorf("Comment with id %v not found", comment.ID)
	}

	return err
}

func (store *CommentStore) DeleteComment(id uuid.UUID) error {
	rows, err := store.Exec(`DELETE FROM comments WHERE Id = ?`, id)
	if err != nil {
		return err
	}

	n, err := rows.RowsAffected()
	if n == 0 {
		return fmt.Errorf("Comment with id %v not found", id)
	}

	return err
}
