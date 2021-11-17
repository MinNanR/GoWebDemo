package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type Filter struct {
	filterMethod func(*HttpContext, *Filter)
	next         *Filter
}

//type FilterChain struct {
//	chain      []func(*HttpContext, *FilterChain)
//	chainIndex int
//}
type FilterChain Filter

func (filter *Filter) doFilter(ctx *HttpContext) {
	if filter == nil {
		if ctx.targetFunction == nil {
			fallBack(ctx)
		} else {
			responseEntity := ctx.targetFunction(ctx)
			ctx.responseEntity = responseEntity
		}
	} else {
		filter.filterMethod(ctx, filter.next)
	}
}

func (filter *Filter) addFilterAt(nextFilter *Filter) {
	nextFilter.next = filter.next
	filter.next = nextFilter
}

func getBody(request *http.Request) []byte {
	bodyBuffer := make([]byte, request.ContentLength)
	request.Body.Read(bodyBuffer)
	return bodyBuffer
}

var router map[string]func(*HttpContext) ResponseInterface
var authUrls map[string]*AuthUrls
var re, _ = regexp.Compile(`\s`)

type AuthUrls struct {
	Urls []string
}

func (auth AuthUrls) isAuthorized(targetUrl string) bool {
	for _, url := range auth.Urls {
		if url == targetUrl {
			return true
		}
	}
	return false
}

/*
构建路由
*/
func buildRoutine() {
	router = make(map[string]func(*HttpContext) ResponseInterface)
	userController := UserController{}
	linkController := LinkController{}
	introductionController := IntroductionController{}
	toolsController := ToolsController{}
	fileController := FileController{}
	router["/login"] = userController.login
	router["/api/getUserInformation"] = userController.getUserInformation
	router["/manager/addUser"] = userController.createUser
	router["/manager/getUserList"] = userController.getUserInformationList
	router["/manager/updateUser"] = userController.updateUser
	router["/manager/validateUser"] = userController.validateUser
	router["/manager/deleteUser"] = userController.deleteUser
	router["/api/getAllLinkList"] = linkController.getAllLinkList
	router["/manager/addLink"] = linkController.addLink
	router["/manager/deleteLink"] = linkController.deleteLink
	router["/manager/updateLink"] = linkController.updateLink
	router["/api/getIntroduction"] = introductionController.getIntroduction
	router["/manager/updateIntroduction"] = introductionController.updateIntroduction
	router["/manager/getImageList"] = introductionController.getImageList
	router["/api/getToolsList"] = toolsController.getToolsList
	router["/api/downloadTools"] = toolsController.downloadTools
	router["/subscribe"] = linkController.getSubscribe
	router["/api/generateSign"] = userController.generateSign
	router["/manager/addImage"] = fileController.insertImage
	router["/manager/deleteImage"] = introductionController.deleteImage
	for key, _ := range router {
		log.Printf("Mapped url %s\n", key)
	}

	adminUrl := []string{"/api/getUserInformation", "/api/getAllLinkList", "/api/getIntroduction",
		"/api/getToolsList", "/api/downloadTools", "/manager/addUser", "/manager/getUserList", "/manager/updateUser",
		"/manager/validateUser", "/manager/deleteUser", "/manager/addLink", "/manager/deleteLink",
		"/manager/updateLink", "/manager/updateIntroduction", "/manager/getImageList", "/subscribe",
		"/api/generateSign", "/manager/addImage", "/manager/deleteImage"}
	userUrl := []string{"/api/getUserInformation", "/api/getAllLinkList", "/api/getIntroduction",
		"/api/getToolsList", "/api/downloadTools", "/subscribe", "/api/generateSign"}

	authUrls = make(map[string]*AuthUrls)
	authUrls["ADMIN"] = &AuthUrls{Urls: adminUrl}
	authUrls["USER"] = &AuthUrls{Urls: userUrl}

}

func buildFilter() {
	filterMethodList = append(filterMethodList, UrlFilter)
	filterMethodList = append(filterMethodList, ParamFilter)
	filterMethodList = append(filterMethodList, UserFilter)
	filterMethodList = append(filterMethodList, SubscribeFilter)
	filterMethodList = append(filterMethodList, AuthorityFilter)
	filterMethodList = append(filterMethodList, LogFilter)
	filterMethodList = append(filterMethodList, ResponseFilter)
}

/*
路径过滤器，负责寻找路径对应的处理函数
*/
func UrlFilter(ctx *HttpContext, filterChain *Filter) {
	if len(router) == 0 {
		buildRoutine()
	}
	url := ctx.url
	//找到路径映射的函数
	if handler, exist := router[url]; exist {
		ctx.targetFunction = handler //设置处理请求的函数
	} else {
		fallBack(ctx)
		return
	}
	filterChain.doFilter(ctx)
}

/*
参数解析过滤器
*/
func ParamFilter(ctx *HttpContext, filterChain *Filter) {
	request := ctx.request
	requestMethod := strings.ToUpper(request.Method)
	contentType := request.Header.Get("Content-Type")
	if requestMethod == "POST" && contentType == "application/json" {
		ctx.paramByte = getBody(request)
	} else if requestMethod == "POST" && contentType == "application/x-www-form-urlencoded" {
		err := request.ParseForm()
		if err != nil {
			form := request.Form
			ctx.paramByte, _ = json.Marshal(form)
		}
	} else if requestMethod == "POST" && strings.Contains(contentType, "multipart/form-data") {
		request.ParseMultipartForm(64 << 20)
	} else if requestMethod == "GET" {
		queryList := request.URL.Query()
		query := make(map[string]string)
		for k, v := range queryList {
			query[k] = v[0]
		}
		ctx.paramByte, _ = json.Marshal(query)
	}
	filterChain.doFilter(ctx)
}

func isAuthorityPath(url string) bool {
	for _, path := range authorityPath {
		if url == path {
			return true
		}
	}
	return false
}

/*
用户解析过滤器
*/
func UserFilter(ctx *HttpContext, filterChain *Filter) {
	if !ctx.isAuthorized {
		isAuthorityPath := isAuthorityPath(ctx.url)
		if !isAuthorityPath {
			user, err := parseUser(ctx.request)
			if err != nil {
				//response := ctx.response
				//responseEntity := message(INVALID_USER, err.Error())
				//responseEntity.export(response)
				ctx.response.WriteHeader(401)
				return
			} else {
				ctx.Principal = user
			}
		} else {
			ctx.isAuthorized = true
		}
		filterChain.doFilter(ctx)
	} else {
		filterChain.doFilter(ctx)
	}
}

/*
权限校验，后续补充
*/
func AuthorityFilter(ctx *HttpContext, filterChain *Filter) {
	if !ctx.isAuthorized {
		isAuthorityPath := isAuthorityPath(ctx.url)
		if !isAuthorityPath {
			principal := ctx.Principal
			authorizedUrls := authUrls[principal.Role]
			isAuthorized := authorizedUrls.isAuthorized(ctx.url)
			if !isAuthorized {
				ctx.response.WriteHeader(401)
				return
			}
		}
	}
	filterChain.doFilter(ctx)
}

/*
日志过滤器
*/
func LogFilter(ctx *HttpContext, filterChain *Filter) {
	request := ctx.request
	startTime := time.Now().UnixNano() / 1e6
	log.Printf("处理接口地址：%s,请求参数[%s]", request.RequestURI, re.ReplaceAllString(string(ctx.paramByte), ""))
	filterChain.doFilter(ctx)
	endTime := time.Now().UnixNano() / 1e6
	responseEntity := ctx.responseEntity
	jsonResponse, _ := json.Marshal(responseEntity)
	log.Printf("处理接口%s完成，返回参数%s,耗时%dms\n", request.RequestURI, string(jsonResponse), endTime-startTime)
}

/*
返回值过滤器
*/
func ResponseFilter(ctx *HttpContext, filterChain *Filter) {
	filterChain.doFilter(ctx)
	accept := ctx.request.Header.Get("Accept")
	responseEntity := ctx.responseEntity
	response := ctx.response
	if accept == "application/json" {
		response.Header().Add("Content-Type", "application/json")
	} else if len(accept) > 0 {
		response.Header().Add("Content-Type", accept)
	} else {
		response.Header().Add("Content-type", "application/text")
	}
	responseEntity.export(response)
}

func SubscribeFilter(ctx *HttpContext, filterChain *Filter) {
	if ctx.url == "/subscribe" {
		user, err := parseSubscribeUser(ctx)
		if err != nil {
			ctx.response.WriteHeader(401)
			return
		}
		ctx.Principal = user
		ctx.isAuthorized = true
	}
	filterChain.doFilter(ctx)
}

/*
解析用户信息
*/
func parseUser(request *http.Request) (AuthUser, error) {
	tokenString := request.Header.Get("Authorization")
	if !strings.HasPrefix(tokenString, "Bearer ") {
		log.Println("JWT Token does not begin with bearer String")
		return AuthUser{}, JwtError{message: "无法识别用户信息"}
	}
	claims := CustomClaims{}
	service := UserService{}
	err := validateJwtToken(strings.Replace(tokenString, "Bearer ", "", 1), &claims)
	if err != nil {
		return AuthUser{}, err
	}
	user, err := service.loadUserByUserName(claims.Subject)
	if err != nil {
		return AuthUser{}, err
	}
	return user, nil
}

func parseSubscribeUser(ctx *HttpContext) (AuthUser, error) {
	param := &SubscribeDTO{}
	paramStr := interface{}(&param)
	ctx.getParam(&paramStr)
	service := UserService{}
	user, err := service.loadUserBySign(param.Sign)
	if err != nil {
		return AuthUser{}, err
	}
	return user, nil
}
