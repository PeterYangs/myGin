package kernel

import (
	"myGin/context"
	"myGin/middleware/exception"
)

// Middleware 全局中间件
var Middleware []context.HandlerFunc

func Load() {

	Middleware = []context.HandlerFunc{
		exception.Exception,
	}

}
