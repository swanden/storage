package storage

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/swanden/storage/internal/adapters"
	"github.com/swanden/storage/pkg/logger"
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

type Storage interface {
	Set(ctx context.Context, key, value string, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}

type UseCase struct {
	log     Logger
	storage Storage
}

func New(options ...Options) (*UseCase, error) {
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

	return &UseCase{
		log:     opts.log,
		storage: opts.storage,
	}, nil
}

func (uc UseCase) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	requestUUID := fmt.Sprintf("%v", ctx.Value(requestUUIDKey))

	uc.log.Debug().
		Str(requestUUIDKey, requestUUID).
		Str(method, "[StorageUseCase] [Set]").
		Str("key", key).
		Str("value", value).
		Str("ttl", ttl.String()).
		Msg("start set")
	defer uc.log.Debug().
		Str(requestUUIDKey, requestUUID).
		Str(method, "[StorageUseCase] [Set]").
		Str("key", key).
		Str("value", value).
		Str("ttl", ttl.String()).
		Msg("stop set")

	err := uc.storage.Set(ctx, key, value, ttl)
	if err != nil {
		uc.log.Error().
			Str(requestUUIDKey, requestUUID).
			Str(method, "[StorageUseCase] [Set]").
			Str("key", key).
			Str("value", value).
			Str("ttl", ttl.String()).
			Err(err).
			Msg(ErrSet.Error())

		return errors.Wrap(ErrSet, err.Error())

	}

	return nil
}

func (uc UseCase) Get(ctx context.Context, key string) (string, error) {
	requestUUID := fmt.Sprintf("%v", ctx.Value(requestUUIDKey))

	uc.log.Debug().
		Str(requestUUIDKey, requestUUID).
		Str(method, "[StorageUseCase] [Get]").
		Str("key", key).
		Msg("start get")
	defer uc.log.Debug().
		Str(requestUUIDKey, requestUUID).
		Str(method, "[StorageUseCase] [Get]").
		Str("key", key).
		Msg("stop get")

	value, err := uc.storage.Get(ctx, key)
	if errors.Is(err, adapters.ErrNotFound) {
		uc.log.Info().
			Str(requestUUIDKey, requestUUID).
			Str(method, "[StorageUseCase] [Get]").
			Str("key", key).
			Msg(ErrNotFound.Error())

		return "", errors.Wrap(ErrNotFound, err.Error())
	}
	if err != nil {
		uc.log.Error().
			Str(requestUUIDKey, requestUUID).
			Str(method, "[StorageUseCase] [Get]").
			Str("key", key).
			Str("value", value).
			Err(err).
			Msg(ErrGet.Error())

		return "", errors.Wrap(ErrGet, err.Error())
	}

	return value, nil
}

func (uc UseCase) Delete(ctx context.Context, key string) error {
	requestUUID := fmt.Sprintf("%v", ctx.Value(requestUUIDKey))

	uc.log.Debug().
		Str(requestUUIDKey, requestUUID).
		Str(method, "[StorageUseCase] [Delete]").
		Str("key", key).
		Msg("start delete")
	defer uc.log.Debug().
		Str(requestUUIDKey, requestUUID).
		Str(method, "[StorageUseCase] [Delete]").
		Str("key", key).
		Msg("stop delete")

	err := uc.storage.Delete(ctx, key)
	if err != nil {
		uc.log.Error().
			Str(requestUUIDKey, requestUUID).
			Str(method, "[StorageUseCase] [Delete]").
			Str("key", key).
			Err(err).
			Msg(ErrDelete.Error())

		return errors.Wrap(ErrDelete, err.Error())

	}

	return nil
}
