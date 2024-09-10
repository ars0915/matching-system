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

	ErrorMatchSameGender = cGin.CustomError{
		Code:     1001,
		HTTPCode: http.StatusBadRequest,
		Message:  "Match same gender",
	}

	ErrorWantedDateLimit = cGin.CustomError{
		Code:     1001,
		HTTPCode: http.StatusBadRequest,
		Message:  "Wanted date limit",
	}

	ErrorHeightCheckFailed = cGin.CustomError{
		Code:     1001,
		HTTPCode: http.StatusBadRequest,
		Message:  "Height check failed",
	}
)
