// @Title socks5.go
// @Description Socks5 代理服务
// @Author Zero - 2024/8/8 23:45:20

package socks5

import (
	"errors"
	"io"
	"net"
	"telepor/connection"
	"telepor/tool"
)

// Socks5Server 仅支持 Socks5 的服务端
type Socks5Server struct {
	AuthMethod        int8   // 认证方式
	ForwardProxy      bool   // 是否转发到上级代理
	SuperiorProxyHost string // 上级代理服务地址
	SuperiorProxyPort int    // 上级代理服务端口
}

// Startup 运行 Socks5 代理服务器
func (s *Socks5Server) Startup() error {
	return nil
}

// Handle 处理 Socks5 请求
func (s *Socks5Server) Handle(conn *connection.Connection) error {
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
func (s *Socks5Server) bridge(conn *connection.Connection) (*connection.Connection, error) {
	return nil, nil
}




// 与 Socks5 客户端进行'协商'处理
func (s *Socks5Server) negotiate(conn *connection.Connection) error {
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
func (s *Socks5Server) negotiationReply(conn net.Conn, method Method) error {
	_, err := conn.Write([]byte{Version, method})
	return err
}


// RequestPack 代理请求报文(1 + 1 + 1 + 1 + 2 + variable)
// +----+-----+-------+------+----------+----------+
// |VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
// +----+-----+-------+------+----------+----------+
// | 1  |  1  | X'00' |  1   | Variable |    2     |
// +----+-----+-------+------+----------+----------+
// Cmd: 代理请求指令
// Type: 目标主机地址类型
// DstHost: 目标主机
// DstPort: 目标端口
type RequestPack struct {
	Version byte
	Cmd     Command
	Rsv     byte
	Type    AddressType
	DstHost string
	DstPort uint16
}
// ProxyRequestUnpack 解析 Socks5 代理请求报文
func ProxyRequestUnpack(c io.Reader) (*RequestPack, error) {
	// Read VER、CMD、RSV、A_TYPE
	buf := make([]byte, 4)
	if _, err := io.ReadFull(c, buf); err != nil {
		return nil, err
	}
	// 报文参数校验
	version, cmd, atyp := buf[0], buf[1], buf[4]
	if Version != version {
		return nil, VersionNotSupportedErr
	}
	return nil, nil
}


// AuthMethodPack Socks5 协商报文
type AuthMethodPack struct {
	Version byte
	NMethod byte
	Methods []byte
}
// Socks5 解析 Socks5 协商报文 (1 + 1 + (1 ~ 255))
// +----+-----------+----------+
// |VER | N_METHODS | METHODS  |
// +----+-----------+----------+
// | 1  |    1      | 1 to 255 |
// +----+-----------+----------+
// VER: Socks版本。
// METHODS: 一个动态的字节序列，表示客户端支持的认证方式。
// N_METHODS: `METHODS`序列的长度。
func (s *Socks5Server) negotiationUnpack(conn net.Conn) (pack *AuthMethodPack, err error) {
	buf := make([]byte, 2)
	// Read VER and N_METHODS
	if _, err = io.ReadFull(conn, buf); err != nil {
		return
	}
	version := buf[0]
	length := buf[1]
	// Read METHODS
	methodBuf := make([]byte, length)
	if _, err = io.ReadFull(conn, methodBuf); err != nil {
		return
	}
	pack = &AuthMethodPack{
		Version: version,
		NMethod: length,
		Methods: methodBuf,
	}
	return
}
