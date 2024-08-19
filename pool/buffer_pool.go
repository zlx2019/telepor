package pool

// @Title       buffer_pool.go
// @Description 缓冲区池
// @Author      Zero.
// @Create      2024-08-12 16:15

import (
	"math/bits"
	"sync"
)

// 按照级别(1~16)创建一批缓冲池，缓解内存的分配频率
// level = (1 ~ 16)
// byte = (1 ~ 65536)

const (
	level    = 17
	maxBytes = 1 << (level - 1)
)

var (
	levels      [level]int       // 不同缓冲池中的缓冲区大小
	bufferPools [level]sync.Pool // 缓冲池序列
)

// init bufferPools
func init() {
	for l := 0; l < level; l++ {
		levelCap := 1 << l
		levels[l] = levelCap
		bufferPools[l].New = func() any {
			//logger.Logger.InfoSf("Creating new buffer pool level: %d cap: %d", l, levelCap)
			return make([]byte, levelCap)
		}
	}
}

// Borrow 从池中获取一个大于size的最小缓冲区,超过[65536]字节则临时分配.
func Borrow(size int) []byte {
	if size >= 1 && size <= maxBytes {
		l := bits.Len32(uint32(size)) - 1
		if levels[l] < size {
			l += 1
		}
		return bufferPools[l].Get().([]byte)[:size]
	}
	return make([]byte, size)
}

// Revert 将缓冲区归还给缓冲池,重复利用
func Revert(buffer []byte) {
	if size := cap(buffer); size >= 1 && size <= maxBytes {
		l := bits.Len32(uint32(size)) - 1
		if levels[l] == size {
			bufferPools[l].Put(buffer)
		}
	}
}
