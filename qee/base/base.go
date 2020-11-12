package base

import (
	"encoding/json"
	"net/http"
)

type Context struct {
	Writer     http.ResponseWriter
	Req        *http.Request
	Path       string
	Method     string
	StatusCode int
	QueryMap   map[string]string
}

// call by ServeHttp only
func NewContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer:     w,
		Req:        req,
		Path:       req.URL.Path,
		Method:     req.Method,
		StatusCode: http.StatusNotFound,
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
