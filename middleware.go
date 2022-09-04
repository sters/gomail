package gomail

import (
	"context"
	"io"
	"net"
	"sync"
)

type (
	DialFunction    func(ctx context.Context, network string, addr string) (net.Conn, error)
	DialMiddleware  func(ctx context.Context, network string, addr string, next DialFunction) (net.Conn, error)
	DialMiddlewares []DialMiddleware

	SendFunction    func(ctx context.Context, from string, to []string, msg io.WriterTo) error
	SendMiddleware  func(ctx context.Context, from string, to []string, msg io.WriterTo, next SendFunction) error
	SendMiddlewares []SendMiddleware
)

func invokeDial(ctx context.Context, middlewares DialMiddlewares, root DialFunction, network string, addr string) (net.Conn, error) {
	wrapper := func(m DialMiddleware, current DialFunction) DialFunction {
		return func(ctx context.Context, network string, addr string) (net.Conn, error) {
			return m(ctx, network, addr, current)
		}
	}

	current := root
	for i := len(middlewares) - 1; i >= 0; i-- {
		current = wrapper(middlewares[i], current)
	}

	return current(ctx, network, addr)
}

func invokeSend(ctx context.Context, middlewares SendMiddlewares, root SendFunction, from string, to []string, msg io.WriterTo) error {
	wrapper := func(m SendMiddleware, current SendFunction) SendFunction {
		return func(ctx context.Context, from string, to []string, msg io.WriterTo) error {
			return m(ctx, from, to, msg, current)
		}
	}

	current := root
	for i := len(middlewares) - 1; i >= 0; i-- {
		current = wrapper(middlewares[i], current)
	}

	return current(ctx, from, to, msg)
}

type OnMemoryDialStatus struct {
	Dialed int
	Errs   int
	m      sync.RWMutex
}

func WithOnMemoryDialStats(s *OnMemoryDialStatus) DialMiddleware {
	return func(ctx context.Context, network string, addr string, next DialFunction) (net.Conn, error) {
		s.m.Lock()
		s.Dialed++
		s.m.Unlock()

		conn, err := next(ctx, network, addr)

		if err != nil {
			s.m.Lock()
			s.Errs++
			s.m.Unlock()
		}

		return conn, err
	}
}

type OnMemorySendStats struct {
	Sent int
	Errs int
	m    sync.RWMutex
}

func WithOnMemorySendStats(s *OnMemorySendStats) SendMiddleware {
	return func(ctx context.Context, from string, to []string, msg io.WriterTo, next SendFunction) error {
		s.m.Lock()
		s.Sent++
		s.m.Unlock()

		err := next(ctx, from, to, msg)

		if err != nil {
			s.m.Lock()
			s.Errs++
			s.m.Unlock()
		}

		return err
	}
}
