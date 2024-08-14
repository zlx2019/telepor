package http

// @Title       request.go
// @Description HTTP Request
// @Author      Zero.
// @Create      2024-08-14 10:38
import (
	"io"
)

// Request 表示 HTTP 请求信息
type Request struct {
}

// 从流中读取字节流，解析为 Request 请求体.
func unpackRequest(reader io.Reader) (*Request, error) {
	return nil, nil
}
