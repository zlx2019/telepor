package main

import (
	"telepor/logger"
	server2 "telepor/server"
)

// @Title       main.go
// @Author      Zero.
// @Create      2024-08-09 15:20

func main() {
	server := server2.MixedServer{
		Addr: "127.0.0.1:9946",
	}
	if err := server.Startup(); err != nil {
		logger.Logger.PanicSf("Startup server failed: %s", err)
	}
}
