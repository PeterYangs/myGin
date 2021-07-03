package kernel

//middleware/session/session.go

import (
	"myGin/context"
	"myGin/middleware/exception"
	"myGin/middleware/session"
)

// Middleware 全局中间件
var Middleware []context.HandlerFunc

func Load() {

	Middleware = []context.HandlerFunc{
		exception.Exception,
		session.Session,
	}

}
