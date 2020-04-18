package repo

import (
	"database/sql"
)

func Setup(db *sql.DB) error {
	if err := setupUsers(db); err != nil {
		return err
	}
	if err := setupSubforum(db); err != nil {
		return err
	}
	if err := setupPosts(db); err != nil {
		return err
	}
	if err := setupSubforumRole(db); err != nil {
		return err
	}
	if err := setupLikes(db); err != nil {
		return err
	}

	return nil
}

func setupLikes(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS "likes" (
		"id"	TEXT NOT NULL UNIQUE,
		"userid"	TEXT NOT NULL,
		"postid"	TEXT NOT NULL,
		"isupvote"	INTEGER NOT NULL,
		FOREIGN KEY("postid") REFERENCES "posts"("id") ON DELETE CASCADE,
		FOREIGN KEY("userid") REFERENCES "users"("id") ON DELETE CASCADE,
		PRIMARY KEY("id")
	)`)

	return err
}

func setupUsers(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS "users" (
		"id"	TEXT NOT NULL UNIQUE,
		"name"	TEXT NOT NULL,
		"email"	TEXT NOT NULL UNIQUE,
		"username"	TEXT NOT NULL UNIQUE,
		"password"	TEXT NOT NULL,
		"avatarurl"	TEXT,
		"role"	INTEGER NOT NULL,
		PRIMARY KEY("id")
	)`)

	return err

}

func setupPosts(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS "posts" (
		"id"	TEXT NOT NULL UNIQUE,
		"parentid"	TEXT NOT NULL,
		"title"	TEXT NOT NULL,
		"content"	TEXT NOT NULL,
		"creationdate"	INTEGER NOT NULL,
		"subforumid"	TEXT NOT NULL,
		"userid"	TEXT NOT NULL,
		FOREIGN KEY("parentid") REFERENCES "posts"("id") ON DELETE CASCADE,
		FOREIGN KEY("userid") REFERENCES "users"("id") ON DELETE CASCADE,
		PRIMARY KEY("id"),
		FOREIGN KEY("subforumid") REFERENCES "subforum"("id") ON DELETE CASCADE
	)`)

	return err

}
func setupSubforum(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS "subforum" (
		"id"	TEXT NOT NULL UNIQUE,
		"name"	TEXT NOT NULL,
		"description"	TEXT NOT NULL,
		PRIMARY KEY("id")
	)`)

	return err

}
func setupSubforumRole(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS "subforumrole" (
		"id"	TEXT NOT NULL UNIQUE,
		"userid"	TEXT NOT NULL,
		"role"	INTEGER NOT NULL,
		"subforumid"	TEXT NOT NULL,
		FOREIGN KEY("subforumid") REFERENCES "subforum"("id") ON DELETE CASCADE,
		PRIMARY KEY("id"),
		FOREIGN KEY("userid") REFERENCES "users"("id") ON DELETE CASCADE
	)`)

	return err

}
