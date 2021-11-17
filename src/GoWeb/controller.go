package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	UserService UserService
	SignService SignService
}

func (controller UserController) login(ctx *HttpContext) ResponseInterface {
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

func (controller UserController) validateUser(ctx *HttpContext) ResponseInterface {
	user := ctx.Principal
	var param ValidateUserParam
	paramPtr := interface{}(&param)
	ctx.getParam(&paramPtr)
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(param.Password))
	result := err == nil
	responseEntity := success("", result)
	return responseEntity
}

func (controller UserController) getUserInformation(ctx *HttpContext) ResponseInterface {
	user := ctx.Principal
	user.Password = ""
	responseEntity := data(SUCCESS, SUCCESS_MESSAGE, user)
	return responseEntity
}

func (controller UserController) getUserInformationList(ctx *HttpContext) ResponseInterface {
	userList, err := controller.UserService.getUserInformationList()
	if err != nil {
		return message(FAIL, err.Error())
	}
	return data(SUCCESS, SUCCESS_MESSAGE, userList)
}

func (controller UserController) createUser(ctx *HttpContext) ResponseInterface {
	param := AddUserParam{}
	paramPtr := interface{}(&param)
	ctx.getParam(&paramPtr)
	err := controller.UserService.addUser(param)
	if err != nil {
		return message(FAIL, err.Error())
	}
	return message(SUCCESS, SUCCESS_MESSAGE)
}

func (controller UserController) updateUser(ctx *HttpContext) ResponseInterface {
	param := UpdateUserParam{}
	paramStr := interface{}(&param)
	ctx.getParam(&paramStr)
	err := controller.UserService.updateUser(param)
	if err != nil {
		return message(FAIL, err.Error())
	}
	return message(SUCCESS, SUCCESS_MESSAGE)
}

func (controller UserController) deleteUser(ctx *HttpContext) ResponseInterface {
	param := DeleteUserParam{}
	paramStr := interface{}(&param)
	ctx.getParam(&paramStr)
	err := controller.UserService.deleteUser(param.Id)
	if err != nil {
		return message(FAIL, err.Error())
	}
	return message(SUCCESS, SUCCESS_MESSAGE)
}

type LinkController struct {
	linkService LinkService
}

func (controller LinkController) getAllLinkList(context *HttpContext) ResponseInterface {
	linkList, err := controller.linkService.getLinkList()
	if err != nil {
		return message(FAIL, "查询失败")
	}
	return success(SUCCESS_MESSAGE, linkList)
}

func (controller LinkController) addLink(ctx *HttpContext) ResponseInterface {
	param := AddLinkParam{}
	paramStr := interface{}(&param)
	ctx.getParam(&paramStr)
	if len(param.Link) == 0 || len(param.Name) == 0 {
		return fail(INVALID_PARAM, "缺少参数")
	}
	err := controller.linkService.addLink(param)
	if err != nil {
		return message(FAIL, err.Error())
	}
	return message(SUCCESS, SUCCESS_MESSAGE)
}

func (controller LinkController) deleteLink(ctx *HttpContext) ResponseInterface {
	param := DeleteLinkParam{}
	paramStr := interface{}(&param)
	ctx.getParam(&paramStr)
	if param.Id == 0 {
		return message(INVALID_PARAM, "缺少参数")
	}
	err := controller.linkService.deleteLink(param.Id)
	if err != nil {
		return message(FAIL, err.Error())
	}
	return message(SUCCESS, SUCCESS_MESSAGE)
}

func (controller LinkController) updateLink(ctx *HttpContext) ResponseInterface {
	param := UpdateLinkParam{}
	paramStr := interface{}(&param)
	ctx.getParam(&paramStr)
	if param.Id == 0 {
		return message(INVALID_PARAM, "缺少参数")
	}
	err := controller.linkService.updateLink(param)
	if err != nil {
		return message(FAIL, err.Error())
	}
	return message(SUCCESS, SUCCESS_MESSAGE)
}

func (controller UserController) generateSign(ctx *HttpContext) ResponseInterface {
	user := ctx.Principal
	param := CreateSignDTO{}
	paramStr := interface{}(&param)
	ctx.getParam(&paramStr)
	param.UserId = user.Id
	url, err := controller.SignService.createSubscribeSign(param)
	if err != nil {
		return message(FAIL, err.Error())
	}
	return success(SUCCESS_MESSAGE, url)
}

func (controller LinkController) getSubscribe(ctx *HttpContext) ResponseInterface {
	param := GetLinkDTO{}
	paramStr := interface{}(&param)
	ctx.getParam(&paramStr)
	subscribe, err := controller.linkService.getSubscribe(param)
	if err != nil {
		return PlainResponse{Data: ""}
	}
	return PlainResponse{Data: subscribe}
}

type IntroductionController struct {
	service IntroductionService
}

func (controller IntroductionController) updateIntroduction(ctx *HttpContext) ResponseInterface {
	param := UpdateIntroductionParam{}
	paramStr := interface{}(&param)
	ctx.getParam(&paramStr)
	controller.service.updateIntroduction(param.Content)
	return message(SUCCESS, SUCCESS_MESSAGE)
}

func (controller IntroductionController) getIntroduction(ctx *HttpContext) ResponseInterface {
	content := controller.service.getIntroduction()
	if len(content) == 0 {
		content = "敬请期待"
	}
	result := map[string]string{"introduction": content}

	return success(SUCCESS_MESSAGE, result)
}

func (controller IntroductionController) getImageList(ctx *HttpContext) ResponseInterface {
	imageList, err := controller.service.getImageList()
	if err != nil {
		return message(FAIL, err.Error())
	}
	return success(SUCCESS_MESSAGE, imageList)
}

type ToolsController struct {
	service ToolsService
}

func (controller ToolsController) getToolsList(ctx *HttpContext) ResponseInterface {
	toolsList, err := controller.service.getToolsList()
	if err != nil {
		return message(FAIL, err.Error())
	}
	return success(SUCCESS_MESSAGE, toolsList)
}

func (controller ToolsController) downloadTools(ctx *HttpContext) ResponseInterface {
	param := DownloadToolsDTO{}
	paramStr := interface{}(&param)
	ctx.getParam(&paramStr)
	vo, err := controller.service.downloadTools(param.Id)
	if err != nil {
		return message(FAIL, FAIL_MESSAGE)
	}
	return success(SUCCESS, vo)
}

type FileController struct {
	toolsService        ToolsService
	introductionService IntroductionService
}

func (controller FileController) insertImage(ctx *HttpContext) ResponseInterface {
	request := ctx.request
	fileForm := request.MultipartForm.File
	imageFile := fileForm["image"][0]
	dto := AddImageDTO{Image: imageFile}
	return controller.introductionService.addImage(dto)
}

func (controller IntroductionController) deleteImage(ctx *HttpContext) ResponseInterface {
	param := DeleteImageDTO{}
	paramStr := interface{}(&param)
	ctx.getParam(&paramStr)
	err := controller.service.deleteImage(param)
	if err != nil {
		return fail(err.Error(), "")
	}
	return success(SUCCESS_MESSAGE, "")
}
