package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

// 1. 基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够一个退出，全部注销退出。
func main() {

	// 产生一个新的group
	g, ctx := errgroup.WithContext(context.Background())

	// 异步监控 linux signal
	g.Go(func() error {
		return monitor(ctx)
	})

	// 启动http1
	g.Go(func() error {
		return listen(ctx, ":8080", &selfHandler{})
	})

	// 启动http2
	g.Go(func() error {
		return listen(ctx, ":8081", &selfHandler{})
	})

	// 等待如果发生错误关闭stop
	if err := g.Wait(); err != nil {
		print(err)
	}

	println(ctx.Err())
}

// 服务器监听
func listen(input context.Context, addr string, handler http.Handler) error {

	s := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go func() {
		<-input.Done()
		s.Shutdown(input)
	}()

	return s.ListenAndServe()
}

// linux signal监控
func monitor(ctx context.Context) error {

	// chan
	ch := make(chan os.Signal)

	// 通过ch监听syscall
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)

	// 新建一个context
	ctx, cancel := context.WithCancel(ctx)

	// 一般情况下堵塞
	select {
	case <-ch:
		cancel()
		return errors.New("systempause")
	case <-ctx.Done():
		return ctx.Err()
	}
}

// 自定义handler
type selfHandler struct {
}

// 实现ServeHTTP
func (s *selfHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello,world"))
}
