// @Title socks5.go
// @Description Socks5 代理服务
// @Author Zero - 2024/8/8 23:45:20

package socks5

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	. "telepor/connection"
	"telepor/logger"
	"telepor/tool"
	"time"
)

// Server 仅支持 Socks5 的服务端
type Server struct {
	AuthMode     Method // 服务支持的认证模式

	Next         bool   // 是否转发到下级代理
	NextProtocol string // 下级代理服务类型 (支持 Socks5 -> HTTP)
	NextHost     string // 下级代理服务地址
	NextPort     string // 下级代理服务端口
}

func NewSocks5Server() *Server {
	return &Server{AuthMode: UserPassword}
}

// Startup 运行 Socks5 代理服务器
func (s *Server) Startup() error {
	return nil
}

// ServeHandle 处理 Socks5 请求
func (s *Server) ServeHandle(c *Connection) {
	if c, ok := c.Conn.(*net.TCPConn); ok {
		_ = c.SetKeepAlive(true)
	}
	// Step1 协商
	if err := s.shakeHands(c); err != nil {
		logger.Logger.ErrorSf("[Socks5] shakeHands fail: %v", err)
		return
	}
	// Step2 与请求目标服务建立连接
	tc, err := s.tunnel(c)
	if err != nil {
		logger.Logger.ErrorSf("[Socks5] connect target server fail: %v", err)
		return
	}
	defer tc.Close()
	// Step3 客户端连接 & 目标服务连接 进行数据交换
	//_ = tc.SetReadDeadline(time.Now().Add(time.Second * 3))
	//_ = c.SetReadDeadline(time.Now().Add(time.Second * 3))
	flow, err := tool.Swap(c, tc)
	if err != nil {
		logger.Logger.ErrorSf("[Socks5] request forward error: %v", err)
		return
	}
	logger.Logger.InfoSf("[Socks5] Client(%s) -> Server(%s) Succeeded, flow: %d", c.RemoteAddr().String(), tc.RemoteAddr().String(), flow)
}

// 与目标服务器建立连接，并且返回连接
func (s *Server) tunnel(conn *Connection) (tc *Connection, e error) {
	// Step1 解析请求报文
	req, e := ProxyRequestUnpack(conn)
	if e != nil {
		e = fmt.Errorf("error on parsing Socks5 request: %w", e)
		return
	}
	// Step2 请求报文校验
	if e = req.Checker(conn); e != nil {
		e = fmt.Errorf("error on invalid request message: %w", e)
		return
	}
	// Step3 与目标服务建立连接
	c, e := net.DialTimeout("tcp", req.Addr(), time.Second*10)
	if e != nil {
		// todo fix 根据不同的错误，响应不同的 REP
		_ = RequestFailureReply(conn, HostUnreachable)
		return nil, fmt.Errorf("failed to establish connection to target server: %w", e)
	}
	tc = WrapConn(c)
	// Step4 连接成功, Send Success Reply
	addrStr := c.RemoteAddr()
	addr := addrStr.(*net.TCPAddr)
	//ip := net.ParseIP(c.RemoteAddr().Network())
	e = RequestSuccessReply(conn, addr.IP, req.DstPort)
	if e != nil {
		_ = tc.Close()
		return
	}
	return
}

// 与 Socks5 客户端握手，协商双方认证方式
func (s *Server) shakeHands(conn *Connection) error {
	// 解析客户端协商报文
	req, err := shakeHandsMessageUnpack(conn)
	if err != nil {
		return err
	}
	if s.AuthMode == NoAuthentication {
		// 无需认证
		return s.shakeHandsReply(conn, NoAuthentication)
	}
	if s.AuthMode == UserPassword && tool.Contains(req.Methods, UserPassword) {
		// 用户名密码认证
		if _, err = s.AuthByUsernamePassword(conn); err != nil {
			// 认证失败
			_ = s.AuthReply(conn, AuthFailure)
			return fmt.Errorf("auth failed: %s", err.Error())
		}
		// 认证成功
		return s.AuthReply(conn, AuthSuccess)
	}
	// 没有可接受的认证方式
	_ = s.shakeHandsReply(conn, NoAcceptableMethods)
	return errors.New("no acceptable authentication method")
}

// RequestSuccessReply 代理请求 建立连接成功报文
// +----+-----+-------+------+----------+----------+
// |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
// +----+-----+-------+------+----------+----------+
// | 1  |  1  | X'00' |  1   | Variable |    2     |
// +----+-----+-------+------+----------+----------+
func RequestSuccessReply(c io.Writer, bindIP net.IP, bindPort uint16) (err error) {
	var addrType AddressType = IPv4
	if len(bindIP) == net.IPv6len {
		addrType = IPv6
	}
	// 缓冲区
	buffer := bytes.Buffer{}
	// Write VER\REP\RSV\ATYPE
	buffer.Write([]byte{Version, Succeeded, Reserved, addrType})
	// Write BND.ADDR
	buffer.Write(bindIP)
	// Write BND.PORT (以大端模式)
	portBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(portBytes, bindPort)
	buffer.Write(portBytes)
	_, err = c.Write(buffer.Bytes())
	return
}

// RequestFailureReply 代理请求处理失败，响应报文
func RequestFailureReply(c io.Writer, reply RequestReply) error {
	_, err := c.Write([]byte{Version, reply, Reserved, IPv4, 0, 0, 0, 0, 0, 0})
	return err
}

// 向客户端响应协商结果报文
// +----+--------+
// |VER | METHOD |
// +----+--------+
// | 1  |   1    |
// +----+--------+
func (s *Server) shakeHandsReply(conn *Connection, method Method) error {
	_, err := conn.Write([]byte{Version, method})
	return err
}
