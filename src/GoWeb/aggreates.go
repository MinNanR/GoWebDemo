package main

type AuthUser struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
