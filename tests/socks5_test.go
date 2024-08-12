// @Title socks5_test.go
// @Description $END$
// @Author Zero - 2024/8/11 14:18:17

package tests

import (
	"bytes"
	"telepor/define"
	"telepor/socks5"
	"testing"
)

// Socks5 协商报文解析
func TestNegotiationUnpack(t *testing.T) {
	t.Run("parse socks5 Negotiate", func(t *testing.T) {
		b := []byte{define.SocksVersion, 0x02, 0x00, 0x01}
		pack, err := socks5.NegotiationUnpack(bytes.NewReader(b))
		if err != nil {
			t.Fatalf("parse fail: %s", err)
		}
		t.Logf("success: %v", pack)
	})
}

// Socks5 代理请求报文解析
func TestProxyRequestUnpack(t *testing.T) {
	var buf bytes.Buffer
	buf.Write([]byte{define.SocksVersion, socks5.Connect, 0x00, socks5.IPv4})
	buf.Write([]byte{192, 168, 0, 1})
	buf.Write([]byte{0x00, 0x50})
	requestMsg, err := socks5.ProxyRequestUnpack(&buf)
	if err != nil {
		t.Fatalf("testing fail: %v", err)
	}
	t.Logf("%v", requestMsg)
}
