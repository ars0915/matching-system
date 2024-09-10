package usecase

import (
	"github.com/ars0915/matching-system/internal/tree"
)

type AppHandler struct {
	Person
}

type NewHandlerOption func(*AppHandler)

func newHandler(optFn ...NewHandlerOption) *AppHandler {
	h := &AppHandler{}

	for _, o := range optFn {
		o(h)
	}

	return h
}

type PersonHandler struct {
	boys  *tree.PersonTree
	girls *tree.PersonTree
}

func NewPersonHandler(boysTree, girlsTree *tree.PersonTree) *PersonHandler {
	return &PersonHandler{
		boys:  boysTree,
		girls: girlsTree,
	}
}

func WithPerson(i PersonHandler) func(h *AppHandler) {
	return func(h *AppHandler) {
		h.Person = i
	}
}
