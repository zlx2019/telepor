package pool

import (
	"bytes"
	"sync"
)

// @Title       bytes_pool.go
// @Description bytes 动态缓冲池
// @Author      Zero.
// @Create      2024-08-19 18:24

const limit = 64 <<10

// bytes.Buffer 动态字节缓冲池
var bytesPools = sync.Pool{New: func() any {return &bytes.Buffer{}}}

// BorrowBytesBuf 获取缓冲区
func BorrowBytesBuf() *bytes.Buffer {
	return bytesPools.Get().(*bytes.Buffer)
}

// RevertBytesBuf 重置并归还缓冲区
func RevertBytesBuf(buf *bytes.Buffer) {
	// 对较小的缓冲区进行复用.
	if buf.Cap() <= limit {
		buf.Reset()
		bytesPools.Put(buf)
	}
}