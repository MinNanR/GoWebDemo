package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func db() *sql.DB {
	db, _ := sql.Open("mysql", "Minnan:minnan@(minnan.site:3306)/link_manager")
	return db
}

var conn = db()

func GetUserByUsername(username string) (AuthUser, error) {
	var authUser AuthUser
	err := conn.QueryRow("select id id, username username, password password from user where username = ?",
		username).Scan(&authUser.Id, &authUser.Username, &authUser.Password)
	if err != nil {
		return authUser, err
	}
	return authUser, nil
}
