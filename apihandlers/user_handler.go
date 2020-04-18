package apihandlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"git.01.alem.school/qjawko/forum/http_errors"
	"git.01.alem.school/qjawko/forum/middleware"
	"git.01.alem.school/qjawko/forum/model"
	"git.01.alem.school/qjawko/forum/service"
	uuid "github.com/satori/go.uuid"
)

type UserHandler struct {
	UserService *service.UserService
	Endpoint    string
}

func NewUserHandler(endpoint string, service *service.UserService) *UserHandler {
	return &UserHandler{
		UserService: service,
		Endpoint:    endpoint,
	}
}

func (uh *UserHandler) Route(w http.ResponseWriter, r *http.Request) {
	var funcToCall func(http.ResponseWriter, *http.Request)

	switch r.Method {
	case http.MethodGet:
		_, ok := r.URL.Query()["username"]
		if ok {
			fmt.Println("username is reached")
			funcToCall = uh.GetUserByUsername
		} else {
			endpoint := r.URL.Path[len(uh.Endpoint):]
			if len(endpoint) == 0 {
				funcToCall = uh.GetAll
			} else {
				funcToCall = uh.GetUserByID
			}
		}

	case http.MethodPost:
		funcToCall = middleware.CheckForAdmin(http.HandlerFunc(uh.CreateUser)).ServeHTTP
	case http.MethodPut:
		funcToCall = uh.Update
	case http.MethodDelete:
		funcToCall = middleware.CheckForModerator(http.HandlerFunc(uh.Delete)).ServeHTTP
	default:
		http.Error(w, "Route Not found", http.StatusNotFound)
		return
	}

	funcToCall(w, r)
}

//CreateUser q
func (uh *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	created, err := uh.UserService.CreateUser(&user)
	if err != nil {
		httpErr := err.(*http_errors.HttpError)
		http.Error(w, httpErr.Error(), httpErr.Code)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

//GetAll qwe
func (uh *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := uh.UserService.GetAll()
	if err != nil {
		httpErr := err.(*http_errors.HttpError)
		http.Error(w, httpErr.Error(), httpErr.Code)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func (uh *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len(uh.Endpoint):]
	user, err := uh.UserService.GetUserByID(uuid.FromStringOrNil(id))
	if err != nil {
		httpErr := err.(*http_errors.HttpError)
		http.Error(w, httpErr.Error(), httpErr.Code)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (uh *UserHandler) GetUserByUsername(w http.ResponseWriter, r *http.Request) {
	username, _ := r.URL.Query()["username"]
	user, err := uh.UserService.GetUserByUsername(username[0])
	if err != nil {
		httpErr := err.(*http_errors.HttpError)
		http.Error(w, httpErr.Error(), httpErr.Code)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (uh *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Bad Body", http.StatusBadRequest)
		return
	}

	idFromURL := r.URL.Path[len(uh.Endpoint):]
	user.ID = uuid.FromStringOrNil(idFromURL)

	updated, err := uh.UserService.UpdateUser(&user)
	if err != nil {
		httpErr := err.(*http_errors.HttpError)
		http.Error(w, httpErr.Error(), httpErr.Code)
		return
	}

	json.NewEncoder(w).Encode(updated)
}

func (uh *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len(uh.Endpoint):]
	if err := uh.UserService.DeleteUser(uuid.FromStringOrNil(id)); err != nil {
		httpErr := err.(*http_errors.HttpError)
		http.Error(w, httpErr.Error(), httpErr.Code)
	}
}

// RegisterHandler handle register
// func RegisterHandler() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var u model.User
// 		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return
// 		}

// 		normalizeUser(&u)

// 		err := service.Register(u)
// 		if err != nil {
// 			http.Error(w, err.Err.Error(), err.Code)
// 			return
// 		}

// 		// Finally return new user in json
// 		u.Password = ""
// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(u)
// 	}
// }

// // LoginHandler handle login
// func LoginHandler() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var u model.User
// 		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return
// 		}

// 		normalizeUser(&u)

// 		service.Login(u)

// 		// TODO send token to user
// 	}
// }

// func normalizeUser(u *model.User) {
// 	trimcharset := " \n\r\t"
// 	u.Username = strings.Trim(u.Username, trimcharset)
// }
