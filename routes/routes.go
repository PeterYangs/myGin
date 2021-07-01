package routes

//routes/routes.go

import (
	"github.com/gin-gonic/gin"
	"myGin/controller"
	"myGin/response"
)

func Load(r *gin.Engine) {

	r.GET("/", convert(controller.Index))

}

func convert(f func(*gin.Context) *response.Response) gin.HandlerFunc {

	return func(c *gin.Context) {

		resp := f(c)

		data := resp.GetData()

		switch item := data.(type) {

		case string:

			c.String(200, item)

		case gin.H:

			c.JSON(200, item)

		}

	}

}
