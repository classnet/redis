//index.go
package main

import (
	"./model"
	"net/http"
)

func main() {

	http.HandleFunc("/reg", model.Reg) //设置访问的路由
	http.HandleFunc("/login", model.Login)
	http.ListenAndServe(":85", nil) //设置监听的端口
}

