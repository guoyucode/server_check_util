package code

import (
	"log"
	"net/smtp"
	"strings"
	"time"
)

var (
	//MailFrom 发送邮件邮箱
	MailFrom = "1234556@qq.com"
	//MailFromName 发送邮件的名字
	MailFromName = "1234556"
	//MailPassword 邮箱密码
	MailPassword = "password"
	//MailHost 发送邮件的地址
	MailHost = "smtp.mxhichina.com:25"
	//MailTo 发送给哪些人
	MailTo = "1234556@qq.com;76543@qq.com"
)

//SendMail 发送邮件方法
func SendMail(subject, body string) {
	m := mailInfo{
		subject: subject,
		body:    body,
	}
	msgs <- m
}

type mailInfo struct {
	subject string
	body    string
}

var msgs = make(chan mailInfo, 500)

func init() {
	go func1()
}

//私有方法
func func1() {
	for true {
		timeout := make(chan bool, 1)
		select {
		case msg := <-msgs:
			go func() {
				time.Sleep(time.Second * 10)
				timeout <- true
			}()
			func2(msg.subject, msg.body)
		case <-timeout:
			log.Fatal("发送邮件超时")
		}
		time.Sleep(time.Second * 10)
	}
}

//私有方法
func func2(subject, body string) {
	subject = MailFromName + "(" + subject + ")" + time.Now().Format("_01-02 15:04:05.999")
	hp := strings.Split(MailHost, ":")
	auth := smtp.PlainAuth("", MailFrom, MailPassword, hp[0])
	var contentType string
	contentType = "Content-Type: text/" + "html" + "; charset=UTF-8"
	msg := []byte("To: " + MailTo + "\r\nFrom: " + MailFromName + "<" + MailFrom + ">\r\nSubject: " + subject + "\r\n" + contentType + "\r\n\r\n" + body)
	sendTo := strings.Split(MailTo, ";")
	if err := smtp.SendMail(MailHost, auth, MailFrom, sendTo, msg); err != nil {
		log.Fatal("发送邮件失败", err)
	}
}
