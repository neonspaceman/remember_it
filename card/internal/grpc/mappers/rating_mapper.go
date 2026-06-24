package mappers

import (
	review_domain "card/internal/domain/review"
	api_card "card/gen/go/card/v1"
	"errors"
	"fmt"
)

var (
	ErrUnknownApiRating    = errors.New("unknown rating type")
	ErrUnknownDomainRating = errors.New("unknown rating representation")
)

func ToRating(rating api_card.Rating) (review_domain.RatingType, error) {
	switch rating {
	case api_card.Rating_RATING_AGAIN:
		return review_domain.AgainRating, nil
	case api_card.Rating_RATING_HARD:
		return review_domain.HardRating, nil
	case api_card.Rating_RATING_GOOD:
		return review_domain.GoodRating, nil
	case api_card.Rating_RATING_EASY:
		return review_domain.EasyRating, nil
	default:
		return "", fmt.Errorf("rating type '%d': %w", rating, ErrUnknownApiRating)
	}
}

func FromRating(rating review_domain.RatingType) (api_card.Rating, error) {
	switch rating {
	case review_domain.AgainRating:
		return api_card.Rating_RATING_AGAIN, nil
	case review_domain.HardRating:
		return api_card.Rating_RATING_HARD, nil
	case review_domain.GoodRating:
		return api_card.Rating_RATING_GOOD, nil
	case review_domain.EasyRating:
		return api_card.Rating_RATING_EASY, nil
	default:
		return api_card.Rating(0), fmt.Errorf("rating type '%s': %w", rating, ErrUnknownDomainRating)
	}
}
