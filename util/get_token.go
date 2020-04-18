package util

import (
	"net/http"

	"git.01.alem.school/qjawko/forum/jwt"
	"git.01.alem.school/qjawko/forum/model"
)

func GetUser(r *http.Request) (*model.User, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		return nil, err
	}

	var user model.User
	if err := jwt.Unmarshal(cookie.Value, "supersecret", &user); err != nil {
		return nil, err
	}

	return &user, nil
}
