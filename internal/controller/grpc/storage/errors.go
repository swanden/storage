package storage

import "github.com/pkg/errors"

var (
	ErrBadLogger         = errors.New("grpc controller: bad logger implementation")
	ErrBadStorageUseCase = errors.New("grpc controller: bad storage use case implementation")
	ErrSet               = errors.New("grpc controller: unable to set key-value pair")
	ErrBadKey            = errors.New("grpc controller: bad key")
	ErrBadValue          = errors.New("grpc controller: bad value")
	ErrBadTTL            = errors.New("grpc controller: ttl must be grater or equals 0")
	ErrGet               = errors.New("grpc controller: unable to get value")
	ErrDelete            = errors.New("grpc controller: unable to delete value")
	ErrNotFound          = errors.New("grpc controller: value not found")
)
