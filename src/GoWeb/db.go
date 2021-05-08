package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
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
	err := conn.QueryRow("select id, username, password, nick_name, role from user where id = ?", id).Scan(&user.Id,
		&user.Username, &user.Password, &user.NickName, &user.Role)
	if err != nil {
		return user, errors.New("用户不存在")
	}
	return user, nil
}

func (userDb UserDB) SelectList() ([]AuthUser, error) {
	rows, err := conn.Query("select id, username, nick_name from user")
	if err != nil {
		return nil, err
	}
	userInformationList := make([]AuthUser, 0)
	for rows.Next() {
		line := new(AuthUser)
		rows.Scan(&line.Id, &line.Username, &line.NickName)
		userInformationList = append(userInformationList, *line)
	}
	return userInformationList, nil
}

func (userDb UserDB) DeleteById(id int) error {
	stmt, err := conn.Prepare("delete from user where id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	return err
}

type LinkDB struct {
}

func (linkDb LinkDB) SelectAllLink() ([]LinkInformation, error) {
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

func (linkDb LinkDB) insert(link LinkInformation) error {
	stmt, err := conn.Prepare("insert into link (name, link, update_time) value (?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(link.Name, link.Link, time.Time(link.UpdateTime))
	return err
}

func (linkDb LinkDB) deleteById(id int) error {
	stmt, err := conn.Prepare("delete from link where id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	return err
}

func (linkDb LinkDB) updateLink(information LinkInformation) error {
	stmt, err := conn.Prepare("update link set name = ?, link =? ,update_time = ? where id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(information.Name, information.Link, time.Now(), information.Id)
	return err
}

func (linkDb LinkDB) selectById(id int) (LinkInformation, error) {
	var link LinkInformation
	err := conn.QueryRow("select id id, name name, link link from link where id = ?", id).Scan(&link.Id, &link.Name,
		&link.Link)
	if err != nil {
		return LinkInformation{}, nil
	}
	return link, err
}

type ImageDB struct {
}

func (imageDb ImageDB) selectList() ([]Image, error) {
	rows, err := conn.Query("select id, url from image")
	if err != nil {
		return nil, err
	}
	imageList := make([]Image, 0)
	for rows.Next() {
		image := Image{}
		rows.Scan(&image.Id, &image.Url)
		imageList = append(imageList, image)
	}
	return imageList, nil
}

type ToolsDB struct {
}

func (ToolsDB) selectList() ([]Tools, error) {
	rows, err := conn.Query("select id, file_name, extension, size, update_time from tools")
	if err != nil {
		return nil, err
	}
	var toolsList []Tools
	for rows.Next() {
		tools := Tools{}
		rows.Scan(&tools.Id, &tools.FileName, &tools.Extension, &tools.Size, &tools.UpdateTime)
		toolsList = append(toolsList, tools)
	}
	return toolsList, nil
}

func (ToolsDB) selectById(id int) (Tools, error) {
	var tools Tools
	err := conn.QueryRow("select id, file_name, extension, oss_key, size, update_time from tools where id = ?",
		id).Scan(&tools.Id, &tools.FileName, &tools.Extension, &tools.OssKey, &tools.Size, &tools.UpdateTime)
	if err != nil {
		return tools, err
	}
	return tools, nil
}
