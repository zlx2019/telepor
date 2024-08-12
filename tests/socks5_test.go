// @Title socks5_test.go
// @Description $END$
// @Author Zero - 2024/8/11 14:18:17

package tests

import (
	"fmt"
	"testing"
)

// Socks5 unit testing
func TestSocks5Conn(t *testing.T) {
	//t.Run("parse socks5 Negotiate", func(t *testing.T) {
	//	b := []byte{define.Socks5Ver, 0x02, 0x00, 0x01}
	//	socks5Server := server.Socks5Server{}
	//	pack, err := socks5Server.UnpackAuthMethod(bytes.NewReader(b))
	//	if err != nil {
	//		t.Fatalf("parse fail: %s", err)
	//	}
	//	t.Logf("success: %v", pack)
	//})
}

func TestBytes(t *testing.T) {
	for i := 0; i < 17; i++ {
		fmt.Println(1 << i)
	}
}
