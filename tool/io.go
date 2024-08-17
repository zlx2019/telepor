// @Title io.go
// @Description IO 相关工具函数
// @Author Zero - 2024/8/17 17:47:46

package tool

import (
	"errors"
	"io"
	"os"
	"sync"
)

// Swap 交换两个连接的数据流，并计算发送与响应的数据总量
func Swap(client, server io.ReadWriter) (n int64 ,err error) {
	var serr, rerr error
	var latch sync.WaitGroup
	var sent int64 // client 发送数据量
	var rece int64 // server 响应数据量
	latch.Add(2)
	go func() {
		// server -> client
		rece, rerr = Transfer(client, server)
		latch.Done()
	}()
	go func() {
		// client -> server
		sent, serr = Transfer(server, client)
		latch.Done()
	}()
	// 等待两个连接关闭.
	latch.Wait()
	n = sent + rece
	if serr != nil && !errors.Is(serr, os.ErrDeadlineExceeded){
		err = serr
	}
	if rerr != nil && !errors.Is(rerr, os.ErrDeadlineExceeded){
		err = rerr
	}
	return
}

// Transfer 转移数据流，并允许指定内部缓冲区大小
func Transfer(dst io.Writer, src io.Reader, sizes ...uint64) (int64,error) {
	if len(sizes) > 0 {
		return transfer(dst, src, make([]byte, sizes[0]))
	}
	return transfer(dst, src, nil)
}
func TransferWithBuf(dst io.Writer, src io.Reader, buffer []byte) (int64,error) {
	return transfer(dst, src, buffer)
}

// TransferOnlyN 数据流拷贝，仅从src中拷贝 n 字节.
func TransferOnlyN(dst io.Writer, src io.Reader, n int64, size ...uint64) (int64,error){
	if len(size) > 0 {
		return TransferOnlyNWithBuf(dst, src, n, make([]byte, size[0]))
	}
	return TransferOnlyNWithBuf(dst, src, n, nil)
}

// TransferOnlyNWithBuf 数据传输，仅传输 n 个字节，并且使用外部提供的缓冲区
func TransferOnlyNWithBuf(dst io.Writer, src io.Reader, n int64, buffer []byte) (int64, error) {
	if n <= 0 {
		return 0, nil
	}
	wrapSrc := io.LimitReader(src, n)
	written, err := transfer(dst, wrapSrc, buffer)
	if written == n {
		// 拷贝指定的数据量完成
		return n, nil
	}
	if written < n && err == nil {
		// src 的所有数据已拷贝完
		err = io.EOF
	}
	return written, err
}


// 将 src 中的数据，转移到 dst 中
// 直至从 src 读取到 EOF，但正常完成后不会返回 EOF
func transfer(dst io.Writer, src io.Reader, buffer []byte) (written int64, err error) {
	// mark 无需通过 WriterTo 和 ReaderFrom 进行传输，产生不必要的调用栈开销
	if buffer == nil {
		// 默认缓冲区大小: 1MB
		size := MB
		if lo, ok := src.(*io.LimitedReader); ok && int64(size) > lo.N {
			if lo.N < 1 {
				size = 1
			} else {
				size = int(lo.N)
			}
		}
		buffer = make([]byte, size)
	}
	// 循环读取写入，直到读EOF || 写入的字节数 >= limit
	for {
		rn, rerr := src.Read(buffer)
		if rn > 0 {
			wn, werr := dst.Write(buffer[0:rn])
			if wn < 0 || rn < wn {
				wn = 0
				if werr == nil {
					werr = errors.New("invalid write result")
				}
			}
			// 写入传输量累加
			written += int64(wn)
			if werr != nil {
				// 写入错误
				err = werr
				break
			}
			if rn != wn {
				// 读取量和写入量不一致
				err = io.ErrShortWrite
				break
			}
		}
		if rerr != nil {
			if rerr != io.EOF {
				err = rerr
			}
			break
		}
	}
	return written, err
}
