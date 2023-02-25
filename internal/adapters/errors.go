package adapters

import "github.com/pkg/errors"

var (
	ErrNotFound = errors.New("memcached adapter: key not found")
)
