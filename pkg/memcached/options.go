package memcached

import "time"

type Option func(*Client)

func WithPort(port int) Option {
	return func(c *Client) {
		c.port = port
	}
}

func WithMaxIdleConns(maxIdleConns int) Option {
	return func(c *Client) {
		c.maxIdleConns = maxIdleConns
	}
}

func WithMaxOpenConns(maxOpenConns int) Option {
	return func(c *Client) {
		c.maxOpenConns = maxOpenConns
	}
}

func WithNewConnTimeout(newConnTimeout time.Duration) Option {
	return func(c *Client) {
		c.newConnTimeout = newConnTimeout
	}
}

func WithConnRetryTimeout(connRetryTimeout time.Duration) Option {
	return func(c *Client) {
		c.connRetryTimeout = connRetryTimeout
	}
}
