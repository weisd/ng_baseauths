package main

import (
	"flag"

	gomail "gopkg.in/gomail.v2"
)

var (
	confPath string
)

func init() {
	flag.StringVar(&confPath, "c", "./config.toml", "-c config.toml")
}

// 定时更新base_auth密码文件，并发送email
func main() {

	flag.Parse()

	InitConfig()

}

// SendEmail SendEmail
func SendEmail(content string) {

	d := gomail.NewDialer(Conf.SMTP.Addr, Conf.SMTP.Port, Conf.SMTP.User, Conf.SMTP.Pass)

}
