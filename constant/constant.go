package constant

type Gender string

const (
	ServiceName        = "matching-system"
	ResponseCodePrefix = 1

	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
)
