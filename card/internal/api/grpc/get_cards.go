package grpc

import (
	"card/internal/grpc/mappers"
	"card/internal/usecase/query"
	api_card "card/gen/go/card/v1"
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *CardImpl) GetByUserId(ctx context.Context, req *api_card.GetByUserIdRequest) (*api_card.GetByUserIdResponse, error) {
	s.log.InfoCtx(
		ctx,
		"Incoming get by user id request",
		zap.String("user_id", req.UserId),
	)

	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, fmt.Errorf("parse user id: %w", err)
	}

	afterCursor := uuid.Nil

	if req.After != "" {
		afterCursor, err = uuid.Parse(req.After)
		if err != nil {
			return nil, fmt.Errorf("parse cursor: %w", err)
		}
	}

	cmd := query.GetCardsByUserIdQuery{
		UserId: userId,
		Limit:  req.Limit + 1,
		After:  afterCursor,
	}

	cards, err := s.getCardsByUserIdQuery.Handle(ctx, cmd)

	if err != nil {
		return nil, fmt.Errorf("get cards query: %w", err)
	}

	var encCursor uuid.UUID
	res := make([]*api_card.Card, 0, req.Limit)
	limit := min(uint64(len(cards)), req.Limit)

	for _, card := range cards[:limit] {
		res = append(res, mappers.FromCard(card))
		encCursor = card.Id
	}

	return &api_card.GetByUserIdResponse{
		Cards:     res,
		HasNext:   uint64(len(cards)) > req.Limit,
		EndCursor: encCursor.String(),
	}, nil
}
