package service

import (
	"net/http"
	"../global"
	"strings"
	"fmt"
	"time"
)

func GetIndex(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.RequestURI, ".html") {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"),
			r.RemoteAddr, "访问", r.RequestURI)
	}
	global.RootFileServer.ServeHTTP(w, r)
}
