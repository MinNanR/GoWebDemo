package main

import (
	"GoWeb/userinterface"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/login", userinterface.Login)
	log.Println("server started, listening on 8989")
	err := http.ListenAndServe(":8989", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
