package context

//context/context.go

import (
	"github.com/gin-gonic/gin"
	"strings"
)

type Context struct {
	*gin.Context
}

type HandlerFunc func(*Context)

func (c *Context) Domain() string {

	return c.Request.Host[:strings.Index(c.Request.Host, ":")]
}
