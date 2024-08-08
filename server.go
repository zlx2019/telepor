// @Title 		server.go
// @Description 服务端
// @Author Zero - 2024/8/8 23:44:14

package telepor

// ProxyServer 代理服务器
type ProxyServer interface {
	Run() error
}