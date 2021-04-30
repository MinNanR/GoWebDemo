package main

import (
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	db UserDB
}

func (service UserService) validateUser(dto LoginPram) (AuthUser, error) {
	user, err := service.loadUserByUserName(dto.Username)
	if err != nil {
		return AuthUser{}, LoginError{message: "用户不存在"}
	}
	check := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password))
	if check != nil {
		return AuthUser{}, LoginError{message: "密码错误"}
	}
	return user, nil
}

func (service UserService) loadUserByUserName(username string) (AuthUser, error) {
	user, err := service.db.SelectUserByUsername(username)
	if err != nil {
		return AuthUser{}, LoginError{message: "用户不存在"}
	}
	return user, nil
}

func (service UserService) addUser(dto AddUserParam) error {
	encodedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}
	user := AuthUser{
		Username: dto.Username,
		Password: string(encodedPassword),
		Role:     "USER",
		NickName: dto.NickName,
	}
	return service.db.insert(user)
}

type LinkService struct {
	db LinkDB
}

func (service LinkService) getLinkList() ([]LinkInformation, error) {
	return service.db.SelectAllLink()
}
