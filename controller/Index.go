package controller

//controller/Index.go

import (
	"myGin/context"
	"myGin/response"
)

func Index(context *context.Context) *response.Response {

	return response.Resp().String(context.Domain())
}
