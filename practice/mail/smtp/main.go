package main

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"strings"
)

func main() {
	// 配置邮件信息
	email := Email{
		From:    "发信地址",
		To:      []string{"收件地址"},
		Subject: "Test Email from Go",
		Body:    "<h1>Hello</h1>",
	}

	// SMTP 服务器配置
	smtpServer := "smtpdm.aliyun.com"
	smtpPort := "80" // 或者 "25", "465", "587" 根据你的配置选择
	smtpUser := "发信地址"
	smtpPass := "SMTP密码"

	// 发送邮件
	err := sendEmail(email, smtpServer, smtpPort, smtpUser, smtpPass)
	if err != nil {
		log.Fatalf("Failed to send email: %v", err)
	}
	fmt.Println("Email sent successfully!")
}

// Email 邮件结构体
type Email struct {
	From    string
	To      []string
	Subject string
	Body    string
}

// 发送邮件函数
func sendEmail(email Email, smtpServer string, smtpPort string, smtpUser string, smtpPass string) error {
	// 构建邮件内容
	msg := bytes.NewBufferString("")
	msg.WriteString(fmt.Sprintf("From: %s\r\n", email.From))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(email.To, ",")))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", email.Subject))
	msg.WriteString("MIME-version: 1.0;\n")
	msg.WriteString("Content-Type: text/html; charset=\"UTF-8\";\n\n")
	msg.WriteString(email.Body)

	// SMTP 认证信息
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpServer)

	// SMTP 服务器地址
	serverAddr := fmt.Sprintf("%s:%s", smtpServer, smtpPort)

	// 发送邮件
	err := smtp.SendMail(serverAddr, auth, email.From, email.To, msg.Bytes())
	if err != nil {
		return err
	}
	return nil
}
