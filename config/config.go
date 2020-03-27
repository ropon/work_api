package config

import "fmt"

//定义错误代码说明
var CodeType = map[int]string{
	4001: "参数不完整",
	4002: "请求过于频繁",
}

var Cfg *Config

type Global struct {
	Host string `conf:"host"`
	Port int64  `conf:"port"`
}

type LimitConfig struct {
	Count  int `conf:"count"`
	Second int `conf:"second"`
}

type RedisConfig struct {
	Addr     string `conf:"addr"`
	RPort    int64  `conf:"port"`
	Password string `conf:"password"`
	DB       int    `conf:"db"`
}

// Config 配置文件结构体
type Config struct {
	Global      `conf:"global"`
	LimitConfig `conf:"limit"`
	RedisConfig `conf:"redis"`
}

func checkConf() {

}

func InitConf(confName string) (err error) {
	Cfg = &Config{}
	err = parseConf(confName, Cfg)
	checkConf()
	return
}

func GHostPort() string {
	return fmt.Sprintf("%s:%d", Cfg.Host, Cfg.Port)
}

func RHostPort() string {
	return fmt.Sprintf("%s:%d", Cfg.Addr, Cfg.RPort)
}
