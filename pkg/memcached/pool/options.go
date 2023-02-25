package pool

import "time"

type Option func(*Pool)

func WithPort(port int) Option {
	return func(p *Pool) {
		p.port = port
	}
}

func WithMaxIdleConns(maxIdleConns int) Option {
	return func(p *Pool) {
		p.maxIdleConns = maxIdleConns
	}
}

func WithMaxOpenConns(maxOpenConns int) Option {
	return func(p *Pool) {
		p.maxOpenConns = maxOpenConns
	}
}

func WithNewConnTimeout(newConnTimeout time.Duration) Option {
	return func(p *Pool) {
		p.newConnTimeout = newConnTimeout
	}
}

func WithConnRetryTimeout(connRetryTimeout time.Duration) Option {
	return func(p *Pool) {
		p.connRetryTimeout = connRetryTimeout
	}
}
