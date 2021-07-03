package session

import (
	uuid "github.com/satori/go.uuid"
	"myGin/context"
)

var cookieName = "my_gin"

func Session(c *context.Context) {

	sessionKey := uuid.NewV4().String()

	c.SetCookie(cookieName, sessionKey, 3600, "/", c.Domain(), false, true)

}
