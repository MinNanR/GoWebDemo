package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func db() *sql.DB {
	var username, password, host, port, database, driverName = dataSource.Username, dataSource.Password, dataSource.Host,
		dataSource.Port, dataSource.DataBase, dataSource.DriverName
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?loc=Asia%%2FShanghai&parseTime=true", username, password, host,
		port, database)
	db, _ := sql.Open(driverName, dataSourceName)
	return db
}

var conn = db()

type UserDB struct {
}

func (userDb UserDB) SelectUserByUsername(username string) (AuthUser, error) {
	var authUser AuthUser
	err := conn.QueryRow("select id id, username username, password password, "+
		"nick_name nickName, role role from user where username = ?",
		username).Scan(&authUser.Id, &authUser.Username, &authUser.Password, &authUser.NickName, &authUser.Role)
	if err != nil {
		return authUser, err
	}
	return authUser, nil
}

func (userDb UserDB) insert(user AuthUser) error {
	stmt, err := conn.Prepare("insert into user(username, password, nick_name, role) value (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(user.Username, user.Password, user.NickName, user.Role)
	return err
}

func (userDb UserDB) updateUser(user AuthUser) error {
	stmt, err := conn.Prepare("update user set password = ?, nick_name = ? where id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(user.Password, user.NickName, user.Id)
	return err
}
func (userDb UserDB) SelectById(id int) (AuthUser, error) {
	var user AuthUser
	err := conn.QueryRow("select (id, username, password, nick_name, role) from user where id = ?", id).Scan(&user.Id,
		&user.Username, &user.NickName, &user.Role)
	if err != nil {
		return user, err
	}
	return user, nil
}

type LinkDB struct {
}

func (linkDB LinkDB) SelectAllLink() ([]LinkInformation, error) {
	rows, err := conn.Query("select id id,name name, link link, update_time updateTime from link")
	if err != nil {
		return nil, err
	}
	LinkInformationList := make([]LinkInformation, 0)
	for rows.Next() {
		line := new(LinkInformation)
		rows.Scan(&line.Id, &line.Name, &line.Link, &line.UpdateTime)
		LinkInformationList = append(LinkInformationList, *line)
	}
	return LinkInformationList, nil
}
