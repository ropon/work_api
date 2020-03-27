package config

import "fmt"

//定义错误代码说明
var CodeType = map[int]string{
	4001: "参数不完整",
	4002: "请求过于频繁",
	4003: "授权异常",
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

type MysqlConfig struct {
	MHost    string `conf:"host"`
	MPort    int64  `conf:"port"`
	UserName string `conf:"username"`
	PassWord string `conf:"password"`
	DBName   string `conf:"dbname"`
}

// Config 配置文件结构体
type Config struct {
	Global      `conf:"global"`
	LimitConfig `conf:"limit"`
	RedisConfig `conf:"redis"`
	MysqlConfig `conf:"mysql"`
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

func GetDSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		Cfg.UserName,
		Cfg.PassWord,
		Cfg.MHost,
		Cfg.MPort,
		Cfg.DBName,
	)
}
