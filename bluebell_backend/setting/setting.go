package setting

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Conf = new(AppConfig)

type AppConfig struct {
	Name         string `mapstructure:"name"`
	Mode         string `mapstructure:"mode"`
	Version      string `mapstructure:"version"`
	StartTime    string `mapstructure:"start_time"`
	MachineID    int64  `mapstructure:"machine_id"`
	Port         string `mapstructure:"port"`
	*MySQLConfig `mapstructure:"mysql"`
	*LogConfig   `mapstructure:"log"`
	*RedisConfig `mapstructure:"redis"`
}

type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DB           string `mapstructure:"dbname"`
	Port         int    `mapstructure:"port"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type LogConfig struct {
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Level      string `mapstructure:"level"`
}

type RedisConfig struct {
	Host         string `mapstructure:"host"`
	Password     string `mapstructure:"password"`
	Port         int    `mapstructure:"port"`
	DB           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
}

func Init() (err error) {
	viper.SetConfigFile("./config.yaml")
	err = viper.ReadInConfig()
	fmt.Println("read configFile success")
	if err != nil {
		fmt.Printf("read configFile failed, err: %v\n", err)
	}

	// 解析结构体
	err = viper.Unmarshal(Conf)
	if err != nil {
		fmt.Printf("unmarshal conf failed, err:%v\n", err)
	}

	// 监控配置文件变化
	viper.WatchConfig()
	// 当配置文件发生变化时，调用一个回调函数
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("config file changed...")
	})
	return
}
