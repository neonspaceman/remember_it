package grpc

import (
	domain_card "card/internal/domain/card"
	"card/internal/grpc/mappers"
	"card/internal/usecase/command"
	api_card "card/gen/go/card/v1"
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *CardImpl) Create(ctx context.Context, req *api_card.CreateRequest) (*api_card.CreateResponse, error) {
	s.log.InfoCtx(
		ctx,
		"Incoming create request",
		zap.String("question", req.Question),
		zap.String("answer", req.Answer),
		zap.String("file_type", req.FileType.String()),
		zap.String("file_id", req.FileId),
	)

	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, fmt.Errorf("parse user id: %w", err)
	}

	cmd := command.CreateCardCmd{
		UserId:   userId,
		Question: req.Question,
		Answer:   req.Answer,
		FileType: domain_card.FileTypeNone,
		FileId:   req.FileId,
	}

	card, err := s.cardCreateHandler.Handle(ctx, cmd)

	if err != nil {
		return nil, fmt.Errorf("create card handler: %w", err)
	}

	return &api_card.CreateResponse{Card: mappers.FromCard(card)}, nil
}
