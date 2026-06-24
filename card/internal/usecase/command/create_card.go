package command

import (
	domain_card "card/internal/domain/card"
	domain_card_state "card/internal/domain/card_state"
	"context"
	"fmt"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/google/uuid"
)

type CreateCardCmd struct {
	UserId   uuid.UUID
	Question string
	Answer   string
	FileType domain_card.FileType
	FileId   string
}

type CardCreateHandler struct {
	cardRepository      domain_card.CardRepositoryInterface
	cardStateRepository domain_card_state.CardStateRepositoryInterface
	trManager           *manager.Manager
}

func NewCreateCardHandler(
	cardRepository domain_card.CardRepositoryInterface,
	cardStateRepository domain_card_state.CardStateRepositoryInterface,
	trManager *manager.Manager,
) *CardCreateHandler {
	return &CardCreateHandler{
		cardRepository:      cardRepository,
		cardStateRepository: cardStateRepository,
		trManager:           trManager,
	}
}

func (h *CardCreateHandler) Handle(ctx context.Context, cmd CreateCardCmd) (*domain_card.Card, error) {
	card, err := domain_card.NewCard(cmd.UserId, cmd.Answer, cmd.Question, cmd.FileType, cmd.FileId)
	if err != nil {
		return nil, fmt.Errorf("new card: %w", err)
	}

	cardState := domain_card_state.NewCardState(card.Id)

	err = h.trManager.Do(ctx, func(ctx context.Context) error {
		if err := h.cardRepository.Create(ctx, card); err != nil {
			return fmt.Errorf("create card: %w", err)
		}

		if err := h.cardStateRepository.Create(ctx, cardState); err != nil {
			return fmt.Errorf("create card state: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return card, nil
}
