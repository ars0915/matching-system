package entity

import "github.com/ars0915/matching-system/constant"

type Person struct {
	ID          uint64
	Name        string
	Height      float64
	Gender      constant.Gender
	WantedDates int64
}
