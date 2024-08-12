// @Title mixed.go
// @Description 支持Socks5 && http 的混合服务端
// @Author Zero - 2024/8/11 10:12:46

package server

import (
	"net"
	"telepor/connection"
	"telepor/define"
	"telepor/http"
	"telepor/logger"
	"telepor/socks5"
	"telepor/tool"
)

// MixedServer 混合协议服务器 [HTTP | HTTPS | Socks5]
type MixedServer struct {
	Addr         string
	Socks5Server *socks5.Server
	HTTPServer   *http.Server
}

// Startup 启动 Mixed 服务
func (ms *MixedServer) Startup() error {
	listen, err := net.Listen("tcp", ms.Addr)
	if err != nil {
		return err
	}
	logger.Logger.InfoSf("[Mixed] HTTP & Socks5 server listen on %s", ms.Addr)
	for {
		conn, err := listen.Accept()
		if err != nil {
			logger.Logger.ErrorSf("[Mixed] Accept error: %v", err)
			continue
		}
		go ms.dispatch(conn)
	}
}

// TCP 连接处理，根据协议派发至具体的服务实现.
func (ms *MixedServer) dispatch(c net.Conn) {
	conn := connection.WrapConn(c)
	defer conn.Close()
	pact, err := tool.IdentifyPact(conn)
	if err != nil {
		logger.Logger.ErrorSf("不支持的协议: %s", err.Error())
		return
	}
	switch pact {
	case define.Socks5:
		err = ms.Socks5Server.Handle(conn)
	case define.HTTP:
	}
	if err != nil {
		logger.Logger.ErrorSf("handler err : %v", err)
	}
}
