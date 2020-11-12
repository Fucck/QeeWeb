package qee

import (
	"QeeWeb/qee/base"
	"QeeWeb/qee/router"
	"fmt"
	"net/http"
)

type qeeWeb struct {
	router.Router
}

func New(r router.Router) *qeeWeb {
	return &qeeWeb{
		r,
	}
}

func (qee *qeeWeb) GetRegistered(rule string, f func(ctx *base.Context)) error {
	err := qee.RegisteredHandler(http.MethodGet, rule, f)
	if err != nil {
		return err
	}
	return nil
}

func (qee *qeeWeb) PostRegistered(rule string, f func(ctx *base.Context)) error {
	err := qee.RegisteredHandler(http.MethodPost, rule, f)
	if err != nil {
		return err
	}
	return nil
}

func BadRequestHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "404 request: %s", req.URL.Path)
}

// generate context from req, and pass to all process
func (qee *qeeWeb) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	routerKey := req.Method + "-" + req.URL.Path
	fmt.Printf("router key: %s\n", routerKey)
	f, queryMap := qee.FindHandler(req.Method, req.URL.Path)
	if f == nil {
		BadRequestHandler(w, req)
	} else {
		ctx := base.NewContext(w, req)
		ctx.QueryMap = queryMap
		fmt.Printf("call function\n")
		f(ctx)
	}

}

func (qee *qeeWeb) Run(bindAddr string) error {
	return http.ListenAndServe(bindAddr, qee)
}
