package routes

////routes/router.go

import (
	"github.com/gin-gonic/gin"
	"myGin/context"
	"myGin/response"
)

type router struct {
	engine *gin.Engine
}

type group struct {
	engine      *gin.Engine
	path        string
	middlewares []context.HandlerFunc
}

type method int

const (
	GET    method = 0x000000
	POST   method = 0x000001
	PUT    method = 0x000002
	DELETE method = 0x000003
	ANY    method = 0x000004
)

func newRouter(engine *gin.Engine) *router {

	return &router{
		engine: engine,
	}
}

func (r *router) Group(path string, callback func(group), middlewares ...context.HandlerFunc) {

	callback(group{
		engine:      r.engine,
		path:        path,
		middlewares: middlewares,
	})

}

func (g group) Group(path string, callback func(group), middlewares ...context.HandlerFunc) {

	//需要注意，父级的中间件在前面
	g.middlewares = append(g.middlewares, middlewares...)

	g.path += path

	callback(g)
}

func (g group) Registered(method method, url string, action func(*context.Context) *response.Response, middlewares ...context.HandlerFunc) {

	//父类中间件+当前action中间件+action
	var handlers = make([]gin.HandlerFunc, len(g.middlewares)+len(middlewares)+1)

	//添加当前action的中间件
	g.middlewares = append(g.middlewares, middlewares...)

	//将中间件转换为gin.HandlerFunc
	for key, middleware := range g.middlewares {

		temp := middleware

		handlers[key] = func(c *gin.Context) {

			temp(&context.Context{Context: c})
		}
	}

	//添加action
	handlers[len(g.middlewares)] = convert(action)

	finalUrl := g.path + url

	switch method {

	case GET:

		g.engine.GET(finalUrl, handlers...)

	case POST:

		g.engine.GET(finalUrl, handlers...)

	case PUT:

		g.engine.PUT(finalUrl, handlers...)

	case DELETE:

		g.engine.DELETE(finalUrl, handlers...)

	case ANY:

		g.engine.Any(finalUrl, handlers...)

	}

}

func convert(f func(*context.Context) *response.Response) gin.HandlerFunc {

	return func(c *gin.Context) {

		resp := f(&context.Context{Context: c})

		data := resp.GetData()

		switch item := data.(type) {

		case string:

			c.String(200, item)

		case gin.H:

			c.JSON(200, item)

		}

	}

}
