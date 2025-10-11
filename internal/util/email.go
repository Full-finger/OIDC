package util

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

// EmailService 邮件服务接口
type EmailService interface {
	// SendVerificationEmail 发送验证邮件
	SendVerificationEmail(email, token string) error
}

// emailService 邮件服务实现
type emailService struct {
	// 可以添加邮件服务器配置等依赖
}

// NewEmailService 创建邮件服务实例
func NewEmailService() EmailService {
	return &emailService{}
}

// SendVerificationEmail 发送验证邮件
func (e *emailService) SendVerificationEmail(email, token string) error {
	// 模拟邮件发送过程
	// 在实际项目中，这里会连接到真实的邮件服务器发送邮件
	fmt.Printf("发送验证邮件到: %s\n", email)
	fmt.Printf("验证链接: http://localhost:8080/verify?token=%s\n", token)
	
	// 模拟网络延迟
	rand.Seed(time.Now().UnixNano())
	delay := time.Duration(rand.Intn(1000)) * time.Millisecond
	time.Sleep(delay)
	
	log.Printf("邮件已发送到 %s，耗时 %v", email, delay)
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