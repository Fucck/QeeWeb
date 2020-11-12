package main

import (
	"QeeWeb/qee/base"
	"QeeWeb/qee/router"
	"fmt"
	"net/http"
	"QeeWeb/qee"
)

func indexHandler(writer http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(writer, "URL.PATH = %s", req.URL.Path)
}

func helloHandler(writer http.ResponseWriter, req *http.Request) {
	for k, v := range req.Header {
		fmt.Fprintf(writer, "Header[%q]=%q\n", k, v)
	}
}

func main(){
	router := router.NewTrieRouter()
	handler := qee.New(router)
	handler.RegisteredHandler("GET", "/index/:name", func(ctx *base.Context) {
		fmt.Println(ctx.QueryMap)
	})
	handler.RegisteredHandler("GET", "/*allpath", func(ctx *base.Context) {
		fmt.Println(ctx.QueryMap)
	})
	handler.Run(":8000")
}