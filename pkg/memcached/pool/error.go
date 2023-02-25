package pool

import "github.com/pkg/errors"

var (
	ErrServerConnect = errors.New("memcached pool: unable to connect to memcached server")
	ErrConnTimeout   = errors.New("memcached pool: connection request timeout")
	ErrConnCanceled  = errors.New("memcached pool: connection request canceled")
)
