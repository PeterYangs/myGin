package controller

//controller/Index.go

import (
	"golang.org/x/time/rate"
	"myGin/component/limiter"
	"myGin/context"
	"myGin/response"
	"strconv"
	"time"
)

func Index(context *context.Context) *response.Response {

	l := limiter.NewLimiter(rate.Every(1*time.Second), 1, context.ClientIP())

	if !l.Allow() {

		return response.Resp().String("您的访问过于频繁")
	}

	return response.Resp().String("success")
}

func Index2(context *context.Context) *response.Response {

	//msg, _ := context.Session().Get("msg")

	//fmt.Println(limiter.GlobalLimiters)

	return response.Resp().String("nice")
}

func Index3(context *context.Context) *response.Response {

	context.Session().Remove("msg")

	return response.Resp().String("")
}

func Index4(context *context.Context) *response.Response {

	session := context.Session()

	for i := 0; i < 100; i++ {

		go func(index int) {

			session.Set("msg"+strconv.Itoa(index), index)

		}(i)
	}

	return response.Resp().String("")
}
