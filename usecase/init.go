package usecase

import (
	"github.com/ars0915/matching-system/internal/tree"
)

func InitHandler(boysTree, girlsTree *tree.PersonTree) Handler {

	h := newHandler()

	return h
}
