package main

import (
	"github.com/gin-gonic/gin"
	"myGin/routes"
)

func main() {
	r := gin.Default()

	routes.Load(r)

	r.Run()
}
