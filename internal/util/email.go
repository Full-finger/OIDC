package util

import (
	"fmt"
	"net/smtp"
	"os"
)

// EmailService 邮件服务
type EmailService struct {
	smtpHost     string
	smtpPort     string
	smtpUser     string
	smtpPassword string
}

// NewEmailService 创建邮件服务实例
func NewEmailService(host, port, user, password string) *EmailService {
	return &EmailService{
		smtpHost:     host,
		smtpPort:     port,
		smtpUser:     user,
		smtpPassword: password,
	}

}

// SendVerificationEmail 发送验证邮件
func (es *EmailService) SendVerificationEmail(to, verificationURL string) error {
	auth := smtp.PlainAuth("", es.smtpUser, es.smtpPassword, es.smtpHost)

	subject := "Bangumoe 邮箱验证"
	body := fmt.Sprintf(`
欢迎注册 Bangumoe！

请点击以下链接验证您的邮箱：
%s

如果不是您本人操作，请忽略此邮件。
`, verificationURL)

	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body + "\r\n")

	return smtp.SendMail(es.smtpHost+":"+es.smtpPort, auth, es.smtpUser, []string{to}, msg)
}

// GetEnv 获取环境变量，如果不存在则使用默认值
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}