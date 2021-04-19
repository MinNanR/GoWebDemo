package main

import (
	"log"
	"net/http"
)

/**
注册路由表，所有接口路径都在这里加入到http路由中
*/
//func buildUpRouter() {
//	http.HandleFunc("/login", Login)
//}

func main() {
	buildRoutine()
	log.Println("server started, listening on 8989")
	http.HandleFunc("/", routine)
	err := http.ListenAndServe(":8989", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	//user, err := GetUserByUsername("Minnan")
	//if err != nil{
	//	fmt.Println(err)
	//	fmt.Println("user not found")
	//}
	//fmt.Println(user)
}
