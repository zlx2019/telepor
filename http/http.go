package http

import "telepor/connection"

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

func (s *Server) Handle(conn *connection.Connection) error {
	// 从TCP中读取 HTTP1.x Request
	//request, err := http.ReadRequest(connection.Reader())
	return nil
}