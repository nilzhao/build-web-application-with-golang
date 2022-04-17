package controller

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

// 接收文件
func Upload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()

		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))
		t, _ := template.ParseFiles("template/upload.html")
		t.Execute(w, token)
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadFile")
		fmt.Println("file", file)
		if err != nil {
			fmt.Println("get file err:", err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		f, openErr := os.OpenFile("./file-tmp/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if openErr != nil {
			fmt.Println(err)
		}
		defer f.Close()
		io.Copy(f, file)
	}
}
