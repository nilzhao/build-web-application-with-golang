package controller

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/nilzhao/build-web-application-with-golang/session"
)

func setCookie(w http.ResponseWriter, username string) {
	expiration := time.Now().AddDate(1, 0, 0)
	cookie := http.Cookie{
		Name:     "username",
		Value:    username,
		Expires:  expiration,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}

func LoginCookie(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Method == "GET" {
		t, _ := template.ParseFiles("template/login.html")
		// 清除 cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "username",
			Value:    "",
			HttpOnly: true,
		})
		log.Println(t.Execute(w, nil))
	} else {
		username := r.FormValue("username")

		fmt.Println("username:", username)
		fmt.Println("password:", r.FormValue("password"))
		w.Header().Set("Content-Type", "text/html")
		setCookie(w, username)
		// XSS 漏洞 <script>alert(1)</script> 会执行
		fmt.Fprintf(w, "welcome back "+username)
		// 转义
		fmt.Fprintf(w, "welcome back "+template.HTMLEscapeString(username))
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	ses := session.GlobalSessions.SessionStart(w, r)
	fmt.Println(ses)
	r.ParseForm()
	if r.Method == "GET" {
		t, _ := template.ParseFiles("template/login.html")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		t.Execute(w, ses.Get("username"))
	} else {
		ses.Set("username", r.FormValue("username"))
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
