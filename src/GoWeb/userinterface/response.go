package userinterface

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ResponseEntity struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

const SUCCESS string = "000"
const FAIL string = "001"
const INVALID_USER = "002"

func success() ResponseEntity {
	return ResponseEntity{
		Code:    SUCCESS,
		Message: "操作成功",
	}
}

func fail() ResponseEntity {
	return ResponseEntity{
		Code:    FAIL,
		Message: "操作失败",
	}
}

func invalid() ResponseEntity {
	return ResponseEntity{
		Code:    INVALID_USER,
		Message: "非法用户",
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

func (entity ResponseEntity) write(writer http.ResponseWriter) {
	jsonObj, _ := json.Marshal(entity)
	fmt.Fprintf(writer, string(jsonObj))
}
