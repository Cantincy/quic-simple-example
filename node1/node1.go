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
	listener, err := quic.ListenAddr(util.Node1Addr, util.GenerateTLSConfig(), nil) // 开启监听器
	if err != nil {
		log.Fatal(err)
	}
	connFrom, err := listener.Accept(context.Background()) // 接收连接
	if err != nil {
		log.Fatal(err)
	}

	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
	}
	connTo, err := quic.DialAddr(util.Node2Addr, tlsConf, nil) // 发送连接
	if err != nil {
		log.Fatal(err)
	}

	for {
		// 接收
		stream, err := connFrom.AcceptStream(context.Background()) // 获取stream
		if err != nil {
			log.Fatal(err)
		}
		buf := make([]byte, 100)
		stream.Read(buf)
		stream.Close()

		log.Println(string(buf))

		// 发送
		stream0, err := connTo.OpenStreamSync(context.Background()) // 创建stream
		if err != nil {
			log.Fatal(err)
		}
		buf0 := []byte("hello from node1 " + time.Now().Format("2006-01-02 15:04:05"))
		stream0.Write(buf0)
		stream0.Close()
		time.Sleep(time.Second)
	}
}
