package dao

import (
	"database/sql"
	"fmt"

	"git.01.alem.school/qjawko/forum/model"
	uuid "github.com/satori/go.uuid"
)

type PostStore struct {
	*sql.DB
}

func NewPostStore(db *sql.DB) *PostStore {
	return &PostStore{DB: db}
}

func (store *PostStore) CreatePost(post *model.Post) error {
	_, err := store.Exec(`INSERT INTO posts (id, title, content, creationdate, subforumid, userid) 
	VALUES (?, ?, ?, ?, ?, ?)`,
		post.ID, post.Title, post.Content,
		post.CreationDate, post.SubforumID, post.UserID)

	return err
}

func (store *PostStore) GetPostById(id uuid.UUID) (*model.Post, error) {
	post := &model.Post{}

	row := store.QueryRow(`
	SELECT id, title, content, creationdate, subforumid, userid 
	FROM posts
	WHERE id = ?`, id)

	err := row.Scan(&post.ID, &post.Title, &post.Content,
		&post.CreationDate, &post.SubforumID, &post.UserID)

	return post, err
}

func (store *PostStore) GetAllPostsByUserId(userID uuid.UUID) ([]model.Post, error) {
	posts := []model.Post{}

	rows, err := store.Query(`
	SELECT id, title, content, creationdate, subforumid, userid 
	FROM posts
	WHERE userid = ?`, userID)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var post model.Post

		err := rows.Scan(&post.ID, &post.Title, &post.Content,
			&post.CreationDate, &post.SubforumID, &post.UserID)

		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (store *PostStore) GetAllPostsBySubforumId(subforumID uuid.UUID) ([]model.Post, error) {
	posts := []model.Post{}

	rows, err := store.Query(`
	SELECT id, title, content, creationdate, subforumid, userid 
	FROM posts
	WHERE subforumid = ?`, subforumID)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var post model.Post

		err := rows.Scan(&post.ID, &post.Title, &post.Content,
			&post.CreationDate, &post.SubforumID, &post.UserID)

		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (store *PostStore) GetAllPosts() ([]model.Post, error) {
	var posts []model.Post

	rows, err := store.Query(`
	SELECT id, title, content, creationdate, subforumid, userid 
	FROM posts`)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var post model.Post

		err := rows.Scan(&post.ID, &post.Title, &post.Content,
			&post.CreationDate, &post.SubforumID, &post.UserID)

		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (store *PostStore) UpdatePost(post *model.Post) error {
	res, err := store.Exec(`
	UPDATE posts SET title = ?, content = ?, creationdate = ?, subforumid = ?
	WHERE id = ?`,
		post.Title,
		post.Content, post.CreationDate,
		post.SubforumID, post.UserID, post.ID)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("Post with id %v not found", post.ID)
	}
	return err
}

func (store *PostStore) DeletePost(id uuid.UUID) error {
	rows, err := store.Exec(`DELETE FROM posts WHERE Id = ?`, id)
	if err != nil {
		return err
	}

	n, err := rows.RowsAffected()
	if n == 0 {
		return fmt.Errorf("Post with id %v not found", id)
	}

	return err
}
