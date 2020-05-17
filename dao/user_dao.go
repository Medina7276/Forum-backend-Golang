package dao

import (
	"database/sql"
	"fmt"

	"git.01.alem.school/qjawko/forum/model"
	uuid "github.com/satori/go.uuid"
)

type UserStore struct {
	*sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{DB: db}
}

func (store *UserStore) CreateUser(user *model.User) error {
	_, err := store.Exec(`INSERT INTO users (id, username, email, name, password, avatarurl) 
	VALUES (?, ?, ?, ?, ?, ?)`,
		user.ID, user.Username, user.Email, user.Name,
		user.Password, user.AvatarURL)

	return err
}

func (store *UserStore) GetUser(id uuid.UUID) (*model.User, error) {
	user := &model.User{}

	row := store.QueryRow(`
	SELECT id, username, email, name, password, avatarurl 
	FROM users
	WHERE id = ?`, id)

	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Name, &user.Password, &user.AvatarURL)

	return user, err
}

func (store *UserStore) GetUserByUsername(username string) (*model.User, error) {
	user := &model.User{}

	row := store.QueryRow(`
	SELECT id, username, email, name, password, avatarurl 
	FROM users
	WHERE username = ?`, username)

	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Name,
		&user.Password, &user.AvatarURL)

	return user, err
}

func (store *UserStore) GetAllUsers() ([]model.User, error) {
	users := []model.User{}

	rows, err := store.Query(`
	SELECT id, username, email, name, password, avatarurl 
	FROM users`)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user model.User

		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Name,
			&user.Password, &user.AvatarURL)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (store *UserStore) UpdateUser(user *model.User) error {
	res, err := store.Exec(`UPDATE users SET username = ?, email = ?, name = ?, 
	password = ?, avatarurl = ?, role = ?  
	WHERE id = ?`, user.Username, user.Email, user.Name,
		user.Password, user.AvatarURL, user.ID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if n == 0 {
		return fmt.Errorf("User with id %v not found", user.ID)
	}

	return err
}

func (store *UserStore) DeleteUser(id uuid.UUID) error {
	rows, err := store.Exec(`DELETE FROM users WHERE Id = ?`, id)
	if err != nil {
		return err
	}

	n, err := rows.RowsAffected()
	if n == 0 {
		return fmt.Errorf("User with id %v not found", id)
	}

	return err
}
