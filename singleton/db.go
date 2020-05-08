package singleton

import (
	"database/sql"
	"sync"
)

var instance *sql.DB
var once sync.Once

func GetDBInstance() *sql.DB {
	once.Do(func() {
		var err error
		instance, err = sql.Open("sqlite3", "./forum.db") //todo config
		if err != nil {
			panic(err)
		}
	})

	return instance
}
