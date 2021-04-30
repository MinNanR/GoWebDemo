package main

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
	NickName string `json:"username"`
	Password string `json:"password"`
}
