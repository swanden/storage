package pool

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"net"
	"sync"
	"time"
)

const (
	protocol          = "tcp"
	maxRequestsLength = 10_000
)

type request struct {
	response chan response
	ctx      context.Context
}

type response struct {
	connection net.Conn
	err        error
}

type Pool struct {
	host string
	port int

	mu        sync.Mutex
	idleConns []net.Conn

	openConns    int
	maxOpenConns int
	maxIdleConns int

	newConnTimeout   time.Duration
	connRetryTimeout time.Duration

	requests chan *request
}

func NewPool(host string, opts ...Option) (*Pool, error) {
	pool := &Pool{
		host:     host,
		requests: make(chan *request, maxRequestsLength),
	}

	for _, opt := range opts {
		opt(pool)
	}

	pool.idleConns = make([]net.Conn, 0, pool.maxIdleConns)

	go pool.handleConnectionRequest()

	return pool, nil
}

func (p *Pool) Put(connection net.Conn) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.maxIdleConns > len(p.idleConns) {
		p.idleConns = append(p.idleConns, connection)
		return
	}

	connection.Close()
	p.openConns--
}

func (p *Pool) Get(ctx context.Context) (net.Conn, error) {
	p.mu.Lock()

	idleConnCount := len(p.idleConns)
	if idleConnCount > 0 {
		for i, c := range p.idleConns {
			p.removeIdleConn(i)
			p.mu.Unlock()
			return c, nil
		}
	}

	if p.maxOpenConns > 0 && p.openConns >= p.maxOpenConns {
		req := &request{
			response: make(chan response, 1),
			ctx:      ctx,
		}

		p.requests <- req

		p.mu.Unlock()

		resp := <-req.response

		return resp.connection, resp.err
	}

	p.openConns++
	p.mu.Unlock()

	newConn, err := p.openNewConnection()
	if err != nil {
		p.mu.Lock()
		p.openConns--
		p.mu.Unlock()

		return nil, err
	}

	return newConn, nil
}

func (p *Pool) removeIdleConn(index int) {
	copy(p.idleConns[index:], p.idleConns[index+1:])
	p.idleConns = p.idleConns[:len(p.idleConns)-1]
}

func (p *Pool) openNewConnection() (net.Conn, error) {
	addr := fmt.Sprintf("%s:%d", p.host, p.port)

	d := net.Dialer{Timeout: p.newConnTimeout}
	c, err := d.Dial(protocol, addr)
	if err != nil {
		return nil, errors.Wrap(ErrServerConnect, err.Error())
	}

	return c, nil

}

func (p *Pool) handleConnectionRequest() {
	for req := range p.requests {
		timeout := time.After(p.connRetryTimeout)

	loop:
		for {
			select {
			case <-req.ctx.Done():
				req.response <- response{
					connection: nil,
					err:        ErrConnCanceled,
				}
				break loop
			case <-timeout:
				req.response <- response{
					connection: nil,
					err:        ErrConnTimeout,
				}
				break loop
			default:
				p.mu.Lock()

				idleConnCount := len(p.idleConns)
				if idleConnCount > 0 {
					for i, c := range p.idleConns {
						p.removeIdleConn(i)
						p.mu.Unlock()
						req.response <- response{
							connection: c,
							err:        nil,
						}
						break loop
					}
				}

				if p.maxOpenConns > 0 && p.openConns < p.maxOpenConns {
					p.openConns++
					p.mu.Unlock()

					c, err := p.openNewConnection()
					if err != nil {
						p.mu.Lock()
						p.openConns--
						p.mu.Unlock()
						break loop
					}
					req.response <- response{
						connection: c,
						err:        nil,
					}
					break loop
				}

				p.mu.Unlock()
			}
		}
	}
}

func (p *Pool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	close(p.requests)
	for range p.requests {
	}

	p.maxIdleConns = 0

	for _, conn := range p.idleConns {
		conn.Close()
	}
	p.idleConns = []net.Conn{}
}
