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

// Reserved 保留位，无任何意义的字节填充
const Reserved = 0x00


// 用户名密码认证状态
const(
	AuthVersion = 0x01
	AuthSuccess = 0x00
	AuthFailure = 0x01
)