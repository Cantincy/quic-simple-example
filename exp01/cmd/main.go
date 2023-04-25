package main

import (
	"context"
	"crypto/tls"
	"github.com/quic-go/quic-go"
	"log"
	"pro01/quic/demo03/exp01/util"
	"time"
)

func main() {
	go func() {
		tlsConf := &tls.Config{
			InsecureSkipVerify: true,
			NextProtos:         []string{"quic-echo-example"},
		}
		connTo, err := quic.DialAddr("localhost:9090", tlsConf, nil)
		if err != nil {
			log.Fatal(err)
		}

		listener, err := quic.ListenAddr(":9091", util.GenerateTLSConfig(), nil)
		if err != nil {
			log.Fatal(err)
		}
		connFrom, err := listener.Accept(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		for {
			stream, err := connTo.OpenStreamSync(context.Background())
			if err != nil {
				log.Fatal(err)
			}
			buf := []byte("hello from client1...")
			stream.Write(buf)
			stream.Close()

			stream0, err := connFrom.AcceptStream(context.Background())
			if err != nil {
				log.Fatal(err)
			}
			buf0 := make([]byte, 100)
			stream0.Read(buf0)
			stream0.Close()
			log.Println("client1: ", string(buf0))
		}
	}()

	listener, err := quic.ListenAddr(":9090", util.GenerateTLSConfig(), nil)
	if err != nil {
		log.Fatal(err)
	}
	connFrom, err := listener.Accept(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
	}
	connTo, err := quic.DialAddr("localhost:9091", tlsConf, nil)
	if err != nil {
		log.Fatal(err)
	}

	for {
		stream, err := connFrom.AcceptStream(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		buf := make([]byte, 100)
		stream.Read(buf)
		stream.Close()
		log.Println("server: ", string(buf))

		stream0, err := connTo.OpenStreamSync(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		buf0 := []byte("reply from server...")
		stream0.Write(buf0)
		stream0.Close()

		time.Sleep(time.Second)
	}
}
