package middleware

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

var onceMiddlewre = &sync.Once{}

type Middleware interface {
	Handler() gin.HandlerFunc
	Recovery() gin.HandlerFunc
	CORS() gin.HandlerFunc
}

type middleware struct {
	log zerolog.Logger
}

func Init(log zerolog.Logger) Middleware {
	var m *middleware

	onceMiddlewre.Do(func() {
		m = &middleware{
			log: log,
		}
	})

	return m
}
