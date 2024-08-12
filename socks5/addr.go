package socks5

import (
	"net"
	"strconv"
)

// @Title       addr.go
// @Description 地址解析
// @Author      Zero.
// @Create      2024-08-12 13:06

// Addr 表示 Socks5 协议中的地址信息 [ATYP + ADDR + PORT]
// +------+----------+----------+
// | ATYP | DST.ADDR | DST.PORT |
// +------+----------+----------+
// |  1   | Variable |    2     |
// +------+----------+----------+
type Addr []byte

// 将 Addr 地址转换为字符串格式
func (addr Addr) String() string {
	var host, port string
	// ADDR可能为 [IPv4 | IPv6 | 域名]
	switch addr[0] {
	case IPv4:
		// 网络传输默认为大端，强转为net.IP即可转为小端
		host = net.IP(addr[1 : 1+net.IPv4len]).String()
		port = strconv.Itoa((int(addr[1+net.IPv4len]) << 8) | int(addr[1+net.IPv4len+1]))
	case IPv6:
		host = net.IP(addr[1 : 1+net.IPv6len]).String()
		port = strconv.Itoa((int(addr[1+net.IPv6len]) << 8) | int(addr[1+net.IPv6len+1]))
	case Domain:
		host = string(addr[2 : 2+int(addr[1])])
		port = strconv.Itoa((int(addr[2+int(addr[1])]) << 8) | int(addr[2+int(addr[1])+1]))
	}
	return net.JoinHostPort(host, port)
}
