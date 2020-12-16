package configs

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

// 函数验证
type Option func(*Config)

// 默认系统
var defaultConfig = &Config{IsDebug: false, ServerConfigs: make([]*ServerConfig, 0, 2)}

//Nginx nginx  配置
type ServerConfig struct {
	Addr string `yaml:"Addr"`
	Port int    `yaml:"Port"`
}

//Config   系统配置配置
type Config struct {
	IsDebug       bool
	ServerConfigs []*ServerConfig
}

// 加载配置
func SetConfiguration(option ...Option) {
	for _, v := range option {
		v(defaultConfig)
	}
}

// 配置设置是否是debug[可选配置]
func SetIsDebug(isDebug bool) Option {
	return func(c *Config) {
		c.IsDebug = isDebug
	}
}

// 配置设置是否是debug[实际可选配置，仅举例]
func SetServerConfigs(configs []*ServerConfig) Option {

	if len(configs) < 2 {
		panic("配置至少需要两个")
	}

	return func(c *Config) {
		c.ServerConfigs = configs
	}
}

// httpserver的配置
func GetHttpConfig() *ServerConfig {
	return defaultConfig.ServerConfigs[0]
}

// grpcserver的配置
func GetGrpcConfig() *ServerConfig {
	return defaultConfig.ServerConfigs[1]
}

// 加载配置文件
func LoadYamlConfig() []*ServerConfig {
	r := make([]*ServerConfig, 0, 2)
	yf, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		panic(fmt.Sprintf("配置文件加载错误，%s", err))
	}

	err = yaml.Unmarshal(yf, &r)
	if err != nil {
		panic(fmt.Sprintf("配置文件解析错误，%s", err))
	}

	return r
}
