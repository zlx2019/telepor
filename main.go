package main

// @Title       main.go
// @Author      Zero.
// @Create      2024-08-09 15:20

import (
	"telepor/config"
	_ "telepor/config"
	"telepor/http"
	"telepor/logger"
	"telepor/server"
	"telepor/socks5"
)

func main() {
	serv := server.MixedServer{
		Addr:         config.Conf.Bind,
		Socks5Server: socks5.NewSocks5Server(),
		HTTPServer:   http.NewHttpServer(),
	}
	if err := serv.Startup(); err != nil {
		logger.Logger.PanicSf("Startup server failed: %s", err)
	}
}
