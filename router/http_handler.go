package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type appRouter struct {
	method   string
	endpoint string
	worker   gin.HandlerFunc
}

func (h HttpHandler) getRouter() (routes []appRouter) {
	return []appRouter{
		{http.MethodPost, "/addPersonAndFindMatch/", h.addPersonAndFindMatchHandler},
		{http.MethodGet, "/print/", h.printHandler},
	}
}
