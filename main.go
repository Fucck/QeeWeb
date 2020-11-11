package main

import (
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
	handler := qee.New()
	handler.GetRegistered("/", func(ctx *qee.Context) {
		ctx.String(http.StatusOK, fmt.Sprintf("URL.PATH = %s", ctx.Path))
	})
	handler.PostRegistered("/hello", func(ctx *qee.Context) {
		name := ctx.PostForm("name")
		if name != "" {
			ctx.Json(200, struct {
				Name string `json:"name"`
				Msg string `json:"msg"`
			}{
				Name: name,
				Msg: "Hello!",
			})
		} else {
			ctx.String(http.StatusBadRequest, "wrong post body")
		}
	})
	handler.Run(":8000")
}