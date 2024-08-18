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
		c := connection.WrapConn(conn)
		go ms.ServeHandle(c)
	}
}

// ServeHandle Mixed 服务请求处理，将请求(TCP)派发至 Socks5 或 HTTP 代理服务
func (ms *MixedServer) ServeHandle(c *connection.Connection) {
	// panic handler
	defer func() {
		if err := recover(); err != nil {
			logger.Logger.ErrorSf("conn [%s] panic: %v", c.RemoteAddr().String(), err)
		}
		_ = c.Close()
	}()
	pact, err := tool.IdentifyPact(c)
	if err != nil {
		logger.Logger.ErrorSf("不支持的协议: %s", err.Error())
		return
	}
	switch pact {
	case define.Socks5:
		ms.Socks5Server.ServeHandle(c)
	case define.HTTP:
	}
}
