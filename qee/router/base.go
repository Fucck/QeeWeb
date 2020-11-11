package router

import "QeeWeb/qee"

type Router interface {
	RegisteredHandler(pattern string, handler func(ctx *qee.Context))
	FindHandler(path string) func(ctx *qee.Context)
}
