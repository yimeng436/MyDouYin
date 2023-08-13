package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var globalConfig = new(GlobalConfig)

type GlobalConfig struct {
	*SvrConfig    `mapstructure:"svr_config"`
	*ConsulConfig `mapstructure:"consul"`
	*DbConfig     `mapstructure:"mysql"`
	*RedisConfig  `mapstructure:"redis"`
	*LogConfig    `mapstructure:"log"`
	RedsyncConfig []*RedsyncConfig `mapstructure:"redsync"`
}

type SvrConfig struct {
	Name string `mapstructure:"name"` // 服务name
	Host string `mapstructure:"host"` // 服务host
	Port int    `mapstructure:"port"`
}

type ConsulConfig struct {
	Host string   `mapstructure:"host" json:"host" yaml:"host"`
	Port int      `mapstructure:"port" json:"port" yaml:"port"`
	Tags []string `mapstructure:"tags" json:"tags" yaml:"tags"`
}

type DbConfig struct {
	Host        string `mapstructure:"host"`
	Port        string `mapstructure:"port"`
	Database    string `mapstructure:"database"`
	Username    string `mapstructure:"username"`
	Password    string `mapstructure:"password"`
	MaxIdleConn int    `mapstructure:"max_idle_conn"`  // 最大空闲连接数
	MaxOpenConn int    ` mapstructure:"max_open_conn"` // 最大打开的连接数
	MaxIdleTime int64  ` mapstructure:"max_idle_time"` // 连接最大空闲时间
}

type RedisConfig struct {
	DB       int `mapstructure:"db"`
	Port     int `mapstructure:"port"`
	PoolSize int `mapstructure:"pool_size"`
	// Expired      int    `mapstructure:"expired"`
	// MinIdleConns int    `mapstructure:"min_idle_conns"`
	Host     string `mapstructure:"host"`
	PassWord string `mapstructure:"password"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	FileName   string `mapstructure:"file_name"`
	LogPath    string `mapstructure:"log_path"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type RedsyncConfig struct {
	Port       int    `mapstructure:"port" json:"port" yaml:"port"`
	LockExpire int    `mapstructure:"expire" json:"expire" yaml:"expire"` // 锁过期时间
	PoolSize   int    `mapstructure:"pool_size" json:"pool_size" yaml:"pool_size"`
	Host       string `mapstructure:"host" json:"host" yaml:"host"`
	Password   string `mapstructure:"password" json:"password" yaml:"password"`
}

func Init() (err error) {
	// 自动推导项目根目录
	configFile := GetRootDir() + "/config/config.yaml"
	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("viper.read failed: %v\n", err)
		return fmt.Errorf("viper.Unmarshal failed: %v\n", err)
	}
	err = viper.Unmarshal(&globalConfig)
	if err != nil {
		fmt.Printf("viper.Unmarshal failed: %v\n", err)
		return fmt.Errorf("viper.Unmarshal failed: %v\n", err)
	}
	viper.WatchConfig() // 实时监控配置文件（热加载）
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改...")
		// 当配置文件信息发生变化 就修改 Conf 变量
		if err := viper.Unmarshal(globalConfig); err != nil {
			fmt.Printf("viper.Unmarshal failed: %v\n", err)
		}
	})
	return
}

func GetGlobalConfig() *GlobalConfig {
	return globalConfig
}
