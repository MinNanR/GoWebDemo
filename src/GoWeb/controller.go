package main

import "log"

func login(context HttpContext) ResponseEntity {
	param := LoginPram{}
	p := interface{}(&param)
	context.getParam(&p)
	user, err := validateUser(param)
	if err != nil {
		responseEntity := message(INVALID_USER, err.Error())
		return responseEntity
	}
	token := generateJwtToken(user)
	deviceType, _ := getDevice(*context.request)
	responseEntity := data(SUCCESS, "登录成功", map[string]interface{}{"jwtToken": token, "deviceType": deviceType})
	return responseEntity
}

func getAllLinkList(context HttpContext) ResponseEntity {
	log.Println("当前用户", context.user.Username)
	return ResponseEntity{}
}
