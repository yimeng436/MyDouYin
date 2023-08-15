package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var globalConfig = new(GlobalConfig)

type GlobalConfig struct {
	*SvrConfig    `mapstructure:"svr_config"`
	*LogConfig    `mapstructure:"log" json:"log" yaml:"log"`
	*ConsulConfig `mapstructure:"consul" json:"consul" yaml:"consul"`
}

type SvrConfig struct {
	Name        string `mapstructure:"name"` // 服务name
	Host        string `mapstructure:"host"` // 服务host
	Port        int    `mapstructure:"port"`
	Mode        string `mapstructure:"mode"`          // gin模式
	UserSvrName string `mapstructure:"user_svr_name"` // 用户服务name
	VideoPath   string `mapstructure:"video_path"`    // 视频存放路径（有耦合，后面处理）
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	FileName   string `mapstructure:"file_name"`
	LogPath    string `mapstructure:"log_path"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}
type ConsulConfig struct {
	Host string   `mapstructure:"host" json:"host" yaml:"host"`
	Port int      `mapstructure:"port" json:"port" yaml:"port"`
	Tags []string `mapstructure:"tags" json:"tags" yaml:"tags"`
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
