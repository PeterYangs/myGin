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

type Session struct {
	Cookie      string                 `json:"cookie"`
	ExpireTime  int64                  `json:"expire_time"`
	SessionList map[string]interface{} `json:"session_list"`
}
