package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	UserService UserService
}

func (controller UserController) login(ctx *HttpContext) ResponseEntity {
	param := LoginPram{}
	paramPtr := interface{}(&param)
	ctx.getParam(&paramPtr)
	user, err := controller.UserService.validateUser(param)
	if err != nil {
		responseEntity := message(INVALID_USER, err.Error())
		return responseEntity
	}
	ctx.Principal = user
	role := user.Role
	//生成登录密钥
	token := generateJwtToken(user)
	//获取路由
	var router []string
	routerPtr := interface{}(&router)
	rc.getObjectValue(fmt.Sprintf("router:%s", role), &routerPtr)
	//判断跳转路径
	deviceType, _ := getDevice(*ctx.request)
	var redirectUrl string
	if role == "USER" {
		redirectUrl = "/"
	} else {
		if deviceType == NORMAL_TYPE {
			redirectUrl = "/manager"
		} else {
			redirectUrl = "/"
		}
	}
	vo := LoginVO{
		JwtToken:    token,
		RedirectUrl: redirectUrl,
		Router:      router,
	}
	responseEntity := data(SUCCESS, "登录成功", &vo)
	return responseEntity
}

func (controller UserController) validateUser(ctx *HttpContext) ResponseEntity {
	user := ctx.Principal
	var param ValidateUserParam
	paramPtr := interface{}(&param)
	ctx.getParam(&paramPtr)
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(param.Password))
	result := err != nil
	responseEntity := success("", result)
	return responseEntity
}

func (controller UserController) getUserInformation(ctx *HttpContext) ResponseEntity {
	user := ctx.Principal
	user.Password = ""
	responseEntity := data(SUCCESS, SUCCESS_MESSAGE, user)
	return responseEntity
}

func (controller UserController) createUser(ctx *HttpContext) ResponseEntity {
	param := AddUserParam{}
	paramPtr := interface{}(&param)
	ctx.getParam(&paramPtr)
	err := controller.UserService.addUser(param)
	if err != nil {
		return message(FAIL, err.Error())
	}
	return message(SUCCESS, SUCCESS_MESSAGE)
}

type LinkController struct {
	linkService LinkService
}

func (controller LinkController) getAllLinkList(context *HttpContext) ResponseEntity {
	linkList, err := controller.linkService.getLinkList()
	if err != nil {
		return message(FAIL, "查询失败")
	}
	return success(SUCCESS_MESSAGE, linkList)
}
