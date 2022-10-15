package main

// @title bluebell
// @version 1.0
// description web项目; 一个博客系统的后端
// termsOfService http://swagger.io/terms/

// @contact.name ropz
// @contact.url http://www.swagger.io/support
// @contact.email 1098985413@qq.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 127.0.0.1:8080
// @BasePath C:\Users\TZH\GolandProjects\bluebell_demo\bluebell_backend

import (
	"bluebell/controller"
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/logger"
	"bluebell/pkg/snowflake"
	"bluebell/router"
	"bluebell/setting"
	"fmt"
	"go.uber.org/zap"
)

func main() {

	// 1.加载配置
	if err := setting.Init(); err != nil {
		fmt.Printf("init setting failed, err:%v\n", err)
		return
	}

	// 2.初始化日志
	if err := logger.Init(setting.Conf.LogConfig, setting.Conf.Mode); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	defer zap.L().Sync() // 把缓冲区的日志追加到我们的日志文件中
	zap.L().Debug("logger init success...")

	// 3.初始化MySQL连接
	if err := mysql.Init(setting.Conf.MySQLConfig); err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}
	defer mysql.Close()

	// 4.初始化Redis连接
	if err := redis.Init(setting.Conf.RedisConfig); err != nil {
		fmt.Printf("init redis failed, err:%v\n", err)
		return
	}
	defer redis.Close()

	if err := snowflake.Init(setting.Conf.StartTime, setting.Conf.MachineID); err != nil {
		//fmt.Println(setting.Conf.StartTime)
		fmt.Printf("init snowflake failed, err:%v\n", err)
		return
	}

	// 初始化gin框架内置的校验器使用的翻译器
	if err := controller.InitTrans("zh"); err != nil {
		fmt.Printf("init validator trans failed, err:%v\n", err)
	}
	// 注册路由
	r := router.SetupRouter(setting.Conf.Mode)
	r.Run()

}
