// @Title errors.go
// @Description Socks5 代理错误
// @Author Zero - 2024/8/12 22:40:32

package socks5

import "errors"

// Socks5 协商阶段错误
var (
	VersionNotSupportedErr = errors.New("the current service only supports Socks5")
	CommandNotSupportedErr = errors.New("the current service only supports tcp proxy")
)

// Socks5 请求阶段错误
var ()
