package main

// @Title       main.go
// @Author      Zero.
// @Create      2024-08-09 15:20

import (
	"telepor/http"
	"telepor/logger"
	"telepor/server"
	"telepor/socks5"
)

func main() {
	serv := server.MixedServer{
		Addr:         "127.0.0.1:15001",
		Socks5Server: socks5.NewSocks5Server(),
		HTTPServer:   http.NewHttpServer(),
	}
	if err := serv.Startup(); err != nil {
		logger.Logger.PanicSf("Startup server failed: %s", err)
	}
}
