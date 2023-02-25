package storage

import "github.com/pkg/errors"

var (
	ErrBadLogger = errors.New("usecase: bad logger implementation")
	ErrSet       = errors.New("usecase: unable to set key-value pair")
	ErrGet       = errors.New("usecase: unable to get value")
	ErrDelete    = errors.New("usecase: unable to delete value")
	ErrNotFound  = errors.New("usecase: value not found")
)
