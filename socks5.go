// @Title socks5.go
// @Description Socks5 服务端
// @Author Zero - 2024/8/8 23:45:20

package telepor

import (
	"fmt"
	"log/slog"
	"net"
)

// Socks5Server Socks5 代理服务端
type Socks5Server struct {
	// socks5 IP
	IP   string
	// socks5 Port
	Port int
}

func (s *Socks5Server) Addr() string {
	return fmt.Sprintf("%s:%d", s.IP, s.Port)
}

// Run 运行 Socks5 代理服务器
func (s *Socks5Server) Run() error {
	listen, err := net.Listen("tcp", s.Addr())
	if err != nil {
		return err
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			slog.Error("accept fail: " + err.Error())
			continue
		}
		// 处理每一个连接
		go func() {
			defer conn.Close()
			err := s.handleConn(conn)
			if err != nil {
				slog.Error("handle connection fail: " + err.Error())
			}
		}()
	}
}

// 处理连接(Socks5)
func (s *Socks5Server) handleConn(conn net.Conn) error {
	// 协商
	if err := s.auth(); err != nil {
		return err
	}

	// 请求转发
	return nil
}

// socks5 协议协商处理
func (s *Socks5Server) auth() error {
	return nil
}