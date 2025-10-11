package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"github.com/Full-finger/OIDC/internal/worker"
	"github.com/Full-finger/OIDC/internal/util"
)

func main() {
	log.Println("启动邮件工作者服务...")

	// 初始化依赖
	emailQueue := util.NewSimpleEmailQueue()
	emailService := util.NewEmailService()
	
	// 创建邮件工作者
	emailWorker := worker.NewEmailWorker(emailQueue, emailService)
	
	// 启动邮件工作者（在goroutine中运行）
	go emailWorker.Start()
	
	// 等待中断信号以优雅地关闭工作者
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	
	// 停止邮件工作者
	emailWorker.Stop()
	log.Println("邮件工作者服务已停止")
}