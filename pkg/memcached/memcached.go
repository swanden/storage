package memcached

import (
	"bytes"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/swanden/storage/pkg/memcached/pool"
	"io"
	"net"
	"strings"
	"time"
)

const (
	defaultPort             = 11211
	defaultMaxIdleConns     = 10
	defaultMaxOpenConns     = 10
	defaultNewConnTimeout   = 3000 * time.Millisecond
	defaultConnRetryTimeout = 3000 * time.Millisecond

	EOL                 = "\r\n"
	ResponseEnd         = "END" + EOL
	ResponseStored      = "STORED" + EOL
	ResponseNotFound    = "NOT_FOUND" + EOL
	ResponseDeleted     = "DELETED" + EOL
	ResponseError       = "ERROR"
	ResponseClientError = "CLIENT_ERROR"
	ResponseServerError = "SERVER_ERROR"
	Metadata            = 0
)

type Client struct {
	host             string
	port             int
	maxIdleConns     int
	maxOpenConns     int
	newConnTimeout   time.Duration
	connRetryTimeout time.Duration
	pool             *pool.Pool
}

func Connect(host string, opts ...Option) (*Client, error) {
	client := &Client{
		host:             host,
		port:             defaultPort,
		maxIdleConns:     defaultMaxIdleConns,
		maxOpenConns:     defaultMaxOpenConns,
		newConnTimeout:   defaultNewConnTimeout,
		connRetryTimeout: defaultConnRetryTimeout,
	}

	for _, opt := range opts {
		opt(client)
	}

	connPool, err := pool.NewPool(
		client.host,
		pool.WithPort(client.port),
		pool.WithMaxIdleConns(client.maxIdleConns),
		pool.WithMaxOpenConns(client.maxOpenConns),
		pool.WithNewConnTimeout(client.newConnTimeout),
		pool.WithConnRetryTimeout(client.connRetryTimeout),
	)
	if err != nil {
		return nil, errors.Wrap(ErNewPool, err.Error())
	}

	client.pool = connPool

	return client, nil
}

// Set sets key-value pair
// ttl - expiration time in seconds, if 0 - no expire time
func (c *Client) Set(ctx context.Context, key string, value string, ttl int) error {
	conn, err := c.pool.Get(ctx)
	if err != nil {
		return errors.Wrap(ErrGetConn, err.Error())
	}
	defer c.pool.Put(conn)

	_, err = fmt.Fprintf(conn, "set %s %d %d %d%s%s%s", key, Metadata, ttl, len(value), EOL, value, EOL)
	if err != nil {
		return errors.Wrap(ErrConnWrite, err.Error())
	}

	resp, err := c.getResponse(conn)
	if err != nil {
		return errors.Wrap(ErrSet, err.Error())
	}
	if bytes.Contains(resp, []byte(ResponseError)) {
		return ErrSet
	}

	return nil
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	conn, err := c.pool.Get(ctx)
	if err != nil {
		return "", errors.Wrap(ErrGetConn, err.Error())
	}
	defer c.pool.Put(conn)

	_, err = fmt.Fprintf(conn, "get %s%s", key, EOL)
	if err != nil {
		return "", errors.Wrap(ErrConnWrite, err.Error())
	}

	resp, err := c.getResponse(conn)
	if err != nil {
		return "", errors.Wrap(ErrGet, err.Error())
	}
	if bytes.Contains(resp, []byte(ResponseError)) {
		return "", ErrGet
	}

	data, ok := getData(string(resp))

	if !ok {
		return "", ErrNotFound
	}

	return data, nil
}

func (c *Client) Delete(ctx context.Context, key string) error {
	conn, err := c.pool.Get(ctx)
	if err != nil {
		return errors.Wrap(ErrGetConn, err.Error())
	}
	defer c.pool.Put(conn)

	_, err = fmt.Fprintf(conn, "delete %s%s", key, EOL)
	if err != nil {
		return errors.Wrap(ErrConnWrite, err.Error())
	}

	_, err = c.getResponse(conn)
	if err != nil {
		return errors.Wrap(ErrDelete, err.Error())
	}

	return nil
}

func getData(resp string) (string, bool) {
	resp = strings.ReplaceAll(resp, ResponseEnd, "")
	data := strings.Split(resp, "\r\n")
	if len(data) > 2 {
		return data[1], true
	}

	return "", false
}

func (c *Client) getResponse(conn net.Conn) ([]byte, error) {
	tmp := make([]byte, 1024)
	data := make([]byte, 0)

	length := 0

	for {
		n, err := conn.Read(tmp)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, errors.Wrap(ErrConnRead, err.Error())
		}

		data = append(data, tmp[:n]...)

		if bytes.Contains(data, []byte(ResponseError)) || bytes.Contains(data, []byte(ResponseClientError)) {
			return nil, ErrClient
		}
		if bytes.Contains(data, []byte(ResponseServerError)) {
			return nil, ErrServer
		}
		if bytes.Contains(data, []byte(ResponseStored)) || bytes.Contains(data, []byte(ResponseEnd)) {
			break
		}
		if bytes.Contains(data, []byte(ResponseNotFound)) || bytes.Contains(data, []byte(ResponseDeleted)) {
			break
		}

		length += n
	}

	return data, nil
}

func (c *Client) Close() {
	c.pool.Close()
}
