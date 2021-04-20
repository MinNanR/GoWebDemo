package main

import (
	"golang.org/x/crypto/bcrypt"
)

func validateUser(dto LoginPram) (AuthUser, error) {
	user, err := loadUserByUserName(dto.Username)
	if err != nil {
		return AuthUser{}, LoginError{message: "用户不存在"}
	}
	check := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password))
	if check != nil {
		return AuthUser{}, LoginError{message: "密码错误"}
	}
	return user, nil
}

func loadUserByUserName(username string) (AuthUser, error) {
	user, err := SelectUserByUsername(username)
	if err != nil {
		return AuthUser{}, LoginError{message: "用户不存在"}
	}
	return user, nil
}

func getLinkList() ([]LinkInformation, error) {
	return SelectAllLink()
}
