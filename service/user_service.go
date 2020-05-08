package service

import (
	"net/http"

	"git.01.alem.school/qjawko/forum/dao"
	"git.01.alem.school/qjawko/forum/http_errors"
	"git.01.alem.school/qjawko/forum/model"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userDao *dao.UserStore
}

func NewUserService(userDao *dao.UserStore) *UserService {
	return &UserService{userDao: userDao}
}

func (this *UserService) CreateUser(u *model.User) (*model.User, error) {
	if err := ValidateUser(u); err != nil {
		return nil, &http_errors.HttpError{Err: err, Code: http.StatusBadRequest}
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(u.Password), 0)
	if err != nil {
		return nil, &http_errors.HttpError{Err: err, Code: http.StatusInternalServerError}
	}
	u.ID = uuid.NewV4()
	u.Password = string(hashedPass)

	if err := this.userDao.CreateUser(u); err != nil {
		return nil, &http_errors.HttpError{Err: err, Code: http.StatusInternalServerError}
	}

	return u, nil
}

func (this *UserService) GetUserByID(id uuid.UUID) (*model.User, error) {
	user, err := this.userDao.GetUser(id)
	if err != nil {
		return nil, &http_errors.HttpError{Err: err, Code: http.StatusNotFound}
	}

	return user, nil
}

func (this *UserService) GetAll() ([]model.User, error) {
	users, err := this.userDao.GetAllUsers()
	if err != nil {
		return nil, &http_errors.HttpError{Err: err, Code: http.StatusInternalServerError}
	}

	return users, nil
}

func (this *UserService) GetUserByUsername(username string) (*model.User, error) {
	user, err := this.userDao.GetUserByUsername(username)

	if err != nil {
		return nil, &http_errors.HttpError{Err: err, Code: http.StatusBadRequest}
	}

	return user, nil
}

func (this *UserService) UpdateUser(u *model.User) (*model.User, error) {

	if u.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, &http_errors.HttpError{Err: err, Code: http.StatusInternalServerError}
		}

		u.Password = string(hash)
	}

	if err := this.userDao.UpdateUser(u); err != nil {
		return nil, &http_errors.HttpError{Err: err, Code: http.StatusInternalServerError}
	}

	return this.GetUserByID(u.ID)
}

func (this *UserService) DeleteUser(id uuid.UUID) error {
	if err := this.userDao.DeleteUser(id); err != nil {
		return &http_errors.HttpError{Err: err, Code: http.StatusInternalServerError}
	}

	return nil
}

// // Register q
// func Register(u model.User) *HTTPError {

// 	if err := validateUser(u); err != nil {
// 		return &HTTPError{err, http.StatusBadRequest}
// 	}

// 	// Generate hash from password
// 	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 0)
// 	if err != nil {
// 		return &HTTPError{err, http.StatusInternalServerError}
// 	}

// 	// Insert new user into db
// 	if err := userRepo.CreateUser(u.Name, u.Username, hash); err != nil {
// 		if strings.Contains(err.Error(), "UNIQUE") {
// 			return &HTTPError{err, http.StatusConflict}
// 		}

// 		return &HTTPError{err, http.StatusInternalServerError}
// 	}

// 	return nil
// }

// // Login q
// func Login(u model.User) *HTTPError {

// 	if err := validateLogin(u); err != nil {
// 		return &HTTPError{err, http.StatusBadRequest}
// 	}

// 	hash, err := repo.GetPassword(u.Username)
// 	if err != nil {
// 		return &HTTPError{err, http.StatusNotFound}
// 	}

// 	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(u.Password))
// 	if err != nil {
// 		return &HTTPError{err, http.StatusConflict}
// 	}

// 	// TODO: create token

// 	return nil
// }

// //4 operations CRUD for users
// func GetUserByID(id string) (*model.User, *HTTPError) {

// 	if len(strings.TrimSpace(Id)) == 0 {
// 		return nil, &HTTPError{errors.New("Invalid Id"), http.StatusBadRequest}
// 	}

// 	user, err := repo.GetUserByID(Id)
// 	if err != nil {
// 		return nil, &HTTPError{err, http.StatusInternalServerError}
// 	}
// 	return user, nil
// }

// //Бизнес логика в сервисах
// //Базы данных в репо
// func GetAllUsers() ([]model.User, *HTTPError) {

// 	users, err := repo.GetAllUsers()
// 	if err != nil {
// 		return nil, &HTTPError{err, http.StatusInternalServerError}
// 	}
// 	return users, nil
// }

// func Update(user model.User) (*model.User, *HTTPError) {

// 	if err := validateUser(user); err != nil {
// 		return nil, &HTTPError{err, http.StatusBadRequest}
// 	}

// 	updatesUser, err := repo.UpdateUser(user)
// 	if err != nil {
// 		return nil, &HTTPError{err, http.StatusInternalServerError}
// 	}
// 	return updatesUser, nil
// }

// func DeleteUser(id string) *HTTPError {
// 	err := repo.DeleteUser(Id)
// 	if err != nil {
// 		return &HTTPError{err, http.StatusInternalServerError}
// 	}
// 	return nil
// }

// func ChangePassword(id, password, repeatPassword, currentPassword string) *HTTPError {
// 	//когда пароли не одинаковы (new password, repeat password)
// 	if password != repeatPassword {
// 		return &HTTPError{errors.New("Passwords are not equal"), http.StatusBadRequest}
// 	}
// 	//когда длина пароля слишком короткая
// 	if len(password) < 6 {
// 		return &HTTPError{errors.New("Password is too short"), http.StatusBadRequest}
// 	}

// 	currentUser, err := repo.GetUserByID(id)
// 	if err != nil {
// 		return &HTTPError{err, http.StatusInternalServerError}
// 	}
// 	//если введенный пароль не правильный (введенный с нынешним)
// 	currPassword, err := bcrypt.GenerateFromPassword([]byte(currentPassword), 0)
// 	if err != nil {
// 		return &HTTPError{err, http.StatusInternalServerError}
// 	}

// 	if string(currPassword) != currentUser.Password {
// 		return &HTTPError{errors.New("Password is incorrect"), http.StatusBadRequest}
// 	}

// 	hash, err := bcrypt.GenerateFromPassword([]byte(password), 0)
// 	if err != nil {
// 		return &HTTPError{err, http.StatusInternalServerError}
// 	}

// 	err = repo.ChangePassword(id, string(hash))
// 	if err != nil {
// 		return &HTTPError{err, http.StatusInternalServerError}
// 	}
// 	return nil
// }
