// @Title define.go
// @Description Socks5 服务相关常量
// @Author Zero - 2024/8/11 22:00:34

package socks5

// Version 协议版本
const (
	Version = 0x05
)

// Method Socks5 认证方式类型
type Method = byte

const (
	NoAuthentication    Method = 0x00 // 无需认证
	GSSAPI                     = 0x01 // GSSAPI认证
	UserPassword               = 0x02 // 用户名密码认证
	ReservedIana               = 0x03 // [0x03 ~ 0x7F] 范围为预留的认证方式(由 IANA 预留)
	ReservedCustom             = 0x80 // [0x80 ~ 0xFE] 范围为用户自定义认证方式
	NoAcceptableMethods        = 0xFF // 没有可接受的认证方式.
)

// Command Socks5 代理请求指令
type Command = byte

const (
	Connect Command = 0x01 // 表示TCP代理
	BIND            = 0x02
	UDP             = 0x03 // 表示UDP代理
)

// AddressType 地址类型
type AddressType = byte

const (
	IPv4   AddressType = 0x01 // 表示为IPv4
	Domain             = 0x03 // 表示为域名
	IPv6               = 0x04 // 表示为IPv6
)

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
