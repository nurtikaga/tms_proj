package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"

	"github.com/nurtikaga/tms_proj/internal/adapters/cache"
	grpcadapter "github.com/nurtikaga/tms_proj/internal/adapters/grpc"
	kafkaadapter "github.com/nurtikaga/tms_proj/internal/adapters/kafka"
	"github.com/nurtikaga/tms_proj/internal/adapters/logger"
	postrepo "github.com/nurtikaga/tms_proj/internal/adapters/repository/postgres"
	"github.com/nurtikaga/tms_proj/internal/config"
	"github.com/nurtikaga/tms_proj/internal/core/usecase"
	tmsmanagerpb "github.com/nurtikaga/tms_proj/tms-protos/go"
)

func Run(cfg *config.Config) error {
	log := logger.New(cfg.Log.Level)

	pool, err := pgxpool.New(context.Background(), cfg.Postgres.DSN)
	if err != nil {
		return fmt.Errorf("postgres connect: %w", err)
	}
	defer pool.Close()

	if err = pool.Ping(context.Background()); err != nil {
		return fmt.Errorf("postgres ping: %w", err)
	}
	log.Info("postgres connected", "dsn", cfg.Postgres.DSN)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	defer redisClient.Close()

	if err = redisClient.Ping(context.Background()).Err(); err != nil {
		return fmt.Errorf("redis ping: %w", err)
	}
	log.Info("redis connected", "addr", cfg.Redis.Addr)

	kafkaProducer := kafkaadapter.New(cfg.Kafka.BrokerAddr, cfg.Kafka.Topic)
	defer kafkaProducer.Close()
	log.Info("kafka producer ready", "broker", cfg.Kafka.BrokerAddr, "topic", cfg.Kafka.Topic)

	shipmentRepo := postrepo.NewShipmentRepo(pool)
	eventRepo := postrepo.NewEventRepo(pool)
	shipmentCache := cache.New(redisClient)

	svc := usecase.New(shipmentRepo, eventRepo, kafkaProducer, shipmentCache, log)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcadapter.RecoveryInterceptor(log),
			grpcadapter.LoggingInterceptor(log),
		),
	)
	tmsmanagerpb.RegisterShipmentServiceServer(grpcServer, grpcadapter.NewHandler(svc))

	addr := fmt.Sprintf(":%d", cfg.GRPC.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen %s: %w", addr, err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info("gRPC server started", "addr", addr)
		if serr := grpcServer.Serve(lis); serr != nil {
			slog.Error("gRPC serve error", "error", serr)
		}
	}()

	<-stop
	log.Info("shutting down...")
	grpcServer.GracefulStop()
	log.Info("server stopped")
	return nil
}
