package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"strings"
	"time"
)

const Email_User = "111111111"
const Email_Password = "1111111111111111"
const Email_Host = "smtphz.qiye.163.com"
const Email_Port = "465"
const Email_Name = "1111111111111"
const Email_Subject = "111111111111111111！"

type Mail struct {
	senderId string
	toIds    []string
	subject  string
	body     string
}

type SmtpServer struct {
	host string
	port string
}

func EmailSendCode(nickname, to, code string) error {

	if nickname != "" {
		nickname = nickname + "，"
	}

	body := `
        <html>
        <body>
        <h3>
        ` + nickname + `您好：
        </h3>
        非常感谢您使用` + Email_Name + `，您的邮箱验证码为：<br/>
        <b>` + code + `</b><br/>
        此验证码有效期30分钟，请妥善保存。<br/>
        如果这不是您本人的操作，请忽略本邮件。<br/>
        </body>
        </html>
        `

	return SendToMail(to, Email_Subject, body)
}

func (s *SmtpServer) ServerName() string {
	return s.host + ":" + s.port
}

func (mail *Mail) BuildMessage() string {
	message := ""
	message += fmt.Sprintf("From: %s<%s>\r\n", Email_Name, mail.senderId)
	if len(mail.toIds) > 0 {
		message += fmt.Sprintf("To: %s\r\n", strings.Join(mail.toIds, ";"))
	}

	message += fmt.Sprintf("Subject: %s\r\n", mail.subject)

	message += "Content-Type: text/html; charset=UTF-8"

	message += "\r\n\r\n" + mail.body

	return message
}

func SendToMail(to, subject, body string) error {

	mail := Mail{}
	mail.senderId = Email_User
	mail.toIds = strings.Split(to, ";")
	mail.subject = Email_Subject
	mail.body = body

	messageBody := mail.BuildMessage()

	smtpServer := SmtpServer{host: Email_Host, port: Email_Port}

	//build an auth
	auth := smtp.PlainAuth("", mail.senderId, Email_Password, smtpServer.host)

	// Gmail will reject connection if it's not secure
	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         smtpServer.host,
	}

	conn, err := tls.Dial("tcp", smtpServer.ServerName(), tlsconfig)
	if err != nil {
		return err
	}

	client, err := smtp.NewClient(conn, smtpServer.host)
	if err != nil {
		return err
	}

	// step 1: Use Auth
	if err = client.Auth(auth); err != nil {
		return err
	}

	// step 2: add all from and to
	if err = client.Mail(mail.senderId); err != nil {
		return err
	}
	for _, k := range mail.toIds {
		if err = client.Rcpt(k); err != nil {
			return err
		}
	}

	// Data
	w, err := client.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(messageBody))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	client.Quit()

	log.Println("Mail sent successfully")

	return nil
}

func tcpGather(ip string, ports []string, body string, to string) map[string]string {
	// 检查 emqx 1883, 8083, 8080, 18083 端口

	results := make(map[string]string)
	for _, port := range ports {
		address := net.JoinHostPort(ip, port)
		// 3 秒超时
		conn, err := net.DialTimeout("tcp", address, 3*time.Second)
		if err != nil {
			results[port] = "failed"
			SendToMail(to, "subject", body)
		} else {
			if conn != nil {
				results[port] = "success"
				_ = conn.Close()
			} else {
				results[port] = "failed"
				SendToMail(to, "subject", body)

			}
		}
	}
	return results
}

func main() {
	HOST := "127.0.0.1"
	ports := []string{"8899"}
	body := "银企对接前置机8899端口检测未启动，请检查"
	to := "m111111111111cn"
	SendToMail(to, "subject", body)
	results := tcpGather(HOST, ports, body, to)
	fmt.Println(results)

}
