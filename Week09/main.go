package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	lis, err := net.Listen("tcp", port)

	if err != nil {
		panic(err)
	}

	for {
		conn, err := lis.Accept()

		if err != nil {
			panic(err)
		}

		go handler(&TcpConn{Conn: conn, Context: context.Background(), rc: make(chan Message), wr: make(chan Message), buf: make([]byte, 1024), active: time.NewTimer(time.Second * time.Duration(timeOut))})
	}
}

var (
	port    string = ":8888"
	timeOut int    = 5
)

type TcpConn struct {
	net.Conn
	context.Context              // 上下文
	rc              chan Message // 接受
	wr              chan Message // 发送
	active          *time.Timer  // 活跃时间
	buf             []byte
}

type Message struct {
	User  string
	Value string
}

// 处理连接
func handler(conn *TcpConn) {
	defer conn.Close()

	eg, _ := errgroup.WithContext(conn.Context)

	// 监控如果上下文关闭，关闭连接
	eg.Go(func() error {

		sig := make(chan os.Signal)
		signal.Notify(sig, os.Kill, os.Interrupt)

		select {
		case <-sig:
			conn.Close()
		case <-conn.active.C:
			conn.Close()
		}

		return nil
	})

	// 写入
	eg.Go(func() error {
		for {
			n, err := conn.Read(conn.buf)

			if err != nil {
				return err
			}

			refresh(conn)

			var msg = &Message{}
			err = json.Unmarshal(conn.buf[:n], &msg)

			if err != nil {
				return err
			}

			fmt.Printf("receive user:%v,message:%v", msg.User, msg.Value)
			conn.wr <- Message{User: "system", Value: "success"}
		}
	})

	// 发送
	eg.Go(func() error {
		for {
			select {
			case msg := <-conn.wr:
				refresh(conn)
				send(conn, &msg)
			case <-conn.Done():
				return nil
			}
		}
	})

	eg.Wait()
}

// 刷新连接存在时间
func refresh(conn *TcpConn) {
	conn.active.Reset(time.Duration(timeOut) * time.Second)
}

// 业务处理message
func handlerMessage(conn *TcpConn, msg *Message) {
	fmt.Printf("%s发来消息%s", msg.User, msg.Value)
	conn.wr <- Message{User: "System", Value: "GetValueSuccess"}
}

// 发送方法
func send(conn *TcpConn, msg *Message) {
	r, err := json.Marshal(msg)
	if err == nil {
		conn.Conn.Write(r)
	}
}
