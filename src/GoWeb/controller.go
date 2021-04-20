package main

import "fmt"

func login(ctx *HttpContext) ResponseEntity {
	param := LoginPram{}
	paramPtr := interface{}(&param)
	ctx.getParam(&paramPtr)
	user, err := validateUser(param)
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

func getAllLinkList(context *HttpContext) ResponseEntity {
	linkList, err := getLinkList()
	if err != nil {
		return message(FAIL, "查询失败")
	}
	return success(SUCCESS_MESSAGE, linkList)
}
