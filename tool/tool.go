// @Title tool.go
// @Description $END$
// @Author Zero - 2024/8/11 01:00:18

package tool

import (
	"strings"
	"telepor/connection"
	"telepor/define"
)

var methods = []string{"CONNECT", "GET", "POST", "PUT", "DELETE", "TRACE", "OPTIONS"}

// IdentifyPact 判断并且获取连接使用的协议类型
func IdentifyPact(c *connection.Connection) (int8, error) {
	// is Socks5？
	ver, _ := c.Peek(1)
	if ver[0] == define.SocksVersion {
		return define.Socks5, nil
	}

	// is HTTP? 通过HTTP请求行 Method 判断
	line, err := c.Peek(7)
	if err != nil {
		return 0, err
	}
	method := strings.ToUpper(string(line))
	isHttp := ContainsBy(methods, func(item string) bool {
		return strings.HasPrefix(method, item)
	})
	if isHttp {
		return define.HTTP, nil
	}
	return define.Unknown, nil
}

func Contains[T comparable](slice []T, target T) bool {
	for _, item := range slice {
		if item == target {
			return true
		}
	}
	return false
}

func ContainsBy[T any](slice []T, predicate func(T) bool) bool {
	for _, item := range slice {
		if predicate(item) {
			return true
		}
	}
	return false
}
