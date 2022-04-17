// go 上传文件
package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

func postFile(filename string, targetUrl string) error {
	bodyBuf := &bytes.Buffer{}

	bodyWriter := multipart.NewWriter(bodyBuf)

	// 关键的一步操作
	fileWriter, err := bodyWriter.CreateFormFile("uploadFile", filename)

	if err != nil {
		fmt.Println("error writting to buffer")
		return err
	}
	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		return err
	}
	defer fh.Close()

	// io copy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(targetUrl, contentType, bodyBuf)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(resp.Status)
	fmt.Println(string(respBody))
	return nil
}

func main() {
	err := postFile("./butter-fly.jpg", "http://localhost:9090/upload")
	if err != nil {
		fmt.Println("upload error", err)
	}
}
