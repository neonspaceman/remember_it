package command

import (
	card_domain "card/internal/domain/card"
	card_state_domain "card/internal/domain/card_state"
	review_domain "card/internal/domain/review"
	"card/internal/service/review"
	"context"
	"fmt"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/google/uuid"
	"time"
)

type ReviewCardCmd struct {
	UserId     uuid.UUID
	CardId     uuid.UUID
	ReviewedAt time.Time
	Rating     review_domain.RatingType
}

type ReviewCardHandler struct {
	cardRepository      card_domain.CardRepositoryInterface
	cardStateRepository card_state_domain.CardStateRepositoryInterface
	reviewLogRepository review_domain.ReviewLogRepositoryInterface
	scheduler           review.SchedulerInterface
	trManager           *manager.Manager
}

func NewReviewCardHandler(
	cardRepository card_domain.CardRepositoryInterface,
	cardStateRepository card_state_domain.CardStateRepositoryInterface,
	reviewLogRepository review_domain.ReviewLogRepositoryInterface,
	scheduler review.SchedulerInterface,
	trManager *manager.Manager,
) *ReviewCardHandler {
	return &ReviewCardHandler{
		cardRepository:      cardRepository,
		cardStateRepository: cardStateRepository,
		reviewLogRepository: reviewLogRepository,
		scheduler:           scheduler,
		trManager:           trManager,
	}
}

func (h *ReviewCardHandler) Handle(ctx context.Context, cmd ReviewCardCmd) (*review_domain.ReviewLog, error) {
	card, err := h.cardRepository.GetById(ctx, cmd.CardId)
	if err != nil {
		return nil, fmt.Errorf("get card '%s': %w", cmd.CardId, err)
	}

	if card.UserId != cmd.UserId {
		return nil, card_domain.NewCardNotFoundError(cmd.CardId)
	}

	cardState, err := h.cardStateRepository.GetById(ctx, cmd.CardId)
	if err != nil {
		return nil, fmt.Errorf("get card state '%s': %w", cmd.CardId, err)
	}

	if !cmd.ReviewedAt.After(cardState.Due) {
		return nil, review_domain.NewReviewPeriodNotStartError(cmd.ReviewedAt, cardState.Due)
	}

	reviewLog, err := h.scheduler.ReviewCard(cardState, cmd.Rating, cmd.ReviewedAt)
	if err != nil {
		return nil, fmt.Errorf("review card '%s': %w", cmd.CardId, err)
	}

	err = h.trManager.Do(ctx, func(ctx context.Context) error {
		if err := h.cardStateRepository.Save(ctx, cardState); err != nil {
			return fmt.Errorf("save card state '%s': %w", cmd.CardId, err)
		}

		if err := h.reviewLogRepository.Create(ctx, reviewLog); err != nil {
			return fmt.Errorf("create review log '%s': %w", cmd.CardId, err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return reviewLog, nil
}
