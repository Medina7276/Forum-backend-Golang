package operations

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"git.01.alem.school/qjawko/forum/dto"
	"git.01.alem.school/qjawko/forum/http_errors"
	"git.01.alem.school/qjawko/forum/jwt"
	"git.01.alem.school/qjawko/forum/model"
	"git.01.alem.school/qjawko/forum/service"
	"golang.org/x/crypto/bcrypt"
)

type UserOperations struct {
	userService *service.UserService
	Endpoint    string
}

func NewUserOperations(userService *service.UserService) *UserOperations {
	return &UserOperations{userService: userService}
}

func (this *UserOperations) Register(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// user.Role = model.ROLE_USER
	_, err := this.userService.CreateUser(&user)
	if err != nil {
		httpErr := err.(*http_errors.HttpError)
		http.Error(w, httpErr.Error(), httpErr.Code)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (this *UserOperations) Login(w http.ResponseWriter, r *http.Request) {
	var creds dto.AuthCredentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := this.userService.GetUserByUsername(creds.Username)
	if err != nil {
		httpErr := err.(*http_errors.HttpError)
		http.Error(w, httpErr.Error(), httpErr.Code)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	userWithClaims := struct {
		model.User
		jwt.DefaultClaims
	}{
		User:          *user,
		DefaultClaims: jwt.DefaultClaimsWithDefaultExp(),
	}

	token, err := jwt.Generate(userWithClaims, "supersecret")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: time.Now().Add(1 * time.Hour),
	}
	http.SetCookie(w, cookie)
	json.NewEncoder(w).Encode(userWithClaims)
}

func (this *UserOperations) Me(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized) //записываем ошибку в респонс
		return
	}
	var user model.User

	fmt.Println(token.Value)
	err = jwt.Unmarshal(token.Value, "supersecret", &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(&user) //создает новый encoder который записывает все hhttp writer, а потом записывает туда user
}

func (this *UserOperations) DeleteYourself(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token") //прочли токен
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized) //записываем ошибку в респонс
		return
	}

	var user model.User
	err = jwt.Unmarshal(token.Value, "supersecret", &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := r.URL.Path[len(this.Endpoint):]

	if id == user.ID.String() {
		err := this.userService.DeleteUser(user.ID)
		if err != nil {
			httpErr := err.(*http_errors.HttpError)
			http.Error(w, httpErr.Error(), httpErr.Code)
			return
		}
	} else {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
