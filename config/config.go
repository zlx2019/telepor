package config

import (
	"flag"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"telepor/logger"
)

// @Title config.go
// @Description Configs
// @Author Zero - 2024/8/11 09:57:54

var Conf *Config

// Config 服务配置
type Config struct {
	Bind   string       `yaml:"bind"`
	Socks5 Socks5Config `yaml:"socks5"`
	Http   HTTPConfig   `yaml:"http"`
}

// Socks5Config Socks5 server config
type Socks5Config struct {
	AuthMode byte
	Auths    map[string]string `yaml:"auths"`
}

// HTTPConfig HTTP server config
type HTTPConfig struct {
}

// AuthUser 认证用户信息
type AuthUser struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

const defaultConfigFile = "config.yml"

var configFile string

func init() {
	// parser commands
	flag.StringVar(&configFile, "c", defaultConfigFile, "server config file path.")
	flag.Parse()
	logger.Logger.InfoSf("config file path: %s", configFile)
	loadConfig()
}

// load config file
func loadConfig() {
	vp := viper.New()
	vp.SetConfigFile(configFile)
	vp.SetConfigType("yaml")
	err := vp.ReadInConfig()
	if err != nil {
		logger.Logger.FatalSf("failed to load the config file [%s]: %v", configFile, err)
	}
	err = vp.Unmarshal(&Conf)
	if err != nil {
		logger.Logger.FatalSf("failed to unmarshal the config file [%s]: %v", configFile, err)
	}

	// Watch config file
	vp.WatchConfig()
	vp.OnConfigChange(func(e fsnotify.Event) {
		if e.Has(fsnotify.Write) {
			if err := vp.Unmarshal(&Conf); err != nil {
				logger.Logger.ErrorSf("[Config Change] config reload failed: %v", err)
				return
			}
			logger.Logger.InfoSf("[Config Change] config reload success: %s", Conf)
		}
	})
}
