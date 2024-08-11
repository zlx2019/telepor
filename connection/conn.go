// @Title conn.go
// @Description Connect 封装
// @Author Zero - 2024/8/11 09:57:54

package connection

import (
	"bufio"
	"io"
	"net"
)

const (
	TCPBufSize = 32 << 10
	UDPBufSize = 2 << 10
)
// Connection 扩展连接,额外增加一个缓冲区
type Connection struct {
	buffer *bufio.Reader
	net.Conn
}

func WrapConn(c net.Conn) *Connection {
	if conn, ok := c.(*Connection); ok {
		return conn
	}
	return &Connection{bufio.NewReader(c), c}
}

func (c *Connection) Reader() *bufio.Reader {
	return c.buffer
}

func (c *Connection) Read(dst []byte) (int,error) {
	return c.buffer.Read(dst)
}

func (c *Connection) WriteTo(w io.Writer) (int64,error) {
	return c.buffer.WriteTo(w)
}

func (c *Connection) Peek(n int) ([]byte,error) {
	return c.buffer.Peek(n)
}

func (c *Connection) Close() error {
	return c.Conn.Close()
}