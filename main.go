package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nilzhao/build-web-application-with-golang/controller"
	"github.com/nilzhao/build-web-application-with-golang/session"
	_ "github.com/nilzhao/build-web-application-with-golang/session/providers/memory"
)

const port = "9090"

func main() {
	session.GlobalSessions, _ = session.NewManager("memory", "go-session-id", 3600)
	go session.GlobalSessions.GC()

	http.HandleFunc("/", controller.SayHello)
	http.HandleFunc("/login", controller.Login)
	http.HandleFunc("/upload", controller.Upload)
	http.HandleFunc("/count", controller.Count)
	fmt.Println("server will run at", "http://localhost:"+port)
	err := http.ListenAndServe(":"+port, nil) // 设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
