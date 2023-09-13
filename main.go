package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"web_app/controller"
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/logger"
	"web_app/routes"
	"web_app/settings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	//1. 加载配置文件
	if err := settings.Init(); err != nil {
		fmt.Printf("Init settings failed, err:%v\n", err)
		return
	}

	//2. 初始化日志
	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		fmt.Printf("Init loggers failed, err:%v\n", err)
	}
	defer zap.L().Sync()
	zap.L().Debug("logger inits successfully...")

	//3. 初始化mysql数据库
	if err := mysql.InitDB(); err != nil {
		fmt.Printf("Init mysql failed, err:%v\n", err)
	}
	defer mysql.Close()

	//4. 初始化redis数据库
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Printf("Init redis failed, err:%v\n", err)
	}
	defer redis.Close()

	//5. 用gin框架注册路由, 初始化一个翻译器
	if err := controller.InitTrans("zh"); err != nil {
		fmt.Printf("init transalator failed, err: %v\n", err)
		return
	}
	r := routes.Setup(settings.Conf.Mode)

	//6. 优雅关机
	srv := &http.Server{
		Addr:                         fmt.Sprintf(":%d", viper.GetInt("app.port")),
		Handler:                      r,
		DisableGeneralOptionsHandler: false,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zap.L().Info("Shutdown Server...")
	//设置超时时间为5秒
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Info("Server Shutdown", zap.Error(err))
	}
	zap.L().Info("Server exiting")
}
