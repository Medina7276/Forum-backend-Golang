package service

import (
	"net/http"

	"git.01.alem.school/qjawko/forum/dao"
	"git.01.alem.school/qjawko/forum/http_errors"
	"git.01.alem.school/qjawko/forum/model"
	uuid "github.com/satori/go.uuid"
)

type SubforumRoleService struct {
	subforumService *SubforumService
	subforumRoleDao *dao.SubforumRoleStore
}

func (sr *SubforumRoleService) Create(role *model.SubforumRole) (*model.SubforumRole, error) {
	err := sr.applyToAll(role, sr.subforumRoleDao.Create)

	return role, err
}

func (sr *SubforumRoleService) GetById(id uuid.UUID) (*model.SubforumRole, error) {
	srole, err := sr.subforumRoleDao.GetById(id)
	if err != nil {
		return nil, &http_errors.HttpError{Err: err, Code: http.StatusInternalServerError}
	}

	return srole, err
}

func (sr *SubforumRoleService) GetBySubforumId(id uuid.UUID) ([]model.SubforumRole, error) {
	sroles, err := sr.subforumRoleDao.GetBySubforumId(id)
	if err != nil {
		return nil, &http_errors.HttpError{Err: err, Code: http.StatusInternalServerError}
	}

	return sroles, err
}

func (sr *SubforumRoleService) Update(role *model.SubforumRole) (*model.SubforumRole, error) {
	err := sr.applyToAll(role, sr.subforumRoleDao.Update)

	return role, err
}

func (sr *SubforumRoleService) Delete(id uuid.UUID) error {
	role, err := sr.GetById(id)
	if err != nil {
		return &http_errors.HttpError{Err: err, Code: http.StatusInternalServerError}
	}

	err = sr.applyToAll(role, sr.deleteByRole)

	return err
}

func (sr *SubforumRoleService) deleteByRole(role *model.SubforumRole) error {
	return sr.subforumRoleDao.Delete(role.ID)
}

func (sr *SubforumRoleService) applyToAll(role *model.SubforumRole, f func(*model.SubforumRole) error) error {
	subforum, err := sr.subforumService.GetSubforumById(role.SubforumID)
	if err != nil {
		return err
	}

	var subforums []model.Subforum
	subforums = append(subforums, *subforum)

	for i := 0; i < len(subforums); i++ {
		children, err := sr.subforumService.GetSubforumsByParentId(subforums[i].ID)
		if err != nil {
			return err
		}

		role.SubforumID = subforums[i].ID
		err = f(role)
		if err != nil {
			return err
		}

		subforums = append(subforums, children...)
	}

	return nil
}
