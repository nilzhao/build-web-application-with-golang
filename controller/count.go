package controller

import (
	"html/template"
	"net/http"
	"time"

	"github.com/nilzhao/build-web-application-with-golang/session"
)

func Count(w http.ResponseWriter, r *http.Request) {
	sess := session.GlobalSessions.SessionStart(w, r)
	createTime := sess.Get("createTime")
	if createTime == nil {
		sess.Set("createTime", time.Now().Unix())
	} else if (createTime.(int64) + 360) < time.Now().Unix() {
		session.GlobalSessions.SessionDestroy(w, r)
		sess = session.GlobalSessions.SessionStart(w, r)
	}
	ct := sess.Get("countnum")
	if ct == nil {
		sess.Set("countnum", 1)
	} else {
		sess.Set("countnum", ct.(int)+1)
	}
	t, _ := template.ParseFiles("template/count.html")
	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, sess.Get("countnum"))
}
