package qee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Context struct {
	Writer http.ResponseWriter
	Req *http.Request
	Path string
	Method string
	StatusCode int
}

type qeeWeb struct {
	router map[string] func(ctx *Context)
}

func New() *qeeWeb {
	return &qeeWeb{
		router: map[string]func(ctx *Context){},
	}
}

func (ctx *Context) PostForm(key string) string {
	return ctx.Req.PostFormValue(key)
}

func (ctx *Context) Status(code int) {
	ctx.StatusCode = code
	ctx.Writer.WriteHeader(code)
}

func (ctx *Context) SetHeader(key, val string) {
	ctx.Writer.Header().Set(key, val)
}

func (ctx *Context) Html(code int, html string) {
	ctx.SetHeader("Content-Type", "text/html")
	ctx.Status(code)
	ctx.Writer.Write([]byte(html))
}

func (ctx *Context) String(code int, str string) {
	ctx.SetHeader("Content-Type", "text/plain")
	ctx.Status(code)
	ctx.Writer.Write([]byte(str))
}

func (ctx *Context) Data(code int, data []byte) {
	ctx.Status(code)
	ctx.Writer.Write(data)
}

func (ctx *Context) Json(code int, obj interface{}) {
	ctx.SetHeader("Content-Type", "text/json")
	ctx.Status(code)
	encoder := json.NewEncoder(ctx.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusInternalServerError)
	}
}

// call by ServeHttp only
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req: req,
		Path: req.URL.Path,
		Method: req.Method,
		StatusCode: http.StatusNotFound,
	}
}

func (qee *qeeWeb) GetRegistered(rule string, f func(ctx *Context)){
	routerKey := http.MethodGet + "-" + rule
	if _, exist := qee.router[routerKey]; exist {
		//TODO: warning log
	}
	qee.router[routerKey] = f
	fmt.Printf("register %s\n", routerKey)
}

func (qee *qeeWeb) PostRegistered(rule string, f func(ctx *Context)) {
	routerKey := http.MethodPost + "-" + rule
	if _, exist := qee.router[routerKey]; exist {
		//TODO: warning log
	}
	qee.router[routerKey] = f
	fmt.Printf("register %s\n", routerKey)
}

func BadRequestHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "404 request: %s", req.URL.Path)
}

// generate context from req, and pass to all process
func (qee *qeeWeb) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	routerKey := req.Method + "-" + req.URL.Path
	fmt.Printf("router key: %s\n", routerKey)
	if f, exist := qee.router[routerKey]; exist {
		ctx := newContext(w, req)
		fmt.Printf("call function\n")
		f(ctx)
	} else {
		BadRequestHandler(w, req)
	}
}

func (qee *qeeWeb) Run(bindAddr string) error {
	return http.ListenAndServe(bindAddr, qee)
}
