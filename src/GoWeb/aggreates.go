package main

import (
	"fmt"
	"time"
)

type jsonTime time.Time

type AuthUser struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	NickName string `json:"nickName"`
	Role     string `json:"role"`
}

type LinkInformation struct {
	Id         int      `json:"id"`
	Name       string   `json:"name"`
	Link       string   `json:"link"`
	UpdateTime jsonTime `json:"updateTime"`
}

func (t jsonTime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02 15:04"))
	return []byte(stamp), nil
}
