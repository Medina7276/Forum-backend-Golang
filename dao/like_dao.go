package dao

import (
	"database/sql"
	"fmt"

	"git.01.alem.school/qjawko/forum/model"
	uuid "github.com/satori/go.uuid"
)

//LikeStore s
type LikeStore struct {
	*sql.DB
}

//NewLikeStore db
func NewLikeStore(db *sql.DB) *LikeStore {
	return &LikeStore{DB: db}
}

//CreateLike like
func (store *LikeStore) CreateLike(like *model.Like) error {

	_, err := store.Exec(`INSERT INTO like (id, userid, postid, isupvote) VALUES (?, ?, ?, ?)`,
		like.ID, like.UserID, like.PostID, like.IsUpVote)
	return err
}

//GetLikeByID id
func (store *LikeStore) GetLikeByID(id uuid.UUID) (*model.Like, error) {

	like := &model.Like{}

	row := store.QueryRow(`
	SELECT id, userid, postid, isupvote
	FROM like
	WHERE id = ?`, id)

	err := row.Scan(&like.ID, &like.UserID, &like.PostID, &like.IsUpVote)
	return like, err
}

func (store *LikeStore) GetLikesByPostId(id uuid.UUID) ([]model.Like, error) {

	var likes []model.Like

	rows, err := store.Query(`
	SELECT id, userid, postid, isupvote
	FROM like
	WHERE postid = ?`, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var like model.Like

		err := rows.Scan(&like.ID, &like.UserID, &like.PostID, &like.IsUpVote)
		if err != nil {
			return nil, err
		}

		likes = append(likes, like)
	}

	return likes, nil
}

//UpdateLike ul
func (store *LikeStore) UpdateLike(like *model.Like) error {

	res, err := store.Exec(`UPDATE like SET id = ?, userid = ?, postid = ?, isupvote = ?
	WHERE id = ?`, like.ID, like.UserID, like.PostID, like.IsUpVote)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("Like with id %v not found", like.ID)
	}
	return err
}

//DeleteLike delete
func (store *LikeStore) DeleteLike(id uuid.UUID) error {

	rows, err := store.Exec(`DELETE FROM like WHERE id = ?`, id)
	if err != nil {
		return err
	}

	n, err := rows.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf(`Like with id %v not found`, id)
	}
	return err
}
