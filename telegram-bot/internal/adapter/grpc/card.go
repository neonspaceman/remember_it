package grpc

import (
	card_api "card/gen/go/card/v1"
	"context"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/metadata"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
	"platform/pkg/request_id"
	card_client "telegram-bot/internal/client/card"
	"telegram-bot/internal/config"
	"telegram-bot/internal/consts"
)

type CardClient struct {
	grpc card_api.CardServiceClient
}

func NewCardClient(cfg *config.Config) *CardClient {
	conn, err := grpc.NewClient(
		cfg.GRPC.CardURI,
		//grpc.WithChainUnaryInterceptor(InterceptorRequestId()),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		panic(err)
	}

	grpcClient := card_api.NewCardServiceClient(conn)

	return &CardClient{grpc: grpcClient}
}

func (c *CardClient) Create(ctx context.Context, userId uuid.UUID, question, answer string) (*card_client.Card, error) {
	response, err := c.grpc.Create(ctx, &card_api.CreateRequest{
		UserId:   userId.String(),
		Question: question,
		Answer:   answer,
	})

	if err != nil {
		return nil, err
	}

	return card_client.ToCard(response.Card), nil
}

func (c *CardClient) GetCards(ctx context.Context, userId uuid.UUID, limit uint64, after string) ([]*card_client.Card, bool, string, error) {
	response, err := c.grpc.GetByUserId(ctx, &card_api.GetByUserIdRequest{
		UserId: userId.String(),
		Limit:  limit,
		After:  after,
	})

	if err != nil {
		return nil, false, "", err
	}

	cards := make([]*card_client.Card, 0, len(response.Cards))

	for _, v := range response.Cards {
		cards = append(cards, card_client.ToCard(v))
	}

	return cards, response.HasNext, response.EndCursor, nil
}

func (c *CardClient) ReviewCard(ctx context.Context, userId uuid.UUID, cardId uuid.UUID, rating string) (*card_client.ReviewLog, error) {
	parsedRating, err := card_client.FromRating(rating)
	if err != nil {
		return nil, err
	}

	response, err := c.grpc.Review(ctx, &card_api.ReviewRequest{
		UserId:   userId.String(),
		CardId:   cardId.String(),
		ReviewAt: timestamppb.Now(),
		Rating:   parsedRating,
	})

	if err != nil {
		if _, ok := AsReason(err, "ERR_CARD_NOT_FOUND"); ok {
			return nil, card_client.ErrCardNotFound
		}

		if _, ok := AsReason(err, "ERR_REVIEW_PERIOD_NOT_START"); ok {
			return nil, card_client.ErrReviewPeriodNotStart
		}

		return nil, err
	}

	return card_client.ToReviewLog(response.ReviewLog), nil
}

// @todo: to platform??
func InterceptorRequestId() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		requestId := request_id.CtxGet(ctx)

		if requestId != "" {
			requestId += ","
		}

		requestId += uuid.NewString()

		ctx = metadata.ExtractOutgoing(ctx).Set(consts.GrpcRequestIdKey, requestId).ToOutgoing(ctx)

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
