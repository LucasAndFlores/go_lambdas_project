package entity

import "github.com/LucasAndFlores/go_lambdas_project/internal/dto"

type Metadata struct {
	FileName string `dynamodbav:"filename"`
	Author   string `dynamodbav:"author"`
	Label    string `dynamodbav:"label"`
	Type     string `dynamodbav:"type"`
	Words    string `dynamodbav:"words"`
}

func (m *Metadata) ConvertToDTO() dto.MetadataDTOInput {
	return dto.MetadataDTOInput{
		FileName: m.FileName,
		Author:   m.Author,
		Label:    m.Label,
		Type:     m.Type,
		Words:    m.Words,
	}

}
