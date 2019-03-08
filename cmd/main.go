package main

import (
	"go.uber.org/zap"
	"net"
	"os"
	"socket/pkg/epoll"
	"socket/pkg/logger"
)

func main() {
	if _, err := os.Stat("/tmp/keyword_match.sock"); err ==nil {
		if err := os.Remove("/tmp/keyword_match.sock"); err != nil {
			logger.Handle.Panic("please remove keyword_match sock",
				zap.String("err_info", err.Error()))
		}
	}
	socket, err := net.Listen("unix", "/tmp/keyword_match.sock")
	if err != nil {
		panic("socket start error")
	}
	epoll.Works.Start()
	defer epoll.Works.Closer()
	go epoll.Start()

	for {
		client, err := socket.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				logger.Handle.Panic("accept temp err",
					zap.String("err_info", ne.Error()))
				continue
			}
			logger.Handle.Panic("socket accept err",
				zap.String("err_info", err.Error()))
			return
		}
		if err := epoll.Epolls.Add(client); err != nil {
			logger.Handle.Panic("failed to add epoll",
				zap.String("err_info", err.Error()))
			err := client.Close()
			logger.Handle.Panic("failed to add epoll and close error ",
				zap.String("err_info", err.Error()))
		}
	}

}

