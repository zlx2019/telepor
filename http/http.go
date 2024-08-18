package http

import . "telepor/connection"

// @Title       http.go
// @Description HTTP 代理服务
// @Author      Zero.
// @Create      2024-08-09 18:34

// Server HTTP 服务端
type Server struct {

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