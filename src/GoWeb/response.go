package main

import (
	"encoding/json"
	"fmt"
	"io"
)

type ResponseEntity struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

const SUCCESS string = "000"
const FAIL string = "001"
const INVALID_USER = "002"
const INVALID_PARAM = "005"
const UNKNOW_ERROR = "500"

const SUCCESS_MESSAGE = "操作成功"
const FAIL_MESSAGE = "操作失败"
const INVALID_USER_MESSAGE = "非法用户"
const UNKNOW_ERROR_MESSAGE = "未知错误"

func success(message string, data interface{}) ResponseEntity {
	return ResponseEntity{
		Code:    SUCCESS,
		Message: message,
		Data:    data,
	}
}

func fail(message string, data interface{}) ResponseEntity {
	return ResponseEntity{
		Code:    FAIL,
		Message: message,
		Data:    data,
	}
}

func invalidUser() ResponseEntity {
	return ResponseEntity{
		Code:    INVALID_USER,
		Message: INVALID_USER_MESSAGE,
	}
}

func message(code string, message string) ResponseEntity {
	return ResponseEntity{
		Code:    code,
		Message: message,
	}
}

func data(code string, message string, data interface{}) ResponseEntity {
	return ResponseEntity{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func defaultData(data interface{}) ResponseEntity {
	return ResponseEntity{
		Code:    SUCCESS,
		Message: "操作成功",
		Data:    data,
	}
}

func (entity ResponseEntity) write(writer io.Writer) {
	jsonObj, _ := json.Marshal(entity)
	fmt.Fprintf(writer, string(jsonObj))
}
