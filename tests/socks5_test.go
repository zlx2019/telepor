// @Title socks5_test.go
// @Description $END$
// @Author Zero - 2024/8/11 14:18:17

package tests

import (
	"bytes"
	"io"
	"net"
	"os"
	"reflect"
	"telepor/define"
	"telepor/socks5"
	"testing"
	"time"
)

// Socks5 协商报文解析
func TestNegotiationUnpack(t *testing.T) {
	t.Run("parse socks5 Negotiate", func(t *testing.T) {
		b := []byte{define.SocksVersion, 0x02, 0x00, 0x01}
		req, err := socks5.NegotiationUnpack(bytes.NewReader(b))
		if err != nil {
			t.Fatalf("parse fail: %s", err)
		}
		t.Logf("success: %v", req)
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

// Socks5 代理请求报文响应
func TestRequestSuccessReply(t *testing.T) {
	buffer := bytes.Buffer{}
	ip := net.IP([]byte{192, 168, 0, 1})
	err := socks5.RequestSuccessReply(&buffer, ip, 7890)
	if err != nil {
		t.Fatalf("error while writing : %s",err)
	}
	want := []byte{socks5.Version, socks5.Succeeded, socks5.Reserved, socks5.IPv4, 192, 168, 0, 1, 0x1e, 0xd2}
	got := buffer.Bytes()
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("message not match: want %v, got %v", want, got)
	}
}

func TestIoPipe(t *testing.T) {
	r, w := io.Pipe()
	go func() {
		time.Sleep(time.Second * 3)
		w.Write([]byte("Hello world\n"))
		time.Sleep(time.Second * 3)
		w.Close()
	}()
	_, _ = io.Copy(os.Stdout, r)
}