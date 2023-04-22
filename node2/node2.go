package main

import (
	"context"
	"crypto/tls"
	"github.com/quic-go/quic-go"
	"log"
	"pro01/quic/demo03/util"
	"time"
)

func main() {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
	}
	connTo, err := quic.DialAddr(util.Node1Addr, tlsConf, nil) // 发送连接
	if err != nil {
		log.Fatal(err)
	}

	listener, err := quic.ListenAddr(util.Node2Addr, util.GenerateTLSConfig(), nil) // 开启监听器
	if err != nil {
		log.Fatal(err)
	}
	connFrom, err := listener.Accept(context.Background()) // 接收连接
	if err != nil {
		log.Fatal(err)
	}

	for {
		// 发送
		stream, err := connTo.OpenStreamSync(context.Background()) // 创建stream
		if err != nil {
			log.Fatal(err)
		}
		buf := []byte("hello from node2 " + time.Now().Format("2006-01-02 15:04:05"))
		stream.Write(buf)
		stream.Close()

		// 接收
		stream0, err := connFrom.AcceptStream(context.Background()) // 获取stream
		if err != nil {
			log.Fatal(err)
		}
		buf0 := make([]byte, 100)
		stream0.Read(buf0)
		stream0.Close()
		log.Println(string(buf0))
		time.Sleep(time.Second)
	}
}
