package main

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
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

func (service UserService) getUserInformationList() ([]AuthUser, error) {
	return service.db.SelectList()
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

func (service UserService) updateUser(dto UpdateUserParam) error {
	user, err := service.db.SelectById(dto.Id)
	if err != nil {
		return err
	}
	if len(dto.NickName) > 0 {
		user.NickName = dto.NickName
	}
	if len(dto.Password) > 0 {
		encodedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.MinCost)
		if err != nil {
			return err
		}
		user.Password = string(encodedPassword)
	}
	return service.db.updateUser(user)
}

func (service UserService) deleteUser(id int) error {
	return service.db.DeleteById(id)
}

type LinkService struct {
	db LinkDB
}

func (service LinkService) getLinkList() ([]LinkInformation, error) {
	return service.db.SelectAllLink()
}

func (service LinkService) addLink(param AddLinkParam) error {
	link := LinkInformation{Name: param.Name, Link: param.Link, UpdateTime: JsonTime(time.Now())}
	return service.db.insert(link)
}

func (service LinkService) deleteLink(id int) error {
	return service.db.deleteById(id)
}

func (service LinkService) updateLink(param UpdateLinkParam) error {
	link, err := service.db.selectById(param.Id)
	if err != nil {
		return err
	}
	if len(param.Link) > 0 {
		link.Link = param.Link
	}
	if len(param.Name) > 0 {
		link.Name = param.Name
	}
	return service.db.updateLink(link)
}

type IntroductionService struct {
	db ImageDB
}

func (service IntroductionService) updateIntroduction(content string) {
	rc.setValue("introduction1", content)
}

func (service IntroductionService) getIntroduction() string {
	return rc.getValue("introduction1")
}

func (service IntroductionService) getImageList() ([]Image, error) {
	imageList, err := service.db.selectList()
	if err != nil {
		return nil, err
	}
	for index, _ := range imageList {
		imageList[index].Url = aliyunConfig.baseUrl + imageList[index].Url
	}
	return imageList, nil
}

type ToolsService struct {
	db ToolsDB
}

func (service ToolsService) getToolsList() ([]Tools, error) {
	return service.db.selectList()
}

func (service ToolsService) downloadTools(id int) (DownloadToolsVO, error) {
	tools, err := service.db.selectById(id)
	if err != nil {
		return DownloadToolsVO{}, errors.New("工具不存在")
	}
	vo := DownloadToolsVO{
		DownloadUrl: fmt.Sprintf("%s/%s", aliyunConfig.baseUrl, tools.OssKey),
		FileName:    fmt.Sprintf("%s.%s", tools.FileName, tools.Extension),
	}
	return vo, nil

}
