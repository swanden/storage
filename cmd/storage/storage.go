package main

import (
	"context"
	"fmt"
	"github.com/swanden/storage/internal/adapters"
	storageController "github.com/swanden/storage/internal/controller/grpc/storage"
	storageUseCase "github.com/swanden/storage/internal/usecase/storage"
	"github.com/swanden/storage/pkg/cache"
	"github.com/swanden/storage/pkg/interceptors"
	"github.com/swanden/storage/pkg/logger"
	"github.com/swanden/storage/pkg/memcached"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpcStorage "github.com/swanden/storage/api/go"
)

const (
	storageTypeMemcached = "memcached"
	storageTypeCache     = "cache"
)

type Storage interface {
	Set(ctx context.Context, key, value string, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	Close()
}

func main() {
	cfg := newConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		os.Exit(0)
	}()

	loggerInst, err := logger.New(
		ctx,
		logger.WithLevel(cfg.LogLevel),
	)
	if err != nil {
		panic("Unable to create logger instance")
	}

	var storage Storage
	if cfg.StorageType == storageTypeMemcached {
		memcachedClient, err := memcached.Connect(
			cfg.MemcachedHost,
			memcached.WithPort(cfg.MemcachedPort),
			memcached.WithMaxIdleConns(cfg.MemcachedMaxIdleConns),
			memcached.WithMaxOpenConns(cfg.MemcachedMaxOpenConns),
			memcached.WithNewConnTimeout(cfg.MemcachedNewConnTimeout),
			memcached.WithConnRetryTimeout(cfg.MemcachedConnRetryTimeout),
		)
		if err != nil {
			loggerInst.Error().Err(err).Msg("Unable to create memcached client")
		}

		memcachedAdapter := adapters.NewMemcachedAdapter(memcachedClient)
		storage = memcachedAdapter
	} else {
		cacheInst := cache.New()
		cacheAdapter := adapters.NewCacheAdapter(cacheInst)
		storage = cacheAdapter
	}
	defer storage.Close()

	storageUseCaseInst, err := storageUseCase.New(
		storageUseCase.WithLogger(loggerInst),
		storageUseCase.WithStorage(storage),
	)

	storageControllerInst, err := storageController.New(
		storageController.WithLogger(loggerInst),
		storageController.WithStorageUseCase(storageUseCaseInst),
	)

	grpcServerOptions := []grpc.ServerOption{
		grpc.UnaryInterceptor(
			interceptors.UnaryRequestUUIDGenerator(),
		),
		grpc.MaxRecvMsgSize(cfg.GRPCServerMaxRcvSize),
		grpc.ConnectionTimeout(cfg.GRPCServerTimeOutConnection),
		grpc.NumStreamWorkers(uint32(cfg.HandlerWorkerPoolSize)),
	}
	serverGRPC := grpc.NewServer(grpcServerOptions...)

	grpcStorage.RegisterStorageServer(serverGRPC, storageControllerInst)

	serverGRPCLis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCServerListenerPort))
	if err != nil {
		loggerInst.Error().Err(err).Msg("Unable to listen TCP for GRPC server")
	}

	go func() {
		loggerInst.Info().Msg("GRPC server start")

		if err := serverGRPC.Serve(serverGRPCLis); err != nil {
			loggerInst.Error().Err(err).Msg("Got GRPC server listener error")
			cancel()
		}
	}()

	loggerInst.Info().
		Str("name", cfg.ServiceName).
		Msg("Service started")

	ctx = newOSSignalContext(ctx)

	<-ctx.Done()

	loggerInst.Info().
		Str("name", cfg.ServiceName).
		Msg("Service shutdown")

	serverGRPC.GracefulStop()
	serverGRPC.Stop()
}

func newOSSignalContext(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		select {
		case <-osSignals:
			cancel()
		case <-ctx.Done():
			signal.Stop(osSignals)
		}
	}()

	return ctx
}
