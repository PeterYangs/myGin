package routes

//routes/routes.go

import (
	"github.com/gin-gonic/gin"
	"myGin/controller"
)

func Load(r *gin.Engine) {

	router := newRouter(r)

	router.Group("/api", func(api group) {

		api.Group("/user", func(user group) {

			user.Registered(GET, "/info", controller.Index)
			user.Registered(GET, "/order", controller.Index)
			user.Registered(GET, "/money", controller.Index)

		})

	})

}
