package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

// DatabaseConfig 对应 database 部分的配置
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type LogConfig struct {
	Level    string `mapstructure:"level"`
	Format   string `mapstructure:"format"`
	Output   string `mapstructure:"output"`
	FilePath string `mapstructure:"file_path"`
}

type SubfinderConfig struct {
	Threads            int    `mapstructure:"threads"`
	Timeout            int    `mapstructure:"timeout"`
	MaxEnumerationTime int    `mapstructure:"max_enumeration_time"`
	AllSources         bool   `mapstructure:"all_sources"`
	ProviderConfigFile string `mapstructure:"provider_config_file"`
	TaskTimeoutSeconds int    `mapstructure:"task_timeout_seconds"` // 新增
	TaskMaxRetry       int    `mapstructure:"task_max_retry"`       // 新增
}

type NaabuConfig struct {
	Ports      string `mapstructure:"ports"`
	Rate       int    `mapstructure:"rate"`
	Timeout    int    `mapstructure:"timeout"`
	Retries    int    `mapstructure:"retries"`
	ScanType   string `mapstructure:"scan_type"`
	EnableNmap bool   `mapstructure:"enable_nmap"`
	NmapCLI    string `mapstructure:"nmap_cli"`
	ExcludeCdn bool   `mapstructure:"exclude_cdn"`
}

type DnsxConfig struct {
	Threads   int      `mapstructure:"threads"`
	Retries   int      `mapstructure:"retries"`
	Resolvers []string `mapstructure:"resolvers"`
}

type ToolsConfig struct {
	Subfinder SubfinderConfig `mapstructure:"subfinder"`
	Naabu     NaabuConfig     `mapstructure:"naabu"`
	Dnsx      DnsxConfig      `mapstructure:"dnsx"`
}

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Log      LogConfig      `mapstructure:"log"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Tools    ToolsConfig    `mapstructure:"tools"`
}

// 全局配置变量
var Cfg *Config

// LoadConfig 从 configs/config.yaml 加载配置
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")    // 配置文件名 (不带后缀)
	viper.SetConfigType("yaml")      // 配置文件类型
	viper.AddConfigPath("./configs") // 配置文件路径

	// 读取配置
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var cfg Config
	// 将读取的配置反序列化到结构体中
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	Cfg = &cfg
	return &cfg, nil
}
