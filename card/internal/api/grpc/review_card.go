package grpc

import (
	card_domain "card/internal/domain/card"
	review_domain "card/internal/domain/review"
	"card/internal/grpc/mappers"
	"card/internal/usecase/command"
	api_card "card/gen/go/card/v1"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *CardImpl) Review(ctx context.Context, req *api_card.ReviewRequest) (*api_card.ReviewResponse, error) {
	s.log.InfoCtx(
		ctx,
		"Incoming create request",
		zap.String("user_id", req.UserId),
		zap.String("card_id", req.CardId),
		zap.String("rating", req.Rating.String()),
		zap.String("review_at", req.ReviewAt.String()),
	)

	cardId, err := uuid.Parse(req.CardId)
	if err != nil {
		return nil, fmt.Errorf("parse card id: %w", err)
	}

	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, fmt.Errorf("parse user id: %w", err)
	}

	rating, err := mappers.ToRating(req.Rating)
	if err != nil {
		return nil, fmt.Errorf("parse rating: %w", err)
	}

	reviewLog, err := s.reviewCardHandler.Handle(ctx, command.ReviewCardCmd{
		UserId:     userId,
		CardId:     cardId,
		ReviewedAt: req.ReviewAt.AsTime(),
		Rating:     rating,
	})
	if err != nil {
		if v, ok := errors.AsType[*card_domain.CardNotFoundError](err); ok {
			return nil, NewErrCardNotFound(v.CardId)
		}
		if v, ok := errors.AsType[*review_domain.ReviewPeriodNotStartError](err); ok {
			return nil, NewErrReviewPeriodNotStart(v.Due)
		}

		return nil, fmt.Errorf("review: %w", err)
	}

	ratingResponse, err := mappers.FromRating(reviewLog.Rating)
	if err != nil {
		return nil, fmt.Errorf("parse rating: %w", err)
	}

	return &api_card.ReviewResponse{
		ReviewLog: &api_card.ReviewLog{
			Id:     reviewLog.Id.String(),
			CardId: reviewLog.CardId.String(),
			Rating: ratingResponse,
		},
	}, nil
}
