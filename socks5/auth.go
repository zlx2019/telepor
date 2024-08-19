package socks5

import (
	"errors"
	"fmt"
	"io"
	"telepor/connection"
	"telepor/logger"
	"telepor/pool"
)

// @Title       auth.go
// @Description Socks5 认证
// @Author      Zero.
// @Create      2024-08-19 15:25

// AuthUserRequest UsernamePassword 认证请求报文
// +----+------+----------+------+----------+
// |VER | ULEN |  UNAME   | PLEN |  PASSWD  |
// +----+------+----------+------+----------+
// | 1  |  1   | 1 to 255 |  1   | 1 to 255 |
// +----+------+----------+------+----------+
// VER: 协商版本，与Socks5版本无关，通常为 `0x01`
// ULEN: 用户名长度
// UNAME: 用户名
// PLEN: 密码长度
// PASSWD: 密码
type AuthUserRequest struct {
	Version     byte
	UsernameLen byte
	Username       string
	PasswdLen   byte
	Password      string
}



// AuthByUsernamePassword 以 UserPassword 模式与客户端进行认证
// 认证消息报文
// +----+------+----------+------+----------+
// |VER | ULEN |  UNAME   | PLEN |  PASSWD  |
// +----+------+----------+------+----------+
// | 1  |  1   | 1 to 255 |  1   | 1 to 255 |
// +----+------+----------+------+----------+
func (s *Server) AuthByUsernamePassword(c *connection.Connection) (req *AuthUserRequest ,err error) {
	// 响应握手报文，告知客户端使用 用户名密码认证
	err = s.shakeHandsReply(c, UserPassword)
	if err != nil {
		return
	}
	// 分配缓冲区
	buf := pool.Borrow(255)
	defer pool.Revert(buf)

	req = &AuthUserRequest{}
	// Read VER、ULEN
	if _, err = io.ReadFull(c, buf[:2]); err != nil {
		return nil, fmt.Errorf("read `VER` `ULEN` failed")
	}
	req.Version = buf[0]
	req.UsernameLen = buf[1]
	if req.UsernameLen <= 0 {
		return nil, errors.New("invalid username length")
	}
	// 读取用户名
	_, err = io.ReadFull(c, buf[:req.UsernameLen])
	if err != nil {
		return nil, errors.New("cannot read `UNAME`")
	}
	req.Username = string(buf[:req.UsernameLen])
	// 读取密码
	_, err = io.ReadFull(c, buf[:1])
	if err != nil {
		return
	}
	req.PasswdLen = buf[0]
	if req.PasswdLen <= 0 {
		return nil, errors.New("invalid password length")
	}
	_, err = io.ReadFull(c, buf[:req.PasswdLen])
	if err != nil {
		return
	}
	req.Password = string(buf[:req.PasswdLen])
	logger.Logger.InfoSf("[Socks5] Auth Version: %d Username: %s, Password: %s",req.Version, req.Username, req.Password)
	return
}

func (s *Server) AuthMessage()  {

}

// AuthReply 认证结果响应
// +----+--------+
// |VER | STATUS |
// +----+--------+
// | 1  |   1    |
// +----+--------+
// VER: 认证版本
// STATUS: 认证结果 `0x00`表示认证成功，其他表示为失败
func (s *Server) AuthReply(c *connection.Connection ,status byte) error  {
	_, err := c.Write([]byte{AuthVersion, status})
	return err
}