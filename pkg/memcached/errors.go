package memcached

import "github.com/pkg/errors"

var (
	ErrClient    = errors.New("memcached: client error")
	ErrServer    = errors.New("memcached: server error")
	ErrSet       = errors.New("memcached: unable to set key-value pair")
	ErNewPool    = errors.New("memcached: unable to create connection pool")
	ErrGetConn   = errors.New("memcached: unable to get connection from pool")
	ErrConnWrite = errors.New("memcached: unable to write to connection")
	ErrConnRead  = errors.New("memcached: unable to read from connection")
	ErrGet       = errors.New("memcached: unable to get value from the store")
	ErrDelete    = errors.New("memcached: unable to delete value")
	ErrNotFound  = errors.New("memcached: value not found")
)
