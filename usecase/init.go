package usecase

import (
	"github.com/ars0915/matching-system/internal/tree"
)

func InitHandler(boysTree, girlsTree *tree.PersonTree) Handler {
	person := NewPersonHandler(boysTree, girlsTree)
	h := newHandler(
		WithPerson(person),
	)

	return h
}
