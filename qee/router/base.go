package router

import (
	"QeeWeb/qee/base"
)

type Router interface {
	RegisteredHandler(method string, pattern string, handler func(ctx *base.Context)) error
	FindHandler(method string, path string) (func(ctx *base.Context), map[string]string)
}
