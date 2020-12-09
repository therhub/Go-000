package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

// 1. 基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够一个退出，全部注销退出。
func main() {

	stop := make(chan struct{})

	eg := &errgroup.Group{}

	// 起一个协程
	eg.Go(func() error {
		return monitor(stop)
	})

	// 另起一个协程
	eg.Go(func() error {
		return listen(":8080", &selfHandler{}, stop)
	})

	// 等待如果发生错误关闭stop
	if err := eg.Wait(); err != nil {
		close(stop)
		os.Exit(1)
	}
}

// 监听
func listen(addr string, handler http.Handler, stop <-chan struct{}) error {

	s := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go func() {
		<-stop
		s.Shutdown(context.Background())
	}()

	return s.ListenAndServe()
}

func monitor(stop chan struct{}) error {

	// chan
	ch := make(chan os.Signal)

	// 通过ch监听syscall
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)

	// 一般情况下堵塞
	select {
	case <-ch:
		stop <- struct{}{}
	case <-stop:
		close(ch)
	}

	return nil
}

// 自定义handler
type selfHandler struct {
}

// 实现ServeHTTP
func (s *selfHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello,world"))
}
