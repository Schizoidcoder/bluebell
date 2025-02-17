package main

import (
	"bluebell/Kafka"
	"bluebell/controller"
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/logger"
	"bluebell/pkg/snowflake"
	"bluebell/routes"
	"bluebell/settings"
	"bluebell/websocket"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

// goweb开发通用脚手架模版
func main() {
	//1。 加载配置
	if err := settings.Init(); err != nil {
		fmt.Println("init setting failed", err)
		return
	}
	//2。初始化日志
	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		fmt.Println("init logger failed", err)
		return
	}
	defer zap.L().Sync()
	zap.L().Debug("init logger success")
	//3。初始化Mysql
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		fmt.Println("init mysql failed", err)
		return
	}
	defer mysql.Close()
	//4。初始化redis
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Println("init redis failed", err)
		return
	}
	defer redis.Close()

	if err := snowflake.Init(settings.Conf.StartTime, settings.Conf.MachineID); err != nil {
		fmt.Println("init snowflake failed", err)
		return
	}
	//注册gin框架内置的校验器使用的翻译器
	if err := controller.InitTrans("zh"); err != nil {
		fmt.Println("init validator trans failed", err)
		return
	}
	Kafka.Init()
	zap.L().Debug("init kafka success")
	go websocket.Init()
	zap.L().Debug("init websocket success")
	//5。注册路由
	r := routes.SetupRouter(settings.Conf.Mode)
	//6。启动服务（优雅关机）
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", settings.Conf.Port),
		Handler: r,
	}

	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	//KafkaTesting
	//err := Kafka.Kafka_writer.WriteMessages(context.Background(), kafka.Message{
	//	Topic: "like_event",
	//	Key:   []byte("1"),
	//	Value: []byte("2"),
	//}, kafka.Message{
	//	Topic: "comment_event",
	//	Key:   []byte("3"),
	//	Value: []byte("4"),
	//})
	//if err != nil {
	//	fmt.Println("kafka write messages failed", err)
	//}

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	Kafka.Close()
	zap.L().Info("Shutdown Server ...")
	// 创建一个2秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	// 2秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过2秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}

	zap.L().Info("Server exiting")

}
