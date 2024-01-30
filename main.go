package main

import (
	"context"
	"io"
	"log"
	"net"

	"github.com/luiscovelo/icecast-server-to-hls/ffmpeg"
	"github.com/luiscovelo/icecast-server-to-hls/parser"
)

var (
	reader *io.PipeReader
	writer *io.PipeWriter
)

func init() {
	reader, writer = io.Pipe()
}

func main() {
	server1()
}

func server1() {
	ln, err := net.Listen("tcp", ":8090")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("listenning on :8090 port...")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handler(conn)
	}
}

func handler(conn net.Conn) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tmp := make([]byte, 4096)
	for {
		select {
		case <-ctx.Done():
			if err := reader.Close(); err != nil {
				log.Println("failed to close reader", err)
			}
			if err := writer.Close(); err != nil {
				log.Println("failed to close writer", err)
			}
			return
		default:
			n, err := conn.Read(tmp)
			if err != nil {
				if err == io.EOF {
					log.Println("connection closed")
					return
				}

				log.Println("error on read", err)
				continue
			}

			//log.Println("rcv raw", string(tmp[:n]))

			req := parser.Parse(tmp[:n])

			if req.Method == "PUT" {
				go ffmpeg.FFmpeg(ctx, reader)
				conn.Write([]byte("HTTP/1.1 100 Continue\r\n"))
				continue
			}

			if req.Method == "GET" {
				conn.Write([]byte("HTTP/1.1 200 OK\r\n"))
				conn.Close()
				return
			}

			if len(req.Body) > 0 {
				if _, err := writer.Write(req.Body); err != nil {
					log.Println("error to write", err)
					continue
				}

				conn.Write([]byte("HTTP/1.1 200 OK\r\n"))
			}
		}
	}
}
