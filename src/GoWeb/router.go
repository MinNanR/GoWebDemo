package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func getBody(request *http.Request) []byte {
	bodyBuffer := make([]byte, request.ContentLength)
	request.Body.Read(bodyBuffer)
	return bodyBuffer
}

type HttpContext struct {
	param   []byte
	user    AuthUser
	request *http.Request
}

func (context HttpContext) getParam(param *interface{}) {
	json.Unmarshal(context.param, param)
}

var router map[string]func(HttpContext) ResponseEntity
var passingUrl []string
var re, _ = regexp.Compile(`\s`)

/*
构建路由
*/
func buildRoutine() {
	router = make(map[string]func(HttpContext) ResponseEntity)
	router["/login"] = login
	router["/api/getAllLinkList"] = getAllLinkList
	passingUrl = []string{"/login"}
}

/*
中央路由
完成路径转发，转发前解析参数，解析当前使用的用户，打印日志
完成接口动作后打印日志并输出到输出流
*/
func routine(response http.ResponseWriter, request *http.Request) {
	url := request.RequestURI
	if handler, exist := router[url]; exist {
		context := HttpContext{
			param:   getBody(request),
			request: request,
		}
		startTime := time.Now().UnixNano() / 1e6
		log.Printf("处理接口地址：%s,请求参数%s", request.RequestURI, re.ReplaceAllString(string(context.param), ""))
		var responseEntity ResponseEntity
		user, err := parseUser(request)
		if err != nil {
			flag := false
			for _, passing := range passingUrl {
				if url == passing {
					flag = true
					break
				}
			}
			if flag {
				responseEntity = handler(context)
			} else {
				responseEntity = message(INVALID_USER, err.Error())
			}
		} else {
			context.user = user
			responseEntity = handler(context)
		}

		finishTime := time.Now().UnixNano() / 1e6
		jsonResponse, _ := json.Marshal(responseEntity)
		log.Printf("处理接口%s完成，返回参数%s,耗时%dms\n", request.RequestURI, string(jsonResponse), finishTime-startTime)
		response.Write(jsonResponse)
	} else {
		response.WriteHeader(404)
	}
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
