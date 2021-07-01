package controller

//controller/Index.go

import (
	"github.com/gin-gonic/gin"
	"myGin/response"
)

func Index(context *gin.Context) *response.Response {

	if 1+1 == 2 {

		return response.Resp().Json(gin.H{"msg": "hello world"})

	}

	return response.Resp().Json(gin.H{"msg": "hello world"})
}
