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

func (rH HttpHandler) getRouter() (routes []appRouter) {
	return []appRouter{
		{http.MethodPost, "/addPersonAndFindMatch/", rH.addPersonAndFindMatchHandler},
		{http.MethodGet, "/print/", rH.printHandler},
		{http.MethodDelete, "/removeSinglePerson/:id/", rH.removePersonHandler},
	}
}
