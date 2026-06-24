package mappers

import (
	domain_card "card/internal/domain/card"
	api_card "card/gen/go/card/v1"
)

func FromCard(e *domain_card.Card) *api_card.Card {
	return &api_card.Card{
		Id:       e.Id.String(),
		Question: e.Question,
		Answer:   e.Answer,
		FileType: fromFileType(e.FileType),
		FileId:   e.FileId,
	}
}

func fromFileType(fileType domain_card.FileType) api_card.FileType {
	switch fileType {
	case domain_card.FileTypePhoto:
		return api_card.FileType_FILE_TYPE_PHOTO
	case domain_card.FileTypeDocument:
		return api_card.FileType_FILE_TYPE_DOCUMENT
	default:
		return api_card.FileType_FILE_TYPE_UNSPECIFIED
	}
}
