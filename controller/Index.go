package controller

//controller/Index.go

import (
	"myGin/context"
	"myGin/response"
)

func Index(context *context.Context) *response.Response {

	//panic("假装错误")

	return response.Resp().String(context.Domain())
}
