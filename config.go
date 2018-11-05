package main

import "github.com/BurntSushi/toml"

// Conf Conf
var Conf Config

// Config Config
type Config struct {
	AuthFile string   // 目录配置文件
	Emails   []string // 发送的email地址
	SMTP     SMTP
}

// SMTP SMTP
type SMTP struct {
	Addr string
	Port int
	User string
	Pass string
}

// InitConfig InitConfig
func InitConfig(fpath string) {
	var err error
	if _, err = toml.DecodeFile(fpath, &Conf); err != nil {
		panic(err)
	}
}
