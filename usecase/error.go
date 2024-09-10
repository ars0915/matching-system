package usecase

import (
	"net/http"

	"github.com/ars0915/matching-system/util/cGin"
)

var (
	ErrorPersonNotFound = cGin.CustomError{
		Code:     1001,
		HTTPCode: http.StatusNotFound,
		Message:  "Person not found",
	}

	ErrorPersonExist = cGin.CustomError{
		Code:     1001,
		HTTPCode: http.StatusNotFound,
		Message:  "Person not found",
	}
)
