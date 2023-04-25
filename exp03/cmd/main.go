package main

import (
	"context"
	"fmt"
	"github.com/quic-go/quic-go"
	"log"
	"pro01/quic/demo03/exp03/util"
	"time"
)

func main() {
	// client协程
	for i := 0; i < util.ClientNum; i++ {
		go func(i int) {

			cliAddr := fmt.Sprintf("%s:%d", util.IP, util.Port+1+i)
			serverAddr := fmt.Sprintf("%s:%d", util.IP, util.Port)

			// 建立client的-->server连接
			connTo, err := quic.DialAddr(serverAddr, util.NewTlsConf(), nil)
			if err != nil {
				log.Fatal(err)
			}

			listener, err := quic.ListenAddr(cliAddr, util.GenerateTLSConfig(), nil)
			if err != nil {
				log.Fatal(err)
			}
			defer listener.Close()

			// 建立client的<--server连接
			connFrom, err := listener.Accept(context.Background())
			if err != nil {
				log.Fatal(err)
			}

			for {
				buf := []byte(fmt.Sprintf("hello from client%d", i))
				streamTo, err := connTo.OpenStreamSync(context.Background())
				if err != nil {
					log.Fatal(err)
				}
				streamTo.Write(buf)
				streamTo.Close()

				buf0 := make([]byte, 1024)
				streamFrom, err := connFrom.AcceptStream(context.Background())
				if err != nil {
					log.Fatal(err)
				}
				streamFrom.Read(buf0)
				streamFrom.Close()
				log.Printf("[Client%d]:%s", i, string(buf0))
			}

		}(i)
	}

	/*
		Client
			=============================================================================
		Server
	*/

	addr := fmt.Sprintf("%s:%d", util.IP, util.Port)
	listener, err := quic.ListenAddr(addr, util.GenerateTLSConfig(), nil)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	connFrom := make([]quic.Connection, util.ClientNum)
	connTo := make([]quic.Connection, util.ClientNum)
	for i := 0; i < util.ClientNum; i++ {
		connFrom[i], err = listener.Accept(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		cliAddr := fmt.Sprintf("%s:%d", util.IP, util.Port+1+i)
		connTo[i], err = quic.DialAddr(cliAddr, util.NewTlsConf(), nil)
		if err != nil {
			log.Fatal(err)
		}
	}

	for {
		for i := 0; i < util.ClientNum; i++ {
			streamFrom, err := connFrom[i].AcceptStream(context.Background())
			if err != nil {
				log.Fatal(err)
			}
			buf := make([]byte, 1024)
			streamFrom.Read(buf)
			streamFrom.Close()
			log.Printf("[Server]:%s", string(buf))

			buf0 := []byte("hello from server")
			streamTo, err := connTo[i].OpenStreamSync(context.Background())
			if err != nil {
				log.Fatal(err)
			}
			streamTo.Write(buf0)
			streamTo.Close()
		}
		time.Sleep(time.Second)
	}

}
