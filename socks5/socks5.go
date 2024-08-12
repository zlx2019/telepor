// @Title socks5.go
// @Description Socks5 代理服务
// @Author Zero - 2024/8/8 23:45:20

package socks5

import (
	"errors"
	"net"
	"telepor/connection"
	"telepor/logger"
	"telepor/tool"
)

// Server 仅支持 Socks5 的服务端
type Server struct {
	AuthMethod        int8   // 认证方式
	ForwardProxy      bool   // 是否转发到上级代理
	SuperiorProxyHost string // 上级代理服务地址
	SuperiorProxyPort int    // 上级代理服务端口
}

// Startup 运行 Socks5 代理服务器
func (s *Server) Startup() error {
	return nil
}

// Handle 处理 Socks5 请求
func (s *Server) Handle(conn *connection.Connection) error {
	if c, ok := conn.Conn.(*net.TCPConn); ok {
		_ = c.SetKeepAlive(true)
	}
	// TODO 协商
	if err := s.negotiate(conn); err != nil {
		return err
	}
	// TODO 尝试与目标服务器建立连接
	_, _ = s.bridge(conn)

	// TODO 双方数据转发
	return nil
}

// 与目标服务器建立连接，并且返回连接
func (s *Server) bridge(conn *connection.Connection) (*connection.Connection, error) {
	msg, err := ProxyRequestUnpack(conn)
	if err != nil {
		logger.Logger.ErrorSf("unpack socks5 proxy request failed: %c", err)
		return nil, err
	}
	logger.Logger.InfoSf("[Socks5 Request] VER: %d, CMD: %d, ATYP: %d, ADDR: %s",
		msg.Version, msg.Cmd, msg.Type, msg.Addr())
	return nil, nil
}

// 与 Socks5 客户端进行'协商'处理
func (s *Server) negotiate(conn *connection.Connection) error {
	// 解析协商报文
	pack, err := s.negotiationUnpack(conn)
	if err != nil {
		return err
	}
	// 回应协商结果
	if pack.NMethod <= 0 || len(pack.Methods) == 0 {
		_ = s.negotiationReply(conn, NoAcceptableMethods)
		return errors.New("the client does not have any authentication method")
	}
	switch {
	case tool.Contains(pack.Methods, NoAuthentication):
		// 无需认证
		err = s.negotiationReply(conn, NoAuthentication)
	case tool.Contains(pack.Methods, UserPassword):
		// 用户名密码认证
		err = s.negotiationReply(conn, UserPassword)
	}
	return err
}

// Socks5 回应协商结果
func (s *Server) negotiationReply(conn net.Conn, method Method) error {
	_, err := conn.Write([]byte{Version, method})
	return err
}
