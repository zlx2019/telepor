// @Title 		server.go
// @Description 服务端
// @Author Zero - 2024/8/8 23:44:14

package server

import "telepor/connection"

// Server 服务端
type Server interface {
	// Startup 启动服务
	Startup() error
}

// ProxyServer 代理服务端
type ProxyServer interface {
	// ServeHandle 代理请求处理
	ServeHandle(*connection.Connection)
}

// MiddleProxyServer 代理转发服务端
type MiddleProxyServer interface {
}
