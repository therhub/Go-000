package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"therhub.com/wk04/configs"
	"therhub.com/wk04/pkg/login"
)

const split = ":"

// 1. 基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够一个退出，全部注销退出。
func main() {

	// 加载配置文件
	loadResult := configs.LoadYamlConfig()
	configs.SetConfiguration(configs.SetServerConfigs(loadResult))

	// 还可以设置是否为debug
	// configs.SetConfiguration(configs.SetIsDebug(true))

	// 产生一个新的group
	g, ctx := errgroup.WithContext(context.Background())

	// 异步监控 linux signal
	g.Go(func() error {
		return monitor(ctx)
	})

	// 启动http1
	g.Go(func() error {
		return listen(ctx, fmt.Sprintf("%s%v", split, configs.GetHttpConfig().Port), &selfHandler{})
	})

	// 启动grpc1
	g.Go(func() error {
		return grpcListen(ctx, fmt.Sprintf("%s%v", split, configs.GetGrpcConfig().Port))
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

	// 获取goroutine管理
	g, _ := errgroup.WithContext(input)

	g.Go(func() error {
		<-input.Done()
		s.Shutdown(input)
		return nil
	})

	return s.ListenAndServe()
}

// 服务器监听
func grpcListen(input context.Context, addr string) error {

	grp := grpc.NewServer()

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		panic(fmt.Errorf("webServer.grpc.net.Listen失败，错误信息为：%s", err))
	}

	// 获取goroutine管理
	g, _ := errgroup.WithContext(input)

	g.Go(func() error {
		<-input.Done()
		grp.Stop()
		return nil
	})

	// 注册
	login.RegisterFunc(grp)

	// 注册反射服务 这个服务是CLI使用的 跟服务本身没有关系
	reflection.Register(grp)

	// 返回
	return grp.Serve(lis)
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
