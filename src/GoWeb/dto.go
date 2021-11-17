package main

import "mime/multipart"

type LoginPram struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ValidateUserParam struct {
	Password string `json:"password"`
}

type AddUserParam struct {
	Username string `json:"username"`
	Password string `json:"password"`
	NickName string `json:"nickName"`
}

type UpdateUserParam struct {
	Id       int    `json:"id"`
	NickName string `json:"nickName"`
	Password string `json:"password"`
}

type DeleteUserParam struct {
	Id int `json:"id"`
}

type AddLinkParam struct {
	Name string `json:"name"`
	Link string `json:"link"`
}

type DeleteLinkParam struct {
	Id int `json:"id"`
}

type UpdateLinkParam struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
}

type UpdateIntroductionParam struct {
	Content string `json:"content"`
}

type DownloadToolsDTO struct {
	Id int `json:"id"`
}

type GetLinkDTO struct {
	Type string `json:"type"`
}

type SubscribeDTO struct {
	Type string `json:"type"`
	Sign string `json:"sign"`
}

type CreateSignDTO struct {
	Type   string `json:"type"`
	UserId int
}

type AddImageDTO struct {
	Image *multipart.FileHeader
}

type DeleteImageDTO struct {
	Id int `json:"id"`
}
