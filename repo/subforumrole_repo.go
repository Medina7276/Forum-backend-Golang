package repo

import (
	"database/sql"
	"fmt"

	"git.01.alem.school/qjawko/forum/model"
	uuid "github.com/satori/go.uuid"
)

type SubforumRoleStore struct {
	*sql.DB
}

func NewSubforumRoleStore(db *sql.DB) *SubforumRoleStore {
	return &SubforumRoleStore{DB: db}
}

func (store *SubforumRoleStore) Create(subforumRole *model.SubforumRole) error {
	_, err := store.Exec(`INSERT INTO subforumrole (id, userId, role, subforumId) VALUES (?, ?, ?, ?)`,
		subforumRole.ID,
		subforumRole.UserID,
		subforumRole.Role,
		subforumRole.SubforumID)

	return err
}

func (store *SubforumRoleStore) GetById(id uuid.UUID) (*model.SubforumRole, error) {
	role := &model.SubforumRole{}

	row := store.QueryRow(`SELECT id, userId, role, subforumId FROM subforumrole WHERE id = ?`, id)

	err := row.Scan(role.ID, role.UserID, role.Role, role.SubforumID)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (store *SubforumRoleStore) GetBySubforumId(id uuid.UUID) ([]model.SubforumRole, error) {
	var roles []model.SubforumRole

	rows, err := store.Query(`SELECT id, userId, role, subforumId FROM subforumrole WHERE subforumId = ?`, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var role *model.SubforumRole

		err := rows.Scan(role.ID, role.UserID, role.Role, role.SubforumID)
		if err != nil {
			return nil, err
		}

		roles = append(roles, *role)
	}

	return roles, nil
}

func (store *SubforumRoleStore) GetAll() ([]model.SubforumRole, error) {
	var roles []model.SubforumRole

	rows, err := store.Query(`SELECT id, userId, role, subforumId FROM subforumrole`)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var role *model.SubforumRole

		err := rows.Scan(role.ID, role.UserID, role.Role, role.SubforumID)
		if err != nil {
			return nil, err
		}

		roles = append(roles, *role)
	}

	return roles, nil
}

func (store *SubforumRoleStore) Update(role *model.SubforumRole) error {
	res, err := store.Exec(`UPDATE userId = ?, role = ?, subforumId 
							FROM subforumRole
							WHERE id = ?`, role.UserID, role.Role, role.SubforumID, role.ID)

	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n == 0 {
		return fmt.Errorf("Role with id %s not found", role.ID)
	}

	return nil
}

func (store *SubforumRoleStore) Delete(id uuid.UUID) error {
	rows, err := store.Exec(`DELETE FROM subforumRole WHERE Id = ?`, id)
	if err != nil {
		return err
	}

	n, err := rows.RowsAffected()
	if err != nil {
		return err
	}

	if n == 0 {
		return fmt.Errorf("Role with id %v not found", id)
	}

	return nil
}
