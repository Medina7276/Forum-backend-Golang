package repo

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	"git.01.alem.school/qjawko/forum/model"
	uuid "github.com/satori/go.uuid"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var store *UserStore

func init() {
	var err error
	db, err = sql.Open("sqlite3", "../test/forum.db")
	if err != nil {
		log.Fatal(err)
	}

	if err := Setup(db); err != nil {
		log.Fatal(err)
	}

	store = NewUserStore(db)
}

func TestCreateUser(t *testing.T) {
	user := &model.User{
		ID:       uuid.NewV4(),
		Username: "rtabuloasdasdvqwe",
		Email:    "rasul.tabulov@gasdasdmail.coqwem",
	}

	err := store.CreateUser(user)
	if err != nil {
		log.Println(err)
		t.Fail()
	}

	all, err := store.GetAllUsers()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(all)
}

func TestGetUser(t *testing.T) {
	id, err := uuid.FromString("b3628432-d90d-471f-8c21-cf3cc6dd42f6")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	user, err := store.GetUser(id)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	fmt.Println(user)
}

func TestGetUserByUsername(t *testing.T) {
	user, err := store.GetUserByUsername("rtabulov")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	fmt.Println(user)
}

func TestDeleteUser(t *testing.T) {
	id, _ := uuid.FromString("b3628432-d90d-471f-8c21-cf3cc6dd42f6")

	err := store.DeleteUser(id)
	if err != nil {
		log.Println(err)
		t.Fail()
	}

	_, err = store.GetUser(id)
	if err == nil {
		fmt.Println("Did not delete")
		t.Fail()
	}

}

func TestUpdateUser(t *testing.T) {

	id, _ := uuid.FromString("b3628432-d90d-471f-8c21-cf3cc6dd42f6")
	user := &model.User{
		ID:       id,
		Username: "updated?",
		Email:    "email updated?",
	}

	err := store.UpdateUser(user)
	if err != nil {
		log.Println(err)
		t.Fail()
	}

	updated, err := store.GetUser(user.ID)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(updated)
}
