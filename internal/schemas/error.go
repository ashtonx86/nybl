package schemas

import "github.com/go-playground/validator"

type BadRequestError struct {
	Field string
	Tag   string
	Value string
}

// IMPORTANT: Only pass errors from validator
func NewBadRequestError(err error) []*BadRequestError {
	var errors []*BadRequestError
	
	for _, err := range err.(validator.ValidationErrors) {
		var badReqErr BadRequestError
		badReqErr.Field = err.Field()
		badReqErr.Tag = err.Tag()
		badReqErr.Value = err.Param()
		errors = append(errors, &badReqErr)
	}
	return errors
}