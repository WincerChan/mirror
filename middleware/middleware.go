package middleware

import (
	"fmt"
	"mirror/config"
	"net/http"
	"runtime"
)

func RecoverMW() HandlerFunc {
	return func(e *Engine) {
		defer func() {
			//吞掉err
			err := recover()
			if err != nil {
				const size = 64 << 10
				buf := make([]byte, size)
				runtime.Stack(buf, false)
				fmt.Printf("runtime error is %+v, stack is %+v", err, string(buf))
			}

		}()
		e.Next()
	}
}

func IdentityMW() HandlerFunc {
	return func(e *Engine) {
		if config.GetConfig().HeaderTokenKey != "" {
			if e.req.Header.Get(config.GetConfig().HeaderTokenKey) != config.GetConfig().Token {
				e.AbortWithStatus(403)
				return
			} else {
				e.req.Header.Del(config.GetConfig().HeaderTokenKey)
			}

		}
		// todo 支持basic auth
	}
}

func CreateHandler(fun func(rw http.ResponseWriter, req *http.Request)) HandlerFunc {
	return func(e *Engine) {
		fun(e.writer, e.req)
	}
}
