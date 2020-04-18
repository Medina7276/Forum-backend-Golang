package repo

import (
	"database/sql"
	"fmt"

	"git.01.alem.school/qjawko/forum/model"
	uuid "github.com/satori/go.uuid"
)

//SubforumStore erfn
type SubforumStore struct {
	*sql.DB
}

//NewSubforumStore dfrf
func NewSubforumStore(db *sql.DB) *SubforumStore {
	return &SubforumStore{DB: db}
}

//CreateSubforum dfjn
func (store *SubforumStore) CreateSubforum(subforum *model.Subforum) error {
	_, err := store.Exec(`INSERT INTO subforum (id, name, description, parentid) 
	VALUES (?, ?, ?)`,
		subforum.ID, subforum.Name, subforum.Description, subforum.ParentID)
	return err
}

// GetSubforum dlfj
func (store *SubforumStore) GetSubforumById(id uuid.UUID) (*model.Subforum, error) {
	subforum := &model.Subforum{}

	row := store.QueryRow(`
	SELECT id, name, description, parentid
	FROM subforum
	WHERE id = ?`, id)

	err := row.Scan(&subforum.ID, &subforum.Name, &subforum.Description, &subforum.ParentID)
	return subforum, err
}

func (store *SubforumStore) GetSubforumByParentId(parentId uuid.UUID) ([]model.Subforum, error) {
	subforums := []model.Subforum{}

	rows, err := store.Query(`
	SELECT id, name, description, parentid
	FROM subforum
	WHERE parentid = ?`, parentId)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var subforum model.Subforum

		err := rows.Scan(&subforum.ID, &subforum.Name, &subforum.Description, &subforum.ParentID)
		if err != nil {
			return nil, err
		}

		subforums = append(subforums, subforum)
	}
	return subforums, nil
}

//GetSubforumByName dflj
func (store *SubforumStore) GetSubforumByName(name string) (*model.Subforum, error) {
	subforum := &model.Subforum{}

	row := store.QueryRow(`
	SELECT id, name, description, parentid
	FROM subdorum
	WHERE name = ?`, name)

	err := row.Scan(&subforum.ID, &subforum.Name, &subforum.Description, &subforum.ParentID)
	return subforum, err
}

//GetAllSubforums dkfgj
func (store *SubforumStore) GetAllSubforums() ([]model.Subforum, error) {
	subforums := []model.Subforum{}

	rows, err := store.Query(`
	SELECT id, name, description, parentid
	FROM subforum`)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var subforum model.Subforum

		err := rows.Scan(&subforum.ID, &subforum.Name, &subforum.Description, &subforum.ParentID)
		if err != nil {
			return nil, err
		}

		subforums = append(subforums, subforum)
	}
	return subforums, nil
}

//UpdateSubforum erlgj
func (store *SubforumStore) UpdateSubforum(subforum *model.Subforum) error {
	res, err := store.Exec(`UPDATE subforum SET name = ?, description = ?, parentid = ?
	WHERE id = ?`, subforum.Name, subforum.Description, subforum.ID, subforum.ParentID)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n == 0 {
		return fmt.Errorf("Subforum with id %v not found", subforum.ID)
	}
	return err
}

//DeleteSubforum dlfj
func (store *SubforumStore) DeleteSubforum(id uuid.UUID) error {
	rows, err := store.Exec(`
	DELETE FROM subforum WHERE id = ?`, id)
	if err != nil {
		return err
	}

	n, err := rows.RowsAffected()
	if n == 0 {
		return fmt.Errorf("Subforum with id %v not found", id)
	}
	return err
}
