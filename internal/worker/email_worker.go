package worker

import (
	"log"
	"time"
	"github.com/Full-finger/OIDC/internal/util"
)

// EmailWorker 邮件工作者
type EmailWorker struct {
	emailQueue util.EmailQueue
	emailService util.EmailService
	running    bool
}

// NewEmailWorker 创建邮件工作者实例
func NewEmailWorker(emailQueue util.EmailQueue, emailService util.EmailService) *EmailWorker {
	return &EmailWorker{
		emailQueue:   emailQueue,
		emailService: emailService,
		running:      false,
	}
}

// Start 启动邮件工作者
func (w *EmailWorker) Start() {
	w.running = true
	log.Println("邮件工作者已启动")

	// 持续监听队列并处理邮件任务
	for w.running {
		// 从队列中获取邮件任务
		item, err := w.emailQueue.Dequeue()
		if err != nil {
			// 队列为空，等待一段时间再尝试
			time.Sleep(1 * time.Second)
			continue
		}

		// 处理邮件任务
		if err := w.processEmail(item); err != nil {
			log.Printf("处理邮件任务失败: %v, 邮箱: %s\n", err, item.Email)
			// 在实际项目中，可能需要实现重试机制或者将失败的任务放入死信队列
			// 这里简化处理，仅记录日志
		} else {
			log.Printf("邮件已发送至: %s\n", item.Email)
		}
	}
}

// Stop 停止邮件工作者
func (w *EmailWorker) Stop() {
	w.running = false
	log.Println("邮件工作者已停止")
}

// processEmail 处理邮件任务
func (w *EmailWorker) processEmail(item *util.EmailQueueItem) error {
	// 调用邮件服务发送验证邮件
	return w.emailService.SendVerificationEmail(item.Email, item.Token)
}