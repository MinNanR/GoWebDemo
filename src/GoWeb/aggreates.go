package main

import (
	"fmt"
	"time"
)

type JsonTime time.Time

func now() JsonTime {
	return JsonTime(time.Now())
}

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
	UpdateTime JsonTime `json:"updateTime"`
}

func (t JsonTime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02 15:04"))
	return []byte(stamp), nil
}

type Image struct {
	Id  int    `json:"id"`
	Url string `json:"url"`
}

type Tools struct {
	Id         int      `json:"id"`
	FileName   string   `json:"fileName"`
	OssKey     string   `json:"ossKey"`
	Extension  string   `json:"extension"`
	Size       int64    `json:"size"`
	UpdateTime JsonTime `json:"updateTime"`
}
