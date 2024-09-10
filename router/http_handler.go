package router

import (
	"github.com/gin-gonic/gin"
)

type appRouter struct {
	method   string
	endpoint string
	worker   gin.HandlerFunc
}

func (h HttpHandler) getRouter() (routes []appRouter) {
	return []appRouter{}
}
