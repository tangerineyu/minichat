package config

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Addr    string `yaml:"addr"`
	GinMode string `yaml:"gin_mode"`
}

type LoggerConfig struct {
	Path       string `yaml:"path"`
	Level      string `yaml:"level"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
	Compress   bool   `yaml:"compress"`
	ToStdout   bool   `yaml:"to_stdout"`
}

type DBConfig struct {
	Driver string `yaml:"driver"`
	DSN    string `yaml:"dsn"`
}

type JWTConfig struct {
	Secret           string `yaml:"secret"`
	AccessTTLMinutes int    `yaml:"access_ttl_minutes"`
	RefreshTTLHours  int    `yaml:"refresh_ttl_hours"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type OSSConfig struct {
	Endpoint        string `yaml:"endpoint"`
	AccessKeyID     string `yaml:"access_key_id"`
	AccessKeySecret string `yaml:"access_key_secret"`
	BucketName      string `yaml:"bucket_name"`
}
type AppConfig struct {
	Server ServerConfig `yaml:"server"`
	Logger LoggerConfig `yaml:"logger"`
	DB     DBConfig     `yaml:"db"`
	JWT    JWTConfig    `yaml:"jwt"`
	Redis  RedisConfig  `yaml:"redis"`
	OSS    OSSConfig    `yaml:"oss"`
}

func defaultConfig() AppConfig {
	return AppConfig{
		Server: ServerConfig{Addr: ":8080", GinMode: "debug"},
		Logger: LoggerConfig{
			Path:       "./logs/minichat.log",
			Level:      "debug",
			MaxSize:    100,
			MaxBackups: 7,
			MaxAge:     30,
			Compress:   false,
			ToStdout:   true,
		},
		DB:    DBConfig{Driver: "sqlite", DSN: "minichat.db"},
		JWT:   JWTConfig{Secret: "minichat-dev-secret-change-me", AccessTTLMinutes: 120, RefreshTTLHours: 168},
		Redis: RedisConfig{Addr: "6379", Password: "", DB: 3},
	}
}

// Load reads yaml config from path. If file doesn't exist, it returns defaults.
func Load(path string) (AppConfig, error) {
	//加载.env文件，没有就忽略
	_ = godotenv.Load()
	cfg := defaultConfig()
	if path == "" {
		path = "internal/config/config.yaml"
	}

	p := path
	// 如果不是绝对路径，则转换为相对于当前工作目录的路径
	if !filepath.IsAbs(p) {
		// 获取当前工作目录
		if wd, err := os.Getwd(); err == nil {
			p = filepath.Join(wd, path)
		}
	}

	b, err := os.ReadFile(p)
	if err != nil {
		// 如果配置文件不存在，则返回默认配置
		if errors.Is(err, os.ErrNotExist) {
			// 读取不到配置文件，使用默认配置并应用环境变量覆盖
			applyEnvOverrides(&cfg)
			return cfg, nil
		}
		return cfg, err
	}
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return cfg, err
	}
	applyEnvOverrides(&cfg)
	return cfg, nil
}

// GetConfig keeps old call sites working: load default location.
func GetConfig() AppConfig {
	cfg, _ := Load("")
	return cfg
}

func applyEnvOverrides(cfg *AppConfig) {
	// server
	if v := os.Getenv("MINICHAT_LISTEN_ADDR"); v != "" {
		cfg.Server.Addr = v
	}
	if v := os.Getenv("GIN_MODE"); v != "" {
		cfg.Server.GinMode = v
	}

	// logger
	if v := os.Getenv("MINICHAT_LOG_PATH"); v != "" {
		cfg.Logger.Path = v
	}
	if v := os.Getenv("MINICHAT_LOG_LEVEL"); v != "" {
		cfg.Logger.Level = v
	}

	// jwt (kept consistent with util/jwt/jwt.go)
	if v := os.Getenv("MINICHAT_JWT_SECRET"); v != "" {
		cfg.JWT.Secret = v
	}
	if v := os.Getenv("MINICHAT_JWT_ACCESS_TTL_MINUTES"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			cfg.JWT.AccessTTLMinutes = n
		}
	}
	if v := os.Getenv("MINICHAT_JWT_REFRESH_TTL_HOURS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			cfg.JWT.RefreshTTLHours = n
		}
	}
}
