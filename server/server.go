// @Title 		server.go
// @Description 服务端
// @Author Zero - 2024/8/8 23:44:14

package server

import "telepor/connection"

// Server 服务端
type Server interface {
	// Startup 启动服务
	Startup() error
	// ServeHandle 处理请求
	ServeHandle(*connection.Connection)
}

