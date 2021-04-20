package main

import (
	"log"
	"net/http"
)

func main() {
	port := server.Port
	buildRoutine()
	buildFilter()
	log.Printf("server started, listening on %s\n", port)
	http.HandleFunc("/", dispatch)
	err := http.ListenAndServe(":"+port, nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
