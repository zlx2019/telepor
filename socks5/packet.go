package socks5

// @Title       packet.go
// @Description Socks5 数据报文处理
// @Author      Zero.
// @Create      2024-08-12 15:11

import (
	"encoding/binary"
	"io"
	"net"
	"strconv"
	"telepor/pool"
)

// Socks5 协议报文解析，参考与 https://datatracker.ietf.org/doc/rfc1928/

// SocksRequestMsg 代理请求报文(1 + 1 + 1 + 1 + 2 + variable)
// +----+-----+-------+------+----------+----------+
// |VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
// +----+-----+-------+------+----------+----------+
// | 1  |  1  | X'00' |  1   | Variable |    2     |
// +----+-----+-------+------+----------+----------+
// Cmd: 代理请求指令
// AddrType: 目标主机地址类型
// DstHost: 目标主机
// DstPort: 目标端口
type SocksRequestMsg struct {
	Version  byte
	Cmd      Command
	Rsv      byte
	AddrType AddressType
	DstHost  string
	DstPort  uint16
}

// Addr 获取代理目标地址信息
func (sr *SocksRequestMsg) Addr() string {
	return net.JoinHostPort(sr.DstHost, strconv.Itoa(int(sr.DstPort)))
}

// ProxyRequestUnpack 解析 Socks5 代理请求报文
func ProxyRequestUnpack(c io.Reader) (*SocksRequestMsg, error) {
	// Read `VER`、`CMD`、`RSV`、`A_TYPE`
	buf := pool.Borrow(4)
	if _, err := io.ReadFull(c, buf); err != nil {
		return nil, err
	}
	msg := &SocksRequestMsg{}
	// 报文参数校验
	msg.Version, msg.Cmd, msg.AddrType = buf[0], buf[1], buf[3]
	if Version != msg.Version {
		return nil, VersionNotSupportedErr
	}
	if Connect != msg.Cmd {
		return nil, CommandNotSupportedErr
	}
	// Read `DST.ADDR` [IPv4 | IPv6 | 域名]
	switch msg.AddrType {
	case IPv4:
	case IPv6:
		pool.Revert(buf)
		buf = pool.Borrow(net.IPv6len)
	case Domain:
		if _, err := io.ReadFull(c, buf[:1]); err != nil {
			return nil, err
		}
		if int(buf[0]) > cap(buf) {
			pool.Revert(buf)
			buf = pool.Borrow(int(buf[0]))
		}
	}
	defer pool.Revert(buf)
	if _, err := io.ReadFull(c, buf); err != nil {
		return nil, err
	}
	if msg.AddrType == Domain {
		msg.DstHost = string(buf)
	} else {
		msg.DstHost = net.IP(buf).String()
	}
	// Read `DST.PORT` (大端处理)
	if _, err := io.ReadFull(c, buf[:2]); err != nil {
		return nil, err
	}
	msg.DstPort = binary.BigEndian.Uint16(buf[:2])
	return msg, nil
}

// AuthMethodMsg Socks5 协商报文
type AuthMethodMsg struct {
	Version byte
	NMethod byte
	Methods []byte
}

// NegotiationUnpack 解析 Socks5 协议协商报文
// +----+-----------+----------+
// |VER | N_METHODS | METHODS  |
// +----+-----------+----------+
// | 1  |    1      | 1 to 255 |
// +----+-----------+----------+
// VER: Socks版本
// N_METHODS: `METHODS`序列的长度
// METHODS: 一个动态的字节序列，表示客户端支持的认证方式
func NegotiationUnpack(conn io.Reader) (pack *AuthMethodMsg, err error) {
	buf := pool.Borrow(2)
	defer pool.Revert(buf)
	// Read `VER` and `N_METHODS`
	if _, err = io.ReadFull(conn, buf); err != nil {
		return
	}
	version := buf[0]
	length := buf[1]
	// Read `METHODS`
	methodBuf := pool.Borrow(int(length))
	defer pool.Revert(methodBuf)
	if _, err = io.ReadFull(conn, methodBuf); err != nil {
		return
	}
	pack = &AuthMethodMsg{
		Version: version,
		NMethod: length,
		Methods: methodBuf,
	}
	return
}
