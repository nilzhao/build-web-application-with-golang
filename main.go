package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nilzhao/build-web-application-with-golang/controller"
)

const port = "9090"

func main() {
	http.HandleFunc("/", controller.SayHello)
	http.HandleFunc("/login", controller.Login)
	http.HandleFunc("/upload", controller.Upload)
	fmt.Println("server will run at", "http://localhost:"+port)
	err := http.ListenAndServe(":"+port, nil) // 设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
