package card

import (
	card_api "card/gen/go/card/v1"
	"fmt"
)

func ToCard(card *card_api.Card) *Card {
	return &Card{Id: card.Id}
}

func ToReviewLog(reviewLog *card_api.ReviewLog) *ReviewLog {
	return &ReviewLog{Id: reviewLog.Id}
}

func FromRating(rating string) (card_api.Rating, error) {
	switch rating {
	case "AGAIN":
		return card_api.Rating_RATING_AGAIN, nil
	case "HARD":
		return card_api.Rating_RATING_HARD, nil
	case "GOOD":
		return card_api.Rating_RATING_GOOD, nil
	case "EASY":
		return card_api.Rating_RATING_EASY, nil
	default:
		return 0, fmt.Errorf("parse '%s': %w", rating, ErrUnknownRating)
	}
}
