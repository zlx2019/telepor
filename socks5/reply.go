// @Title reply.go
// @Description Socks5 响应报文体
// @Author Zero - 2024/8/17 15:37:37

package socks5

import "net"

// RequestReply 代理请求连接响应
type RequestReply = byte

const (
	Succeeded               RequestReply = 0x00 // 连接成功
	SocksServerFailure                   = 0x01 // 当前代理服务器错误
	NotAllowedByRuleset                  = 0x02 // 不支持的规则
	NetworkNotReachable                  = 0x03 // 网络不可达错误
	HostUnreachable                      = 0x04 // 目标主机不可访问
	ConnectionRefused                    = 0x05 // 连接目标主机被拒绝
	TTLExpired                           = 0x06 // 访问超时
	CommandNotSupported                  = 0x07 // 代理请求的`CMD`不支持
	AddressTypeNotSupported              = 0x08 // 代理请求的目标地址类型不支持
)

// Socks5 代理请求响应报文
// +----+-----+-------+------+----------+----------+
// |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
// +----+-----+-------+------+----------+----------+
// | 1  |  1  | X'00' |  1   | Variable |    2     |
// +----+-----+-------+------+----------+----------+
// VER: 协议版本，代理流程应保持一致
// REP: 响应码，告诉客户端本次请求的结果
// RSV: 保留位
// ATYP: `BND.ADDR` 地址类型
// BND.ADDR: 目标服务地址
// BND.PORT: 目标服务端口
type requestReply struct {
	Version  byte
	Reply    RequestReply
	Rsv      byte
	Type     AddressType
	BindAddr net.IP
	BindPort uint16
}
