package service

import (
	"net/http"

	"git.01.alem.school/qjawko/forum/dao"
	"git.01.alem.school/qjawko/forum/http_errors"
	"git.01.alem.school/qjawko/forum/model"
	uuid "github.com/satori/go.uuid"
)

type SubforumService struct {
	subforumDao         *dao.SubforumStore
	SubforumRoleService *SubforumRoleService
}

func NewSubforumService(subforumDao *dao.SubforumStore, service *SubforumRoleService) *SubforumService {
	return &SubforumService{
		subforumDao:         subforumDao,
		SubforumRoleService: service,
	}
}

func (this *SubforumService) CreateSubforum(s *model.Subforum) (*model.Subforum, error) {
	s.ID = uuid.NewV4()
	if err := this.subforumDao.CreateSubforum(s); err != nil {
		return nil, &http_errors.HttpError{Err: err, Code: http.StatusInternalServerError}
	}

	sroles, err := this.SubforumRoleService.GetBySubforumId(s.ParentID)
	if err != nil {
		return nil, &http_errors.HttpError{Err: err, Code: http.StatusInternalServerError}
	}

	for _, srole := range sroles {
		srole.ID = uuid.NewV4()
		srole.SubforumID = s.ID
		_, err := this.SubforumRoleService.Create(&srole)
		if err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (this *SubforumService) GetSubforumById(id uuid.UUID) (*model.Subforum, error) {
	subforum, err := this.subforumDao.GetSubforumById(id)
	if err != nil {
		return nil, &http_errors.HttpError{Err: err, Code: http.StatusNotFound}
	}
	return subforum, nil
}

func (this *SubforumService) GetSubforumsByParentId(parentid uuid.UUID) ([]model.Subforum, error) {
	subforums, err := this.subforumDao.GetSubforumByParentId(parentid)
	if err != nil {
		return nil, &http_errors.HttpError{Err: err, Code: http.StatusNotFound}
	}
	return subforums, nil
}

func (this *SubforumService) GetAllSubforums() ([]model.Subforum, error) {
	subforums, err := this.subforumDao.GetAllSubforums()
	if err != nil {
		return nil, &http_errors.HttpError{Err: err, Code: http.StatusInternalServerError}
	}
	return subforums, nil
}

func (this *SubforumService) GetSubforumByName(name string) (*model.Subforum, error) {
	subforum, err := this.subforumDao.GetSubforumByName(name)
	if err != nil {
		return nil, &http_errors.HttpError{Err: err, Code: http.StatusBadRequest}
	}
	return subforum, nil
}

func (this *SubforumService) UpdateSubforum(s *model.Subforum) (*model.Subforum, error) {
	if err := this.subforumDao.UpdateSubforum(s); err != nil {
		return nil, &http_errors.HttpError{Err: err, Code: http.StatusInternalServerError}
	}
	return this.GetSubforumById(s.ID)
}

func (this *SubforumService) DeleteSubforum(id uuid.UUID) error {
	if err := this.subforumDao.DeleteSubforum(id); err != nil {
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
