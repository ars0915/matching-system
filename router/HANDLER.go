package router

import (
	"github.com/ars0915/matching-system/config"
	"github.com/ars0915/matching-system/usecase"
)

type HttpHandler struct {
	conf config.ConfENV
	h    usecase.Handler
}

func newHttpHandler(conf config.ConfENV, h usecase.Handler) *HttpHandler {
	return &HttpHandler{
		conf: conf,
		h:    h,
	}
}

func (rH *HttpHandler) Usecase() usecase.Handler {
	return rH.h
}

type Handler struct {
	http *HttpHandler
}

func NewHandler(conf config.ConfENV, h usecase.Handler) Handler {
	return Handler{
		http: newHttpHandler(conf, h),
	}
}
