package entity

import (
	"github.com/go-playground/validator/v10"
)

type AudioDTOInput struct {
	FileName string `json:"filename" validate:"required"`
	Author   string `json:"author" validate:"required"`
	Label    string `json:"label" validate:"required"`
	Type     string `json:"type" validate:"required"`
	Words    string `json:"words" validate:"required"`
}

type AudioInputError struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value"`
}

var Validator = validator.New()

func (a *AudioDTOInput) Validate() []AudioInputError {
	var errors []AudioInputError

	err := Validator.Struct(a)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var el AudioInputError
			el.Field = err.Field()
			el.Tag = err.Tag()
			el.Value = err.Param()
			errors = append(errors, el)
		}
		return errors
	}

	return nil
}
