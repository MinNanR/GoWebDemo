package main

import (
	"encoding/json"
	"net/http"
)

type HttpContext struct {
	url            string                            //请求地址
	paramByte      []byte                            //参数的byte数组，可以来自地址参数（get方法），表单或body（Post方法），具体由请求头决定
	Principal      AuthUser                          //当前操作的用户
	request        *http.Request                     //request域
	response       http.ResponseWriter               //response域
	targetFunction func(*HttpContext) ResponseEntity //接口执行的函数
	isAuthorized   bool                              //是否已通过权限校验
	responseEntity ResponseEntity                    //返回值
}

func (context HttpContext) getParam(param *interface{}) {
	json.Unmarshal(context.paramByte, param)
}

/*
中央路由
完成路径转发，转发前解析参数，解析当前使用的用户，打印日志
完成接口动作后打印日志并输出到输出流
*/
func dispatch(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Access-Control-Allow-Origin", "*")                                                                                   //允许访问所有域
	response.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization") //header的类型
	response.Header().Set("Access-Control-Allow-Methods", "POST,GET,OPTIONS")
	url := request.RequestURI
	if request.Method == http.MethodOptions {
		return
	}
	context := HttpContext{
		url:          url,
		request:      request,
		response:     response,
		isAuthorized: false,
	}
	filterChain := createFilterChain()
	filterChain.doFilter(&context)
}

/*
过滤器最后没有映射的函数时的处理方法
*/
func fallBack(ctx *HttpContext) {
	response := ctx.response
	response.WriteHeader(404)
}
