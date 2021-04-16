package userinterface

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func getBody(request *http.Request) []byte {
	bodyBuffer := make([]byte, request.ContentLength)
	request.Body.Read(bodyBuffer)
	return bodyBuffer
}

func Login(response http.ResponseWriter, request *http.Request) {
	type LoginPram struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	bodyByte := getBody(request)
	var param LoginPram
	json.Unmarshal(bodyByte, &param)
	fmt.Println(param.Username, param.Password)
	responseEntity := message(SUCCESS, "登录成功")
	responseEntity.write(response)
}
