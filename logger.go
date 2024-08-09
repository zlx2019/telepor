package main

import (
	"github.com/zlx2019/spoor"
	"go.uber.org/zap"
)

// @Title       logger.go
// @Description Logger组件
// @Author      Zero.
// @Create      2024-08-09 15:13

var Logger *spoor.Spoor

// 初始化日志组件
func init() {
	var err error
	Logger, err = spoor.NewSpoor(&spoor.Config{
		Level:         spoor.DEBUG,
		LogTimeFormat: "",
		Plugins:       []zap.Option{zap.AddCaller(), zap.AddStacktrace(spoor.ERROR)},
		WrapSkip:      1,
	})
	if err != nil {
		panic("init log failed.")
	}
}
