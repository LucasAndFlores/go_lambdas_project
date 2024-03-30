package dto

import (
	"github.com/go-playground/validator/v10"
)

type MetadataDTOInput struct {
	FileName string `json:"filename" validate:"required"`
	Author   string `json:"author" validate:"required"`
	Label    string `json:"label" validate:"required"`
	Type     string `json:"type" validate:"required"`
	Words    string `json:"words" validate:"required"`
}

type MetadataDTOOutput struct {
	FileName string `json:"filename" validate:"required"`
	Author   string `json:"author" validate:"required"`
	Label    string `json:"label" validate:"required"`
	Type     string `json:"type" validate:"required"`
	Words    string `json:"words" validate:"required"`
}

type MetadataInputError struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value"`
}

var MetadataValidator = validator.New()

func (a *MetadataDTOInput) Validate() []MetadataInputError {
	var errors []MetadataInputError

	err := MetadataValidator.Struct(a)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var el MetadataInputError
			el.Field = err.Field()
			el.Tag = err.Tag()
			el.Value = err.Param()
			errors = append(errors, el)
		}
		return errors
	}

	return nil
}
