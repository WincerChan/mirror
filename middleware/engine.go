package middleware

import (
	"net/http"
	"sync"
)

type Engine struct {
	handlers HandlersChain
	index    int
	req      *http.Request
	writer   http.ResponseWriter
}

type HandlersChain []HandlerFunc
type HandlerFunc func(*Engine)

func (e *Engine) Use(middleware ...HandlerFunc) {
	e.handlers = append(e.handlers, middleware...)
	return
}

func (e *Engine) Run() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		e.req = r
		e.writer = w
		e.index = -1
		e.Next()
	}
}

func (e *Engine) Next() {
	e.index++
	for e.index < len(e.handlers) {
		e.handlers[e.index](e)
		e.index++
	}
}

func (e *Engine) Abort() {
	e.index = len(e.handlers)
}
func (e *Engine) Status(code int) {
	e.writer.WriteHeader(code)
}
func (e *Engine) AbortWithStatus(code int) {
	e.Status(code)
	e.Abort()
}

var (
	once sync.Once
	eng  *Engine
)

func GetEngine(middleware ...HandlerFunc) *Engine {
	if eng != nil {
		return eng
	}
	once.Do(func() {
		eng = &Engine{}
		eng.Use(middleware...)
	})
	return eng
}
