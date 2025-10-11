package util

import (
	"fmt"
	"log"
	"math/rand"
	"time"
	"net/smtp"
	"os"
)

// EmailService 邮件服务接口
type EmailService interface {
	// SendVerificationEmail 发送验证邮件
	SendVerificationEmail(email, token string) error
}

// emailService 邮件服务实现
type emailService struct {
	smtpHost     string
	smtpPort     string
	senderEmail  string
	senderPassword string
}

// NewEmailService 创建邮件服务实例
func NewEmailService() EmailService {
	return &emailService{
		smtpHost:       getEnv("SMTP_HOST", "smtp.gmail.com"),
		smtpPort:       getEnv("SMTP_PORT", "587"),
		senderEmail:    getEnv("SENDER_EMAIL", "noreply@example.com"),
		senderPassword: getEnv("SENDER_PASSWORD", ""),
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// SendVerificationEmail 发送验证邮件
func (e *emailService) SendVerificationEmail(email, token string) error {
	// 邮件主题
	subject := "请验证您的邮箱地址"
	
	// 邮件内容
	verificationURL := fmt.Sprintf("http://localhost:8080/api/v1/verify?token=%s", token)
	body := fmt.Sprintf("请点击以下链接验证您的邮箱地址：\n%s\n\n如果您没有注册我们的服务，请忽略此邮件。", verificationURL)
	
	// 构建邮件
	message := fmt.Sprintf(
		"From: %s\n"+
			"To: %s\n"+
			"Subject: %s\n"+
			"\n"+
			"%s",
		e.senderEmail,
		email,
		subject,
		body,
	)
	
	// 发送邮件
	auth := smtp.PlainAuth("", e.senderEmail, e.senderPassword, e.smtpHost)
	err := smtp.SendMail(e.smtpHost+":"+e.smtpPort, auth, e.senderEmail, []string{email}, []byte(message))
	if err != nil {
		log.Printf("发送邮件失败: %v", err)
		return err
	}
	
	log.Printf("验证邮件已发送到 %s", email)
	return nil
}

// EmailQueueItem 邮件队列项
type EmailQueueItem struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

// EmailQueue 邮件队列接口
type EmailQueue interface {
	// Enqueue 将邮件任务加入队列
	Enqueue(item EmailQueueItem) error
	// Dequeue 从队列中取出邮件任务
	Dequeue() (*EmailQueueItem, error)
}

// SimpleEmailQueue 简单邮件队列实现（使用内存队列模拟）
type SimpleEmailQueue struct {
	queue []EmailQueueItem
}

// NewSimpleEmailQueue 创建简单邮件队列实例
func NewSimpleEmailQueue() EmailQueue {
	return &SimpleEmailQueue{
		queue: make([]EmailQueueItem, 0),
	}
}

// Enqueue 将邮件任务加入队列
func (q *SimpleEmailQueue) Enqueue(item EmailQueueItem) error {
	q.queue = append(q.queue, item)
	fmt.Printf("邮件任务已加入队列: %s\n", item.Email)
	return nil
}

// Dequeue 从队列中取出邮件任务
func (q *SimpleEmailQueue) Dequeue() (*EmailQueueItem, error) {
	if len(q.queue) == 0 {
		return nil, fmt.Errorf("队列为空")
	}
	
	item := q.queue[0]
	q.queue = q.queue[1:]
	return &item, nil
}