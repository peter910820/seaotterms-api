package blog

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

func GenerateRandomKey(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func SendRegisterEmail(to string, registerKey string) {
	registerUrl := fmt.Sprintf("%sapi/gal/register/%s/%s", os.Getenv("SITE_DOMAIN"), strings.Split(to, "@")[0], registerKey)

	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("MAIL_SEND_NOREPLY_ACCOUNT"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", "註冊驗證信件")
	m.SetBody("text/html", fmt.Sprintf(`			
	<body style="margin: 0; padding: 0; height: 100vh; display: flex; justify-content: center; align-items: center; position: relative; overflow: hidden;">
  <img src="https://fs.seaotterms.com/resource/image/premium_photo-1673029927496-a9b574091d12.jpg"
       style="position: absolute; width: 100%%; height: auto; top: 50%%; left: 50%%; transform: translate(-50%%, -50%%); opacity: 0.3; z-index: -1;">
  <div style="
    background-color: rgba(0, 0, 0, 0.6);
    padding: 40px;
    border-radius: 15px;
    color: white;
    text-align: center;
    font-family: sans-serif;
    text-shadow: 1px 1px 5px black;
  ">
    <h1>信箱註冊成功</h1>
    <p>您的信箱已經註冊成功!</p>
    <p><a href="%s" style="color: yellow; font-weight: bold; text-decoration: underline;">點擊這裡</a> 完成驗證。</p>
  </div>
</body>
	`, registerUrl))

	d := gomail.NewDialer(os.Getenv("MAIL_HOST"), 587, os.Getenv("MAIL_SEND_NOREPLY_ACCOUNT"), os.Getenv("MAIL_SEND_NOREPLY_PASSWORD"))
	d.TLSConfig = &tls.Config{
		ServerName: os.Getenv("MAIL_HOST"),
		MinVersion: tls.VersionTLS12,
	}

	if err := d.DialAndSend(m); err != nil {
		logrus.Error(err)
	} else {
		logrus.Infof("寄信成功! 目標: %s", to)
	}
}
