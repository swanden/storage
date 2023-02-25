package storage

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	grpcStorage "github.com/swanden/storage/api/go"
	storageUseCase "github.com/swanden/storage/internal/usecase/storage"
	"github.com/swanden/storage/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
	"time"
)

const (
	method         = "method"
	requestUUIDKey = "requestUUID"
)

type Logger interface {
	Debug() logger.Event
	Info() logger.Event
	Error() logger.Event
}

type StorageUseCase interface {
	Set(ctx context.Context, key, value string, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}

type Controller struct {
	grpcStorage.UnimplementedStorageServer
	log            Logger
	storageUseCase StorageUseCase
}

func New(options ...Options) (*Controller, error) {
	opts := getDefaultOptions()

	for _, opt := range options {
		if opt != nil {
			if err := opt(&opts); err != nil {
				return nil, err
			}
		}
	}

	if err := validate(opts); err != nil {
		return nil, err
	}

	return &Controller{
		log:            opts.log,
		storageUseCase: opts.storageUseCase,
	}, nil
}

func (c Controller) Set(ctx context.Context, request *grpcStorage.SetRequest) (*empty.Empty, error) {
	requestUUID := fmt.Sprintf("%v", ctx.Value(requestUUIDKey))

	c.log.Debug().
		Str(requestUUIDKey, requestUUID).
		Str(method, "[StorageController] [Set]").
		Msg("start set")
	defer c.log.Debug().
		Str(requestUUIDKey, requestUUID).
		Str(method, "[StorageController] [Set]").
		Msg("stop set")

	key := strings.TrimSpace(request.GetKey())
	if key == "" {
		c.log.Error().
			Str(requestUUIDKey, requestUUID).
			Str(method, "[StorageController] [Set]").
			Str("key", key).
			Msg(ErrBadKey.Error())

		return nil, status.Errorf(codes.InvalidArgument, ErrBadKey.Error())
	}

	value := strings.TrimSpace(request.GetValue())
	if value == "" {
		c.log.Error().
			Str(requestUUIDKey, requestUUID).
			Str(method, "[StorageController] [Set]").
			Str("key", key).
			Str("value", value).
			Msg(ErrBadValue.Error())

		return nil, status.Errorf(codes.InvalidArgument, ErrBadValue.Error())
	}

	ttl := time.Duration(request.GetTtl()) * time.Second
	if ttl < 0 {
		c.log.Error().
			Str(requestUUIDKey, requestUUID).
			Str(method, "[StorageController] [Set]").
			Str("key", key).
			Str("value", value).
			Str("ttl", ttl.String()).
			Msg(ErrBadTTL.Error())

		return nil, status.Errorf(codes.InvalidArgument, ErrBadTTL.Error())
	}

	err := c.storageUseCase.Set(ctx, key, value, ttl)
	if err != nil {
		c.log.Error().
			Str(requestUUIDKey, requestUUID).
			Str(method, "[StorageController] [Set]").
			Str("key", key).
			Str("value", value).
			Str("ttl", ttl.String()).
			Err(err).
			Msg(ErrSet.Error())

		return nil, status.Errorf(codes.Internal, ErrSet.Error())
	}

	return &empty.Empty{}, nil
}

func (c Controller) Get(ctx context.Context, request *grpcStorage.GetRequest) (*grpcStorage.GetResponse, error) {
	requestUUID := fmt.Sprintf("%v", ctx.Value(requestUUIDKey))

	c.log.Debug().
		Str(requestUUIDKey, requestUUID).
		Str(method, "[StorageController] [Get]").
		Msg("start get")
	defer c.log.Debug().
		Str(requestUUIDKey, requestUUID).
		Str(method, "[StorageController] [Get]").
		Msg("stop get")

	key := strings.TrimSpace(request.GetKey())
	if key == "" {
		c.log.Error().
			Str(requestUUIDKey, requestUUID).
			Str(method, "[StorageController] [Get]").
			Str("key", key).
			Msg(ErrBadKey.Error())

		return nil, status.Errorf(codes.InvalidArgument, ErrBadKey.Error())
	}

	value, err := c.storageUseCase.Get(ctx, key)
	if errors.Is(err, storageUseCase.ErrNotFound) {
		c.log.Info().
			Str(requestUUIDKey, requestUUID).
			Str(method, "[StorageController] [Get]").
			Str("key", key).
			Msg(ErrNotFound.Error())

		return nil, status.Errorf(codes.NotFound, ErrNotFound.Error())
	}
	if err != nil {
		c.log.Error().
			Str(requestUUIDKey, requestUUID).
			Str(method, "[StorageController] [Get]").
			Str("key", key).
			Str("value", value).
			Err(err).
			Msg(ErrGet.Error())

		return nil, status.Errorf(codes.Internal, ErrGet.Error())
	}

	return &grpcStorage.GetResponse{Value: value}, nil
}

func (c Controller) Delete(ctx context.Context, request *grpcStorage.DeleteRequest) (*empty.Empty, error) {
	requestUUID := fmt.Sprintf("%v", ctx.Value(requestUUIDKey))

	c.log.Debug().
		Str(requestUUIDKey, requestUUID).
		Str(method, "[StorageController] [Delete]").
		Msg("start delete")
	defer c.log.Debug().
		Str(requestUUIDKey, requestUUID).
		Str(method, "[StorageController] [Delete]").
		Msg("stop delete")

	key := strings.TrimSpace(request.GetKey())
	if key == "" {
		c.log.Error().
			Str(requestUUIDKey, requestUUID).
			Str(method, "[StorageController] [Delete]").
			Str("key", key).
			Msg(ErrBadKey.Error())

		return nil, status.Errorf(codes.InvalidArgument, ErrBadKey.Error())
	}

	err := c.storageUseCase.Delete(ctx, key)
	if err != nil {
		c.log.Error().
			Str(requestUUIDKey, requestUUID).
			Str(method, "[StorageController] [Delete]").
			Str("key", key).
			Err(err).
			Msg(ErrDelete.Error())

		return nil, status.Errorf(codes.Internal, ErrDelete.Error())
	}

	return &empty.Empty{}, nil
}
