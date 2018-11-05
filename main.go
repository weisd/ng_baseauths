package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/smtp"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/robfig/cron"
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

	InitConfig(confPath)

	fmt.Println(Conf)

	do()

	cron := cron.New()
	cron.AddFunc(Conf.CronEntry, do)
	cron.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			cron.Stop()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}

}

func do() {

	user := fmt.Sprintf("uc%d", rand.Intn(60))
	pass := RandStringBytesMaskImprSrc(6)

	authBuff := genBaseAuth(user, pass)

	log.Println("up:", user, pass, authBuff.String())

	body := fmt.Sprintf("%s %s", user, pass)

	// SendToMail
	if err := SendToMail("大婶，还学得动吗？", body); err != nil {
		log.Println(errors.Wrap(err, "SendToMail"))
		return
	}

	// 更新配置文件失败发邮件
	if err := ioutil.WriteFile(Conf.AuthFile, authBuff.Bytes(), 0644); err != nil {
		SendToMail("大婶，密码更新失败了", body)
	}

}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// RandStringBytesMaskImprSrc RandStringBytesMaskImprSrc
func RandStringBytesMaskImprSrc(n int) string {
	var src = rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func genBaseAuth(user, pass string) *bytes.Buffer {
	b := new(bytes.Buffer)
	b.WriteString(user)
	b.WriteString(":")
	b.WriteString("{SHA}")
	b.WriteString(GetSha(pass))
	// b.WriteString("\n")
	return b
}

// GetSha GetSha
func GetSha(password string) string {
	s := sha1.New()
	s.Write([]byte(password))
	passwordSum := []byte(s.Sum(nil))
	return base64.StdEncoding.EncodeToString(passwordSum)
}

// SendToMail SendToMail
func SendToMail(subject, body string) error {
	if len(Conf.Emails) == 0 {
		return nil
	}
	hp := strings.Split(Conf.SMTP.Host, ":")
	auth := smtp.PlainAuth("", Conf.SMTP.User, Conf.SMTP.Pass, hp[0])

	msg := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s\r\n", Conf.Emails[0], subject, body)

	return smtp.SendMail(Conf.SMTP.Host, auth, Conf.SMTP.User, Conf.Emails, []byte(msg))
}
