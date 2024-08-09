// @Title 		server.go
// @Description 服务端
// @Author Zero - 2024/8/8 23:44:14

package main

import "net"

type Server interface {
	Startup() error
}

// MixedServer 混合协议服务器 [HTTP | HTTPS | Socks5]
type MixedServer struct {
	Addr         string
	Socks5Server *Socks5Server
}

// Startup 启动服务
func (ms *MixedServer) Startup() error {
	listen, err := net.Listen("tcp", ms.Addr)
	if err != nil {
		return err
	}
	Logger.InfoSf("[Mixed] HTTP & Socks5 server listen on %s", ms.Addr)
	for {
		conn, err := listen.Accept()
		if err != nil {
			Logger.ErrorSf("[Mixed] Accept error: %v", err)
			continue
		}
		go ms.serve(conn)
	}
}

// 处理不同的协议
func (ms *MixedServer) serve(conn net.Conn) {

}
