package context

//context/context.go

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"myGin/redis"
	"strings"
	"sync"
	"time"
)

type Context struct {
	*gin.Context
}

type HandlerFunc func(*Context)

func (c *Context) Domain() string {

	return c.Request.Host[:strings.Index(c.Request.Host, ":")]
}

// Session 返回一个session实例
func (c *Context) Session() *Session {

	var session Session

	cookie, ok := c.Get("_session")

	if !ok {

		return nil
	}

	session = cookie.(Session)

	session.Lock = &sync.Mutex{}

	return &session

}

type Session struct {
	Cookie      string                 `json:"cookie"`
	ExpireTime  int64                  `json:"expire_time"`
	SessionList map[string]interface{} `json:"session_list"`
	Lock        *sync.Mutex
}

func (s *Session) Set(key string, value interface{}) error {

	//加锁，防止读取并行
	s.Lock.Lock()

	defer s.Lock.Unlock()

	sessionString, err := redis.Client().Get(context.TODO(), s.Cookie).Result()

	if err != nil {

		return err
	}

	var session Session

	err = json.Unmarshal([]byte(sessionString), &session)

	if err != nil {

		return err
	}

	//设置
	session.SessionList[key] = value

	sessionStringNew, err := json.Marshal(session)

	//计算新的过期时间
	e := s.ExpireTime - time.Now().Unix()

	if e < 0 {

		return errors.New("the session has expired")
	}

	redis.Client().Set(context.TODO(), s.Cookie, sessionStringNew, time.Duration(e)*time.Second)

	return nil
}

func (s *Session) Get(key string) (interface{}, error) {

	sessionString, err := redis.Client().Get(context.TODO(), s.Cookie).Result()

	if err != nil {

		return nil, err
	}

	var session Session

	err = json.Unmarshal([]byte(sessionString), &session)

	if err != nil {

		return nil, err
	}

	value, ok := session.SessionList[key]

	if ok {

		return value, nil
	}

	return nil, errors.New("not found key :" + key)

}

func (s *Session) Remove(key string) error {

	s.Lock.Lock()

	defer s.Lock.Unlock()

	sessionString, err := redis.Client().Get(context.TODO(), s.Cookie).Result()

	if err != nil {

		return err
	}

	var session Session

	err = json.Unmarshal([]byte(sessionString), &session)

	if err != nil {

		return err
	}

	delete(session.SessionList, key)

	sessionStringNew, err := json.Marshal(session)

	if err != nil {

		return err
	}

	//计算新的过期时间
	e := s.ExpireTime - time.Now().Unix()

	if e < 0 {

		return errors.New("the session has expired")
	}

	redis.Client().Set(context.TODO(), s.Cookie, sessionStringNew, time.Duration(e)*time.Second)

	return nil

}
