// @Title mixed.go
// @Description 支持Socks5 && http 的混合服务端
// @Author Zero - 2024/8/11 10:12:46

package server

import (
	"net"
	"strings"
	"telepor/connection"
	"telepor/define"
	"telepor/logger"
	"telepor/socks5"
)

// MixedServer 混合协议服务器 [HTTP | HTTPS | Socks5]
type MixedServer struct {
	Addr         string
	Socks5Server *socks5.Socks5Server
	HTTPServer   *HTTPServer
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
	pact, err := IdentifyPact(conn)
	if err != nil {
		logger.Logger.ErrorSf("不支持的协议: %s", err.Error())
		return
	}
	logger.Logger.InfoSf("conn is %s", conn.RemoteAddr().String(), pact)
	switch pact {
	case define.Socks5:
		err = ms.Socks5Server.Handle(conn)
	case define.HTTP:
	}
	if err != nil {
		logger.Logger.ErrorSf("handler err : %v", err)
	}
}


// IdentifyPact 判断并且获取连接使用的协议类型
func IdentifyPact(c *connection.Connection) (int8,error) {
	// is Socks5？
	ver, _ := c.Peek(1)
	if ver[0] == socks5.Version {
		return define.Socks5, nil
	}

	// is HTTP? 通过HTTP请求行 Method 判断
	line, err := c.Peek(7)
	if err != nil {
		return 0, err
	}
	method := strings.ToUpper(string(line))
	if strings.HasPrefix(method, "GET") ||
		strings.HasPrefix(method, "POST") ||
		strings.HasPrefix(method, "PUT") ||
		strings.HasPrefix(method, "DELETE") ||
		strings.HasPrefix(method, "CONNECT") ||
		strings.HasPrefix(method, "OPTIONS") ||
		strings.HasPrefix(method, "TRACE") {
		return define.HTTP, nil
	}
	return define.Unknown, nil
}
