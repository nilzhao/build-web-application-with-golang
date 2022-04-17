package controller

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) // 获取请求的方法
	if r.Method == "GET" {
		t, _ := template.ParseFiles("template/login.html")
		log.Println(t.Execute(w, nil))
	} else {
		username := r.FormValue("username")

		fmt.Println("username:", username)
		fmt.Println("password:", r.FormValue("password"))
		w.Header().Set("Content-Type", "text/html")
		// XSS 漏洞 <script>alert(1)</script> 会执行
		fmt.Fprintf(w, "welcome back "+username)
		// 转义
		fmt.Fprintf(w, "welcome back "+template.HTMLEscapeString(username))
	}
}
