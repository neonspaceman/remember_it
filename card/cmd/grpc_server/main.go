package main

import (
	"buf.build/go/protovalidate"
	api_card "card/gen/go/card/v1"
	"card/internal/adapter/postgresql"
	grpc_api "card/internal/api/grpc"
	"card/internal/config"
	"card/internal/consts"
	"card/internal/grpc/interceptors"
	review_service "card/internal/service/review"
	"card/internal/usecase/command"
	"card/internal/usecase/query"
	"context"
	"flag"
	"fmt"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_logging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpc_protovalidate "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"os/signal"
	log_pkg "platform/pkg/logger"
	"slices"
	"syscall"
	"time"
)

const defaultEnvFile = ".env"

func main() {
	var overrideEnv string
	flag.StringVar(&overrideEnv, "env", "", "path to override .env file")
	flag.Parse()

	envs := []string{defaultEnvFile}

	if overrideEnv != "" {
		envs = slices.Insert(envs, 0, overrideEnv)
	}

	cfg, errors := config.Load(envs...)

	if errors != nil {
		panic(errors)
	}

	errors = run(cfg)

	if errors != nil {
		panic(errors)
	}
}

func run(cfg *config.Config) error {
	log, err := log_pkg.NewLogger(
		log_pkg.WithDevelopmentMode(cfg.App.Env == consts.EnvDev),
		log_pkg.WithLevel(cfg.Log.Level),
	)
	if err != nil {
		return err
	}
	defer log.Close()

	pool, err := pgxpool.New(context.Background(), cfg.Database.DSN)
	if err != nil {
		return err
	}
	err = pool.Ping(context.Background())
	if err != nil {
		return err
	}
	defer pool.Close()

	protoValidator, err := protovalidate.New()
	if err != nil {
		return fmt.Errorf("failed to init proto validator: %w", err)
	}

	trManager, err := manager.New(trmpgx.NewDefaultFactory(pool))
	trGetter := trmpgx.DefaultCtxGetter

	if err != nil {
		return err
	}

	cardRepository := postgresql.NewCardRepository(pool, trGetter)
	cardStateRepository := postgresql.NewCardStateRepository(pool, trGetter)
	reviewLogRepository := postgresql.NewReviewLogRepository(pool, trGetter)

	createCardHandler := command.NewCreateCardHandler(cardRepository, cardStateRepository, trManager)
	getCardsByUserIdHandler := query.NewGetCardByUserIdHandler(pool)
	reviewCardHandler := command.NewReviewCardHandler(cardRepository, cardStateRepository, reviewLogRepository, review_service.NewScheduler(), trManager)

	api := grpc_api.NewCardImpl(grpc_api.CardImplProps{
		CreateCardHandler:     createCardHandler,
		GetCardsByUserIdQuery: getCardsByUserIdHandler,
		ReviewCardHandler:     reviewCardHandler,
		Log:                   log,
	})

	// GRPC Server
	l, err := net.Listen("tcp", cfg.GRPC.Addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	defer l.Close()

	server := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: time.Duration(cfg.GRPC.MaxConnectionIdle) * time.Minute,
			Timeout:           time.Duration(cfg.GRPC.Timeout) * time.Second,
			MaxConnectionAge:  time.Duration(cfg.GRPC.MaxConnectionAge) * time.Minute,
			Time:              time.Duration(cfg.GRPC.Timeout) * time.Minute,
		}),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
			//InterceptorRequestId(),
			grpc_logging.UnaryServerInterceptor(interceptors.LoggerInterceptor(log), grpc_logging.WithLogOnEvents(grpc_logging.StartCall, grpc_logging.FinishCall)),
			grpc_protovalidate.UnaryServerInterceptor(protoValidator),
		)),
	)

	api_card.RegisterCardServiceServer(server, api)

	if cfg.App.Env == consts.EnvDev {
		reflection.Register(server)
	}

	g, _ := errgroup.WithContext(context.Background())

	g.Go(func() error {
		log.Info("GRPC server is listening...", zap.String("Addr", cfg.GRPC.Addr))
		err := server.Serve(l)

		if err != nil {
			return err
		}

		log.Info("GRPC server shut down correctly")

		return err
	})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	server.GracefulStop()

	return g.Wait()
}
