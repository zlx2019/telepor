package http

import . "telepor/connection"

// @Title       http.go
// @Description HTTP 代理服务
// @Author      Zero.
// @Create      2024-08-09 18:34

// Server HTTP 服务端
type Server struct {
	Next         bool		// 是否转发到下级代理
	NextProtocol string		// 下级代理协议(支持 HTTP -> Socks5)
	NextHost     string		// 下级代理服务地址
	NextPort     uint16		// 下级代理服务端口
}

func NewHttpServer() *Server {
	return &Server{}
}

// Startup 启动 HTTP 代理服务器
func (s *Server) Startup() error {
	//TODO implement me
	panic("implement me")
}

// ServeHandle 处理 HTTP/S 请求
func (s *Server) ServeHandle(c *Connection) {
	// 从TCP中读取 HTTP1.x Request
	//request, err := http.ReadRequest(connection.Reader())
}
