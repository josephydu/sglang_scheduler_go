package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sglang_scheduler_go/controller"
	"sglang_scheduler_go/server_args"
	"syscall"
	"time"
)

var ctrl *controller.Controller

func main() {
	// 解析命令行参数
	serverArgs := server_args.ParseArgs()

	// 初始化控制器
	ctrl = controller.NewController(serverArgs)

	// 设置gin引擎和路由
	r := setupRouter()

	// 启动服务器
	srv := &http.Server{
		Addr:    serverArgs.Host + ":" + fmt.Sprintf("%d", serverArgs.Port),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以关闭服务器（优雅）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")
}
