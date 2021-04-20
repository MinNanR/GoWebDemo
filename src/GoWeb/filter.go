package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type FilterChain struct {
	chain      []func(*HttpContext, *FilterChain)
	chainIndex int
}

func (filterChain *FilterChain) doFilter(ctx *HttpContext) {
	chainLength := len(filterChain.chain)
	filterChain.chainIndex++
	if filterChain.chainIndex == chainLength {
		if ctx.targetFunction == nil {
			fallBack(ctx)
		} else {
			responseEntity := ctx.targetFunction(ctx)
			ctx.responseEntity = responseEntity
		}
	} else {
		filterChain.chain[filterChain.chainIndex](ctx, filterChain)
	}
}

func getBody(request *http.Request) []byte {
	bodyBuffer := make([]byte, request.ContentLength)
	request.Body.Read(bodyBuffer)
	return bodyBuffer
}

var router map[string]func(*HttpContext) ResponseEntity
var re, _ = regexp.Compile(`\s`)

/*
构建路由
*/
func buildRoutine() {
	router = make(map[string]func(*HttpContext) ResponseEntity)
	router["/login"] = login
	router["/api/getAllLinkList"] = getAllLinkList
	for key, _ := range router {
		log.Printf("Mapped url %s\n", key)
	}
}

func buildFilter() {
	filterList = append(filterList, UrlFilter)
	filterList = append(filterList, ParamFilter)
	filterList = append(filterList, UserFilter)
	filterList = append(filterList, AuthorityFilter)
	filterList = append(filterList, LogFilter)
	filterList = append(filterList, ResponseFilter)
}

/*
路径过滤器，负责寻找路径对应的处理函数
*/
func UrlFilter(ctx *HttpContext, filterChain *FilterChain) {
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
func ParamFilter(ctx *HttpContext, filterChain *FilterChain) {
	request := ctx.request
	requestMethod := request.Method
	contentType := request.Header.Get("Content-Type")
	for key, value := range request.Header {
		fmt.Printf("%s:%s\n", key, value)
	}
	if strings.ToUpper(requestMethod) == "POST" && contentType == "application/json" {
		ctx.paramByte = getBody(request)
	} else {
		err := request.ParseForm()
		if err != nil {
			form := request.Form
			ctx.paramByte, _ = json.Marshal(form)
		}
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
func UserFilter(ctx *HttpContext, filterChain *FilterChain) {
	isAuthorityPath := isAuthorityPath(ctx.url)
	if !isAuthorityPath {
		user, err := parseUser(ctx.request)
		if err != nil {
			response := ctx.response
			responseEntity := message(INVALID_USER, err.Error())
			responseEntity.export(response)
			return
		} else {
			ctx.Principal = user
		}
	} else {
		ctx.isAuthorized = true
	}
	filterChain.doFilter(ctx)
}

/*
权限校验，后续补充
*/
func AuthorityFilter(ctx *HttpContext, filterChain *FilterChain) {
	filterChain.doFilter(ctx)
}

/*
日志过滤器
*/
func LogFilter(ctx *HttpContext, filterChain *FilterChain) {
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
func ResponseFilter(ctx *HttpContext, filterChain *FilterChain) {
	filterChain.doFilter(ctx)
	accept := ctx.request.Header.Get("Accept")
	responseEntity := ctx.responseEntity
	response := ctx.response
	if len(accept) == 0 {
		response.Header().Add("Content-Type", "application/json")
	} else {
		response.Header().Add("Content-Type", accept)
	}
	responseEntity.export(response)
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
	err := validateJwtToken(strings.Replace(tokenString, "Bearer ", "", 1), &claims)
	if err != nil {
		return AuthUser{}, err
	}
	user, err := loadUserByUserName(claims.Subject)
	if err != nil {
		return AuthUser{}, err
	}
	return user, nil
}