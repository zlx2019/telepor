// @Title define.go
// @Description 全局常量
// @Author Zero - 2024/8/10 17:53:26

package define

// SocksVersion Socks协议版本
const SocksVersion = 0x05

// ProtocolType 支持的代理协议类型
type ProtocolType = int8

const (
	Unknown ProtocolType = iota - 1
	Socks5
	HTTP
	HTTPS
)
