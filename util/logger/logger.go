package logger

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	Path       string
	Level      string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool

	// ToStdout 为 true 时会同时输出到控制台（带颜色）。
	ToStdout bool
}

func defaultConfig() *Config {
	return &Config{
		Path:       "./logs/minichat.log",
		Level:      "debug",
		MaxSize:    100,
		MaxBackups: 7,
		MaxAge:     30,
		Compress:   false,
		ToStdout:   true,
	}
}

// InitLogger 初始化全局 logger，并返回实例。
// 如果 cfg 为 nil，则使用默认配置。
func InitLogger(cfg *Config) *zap.Logger {
	if cfg == nil {
		cfg = defaultConfig()
	}
	if cfg.Path == "" {
		cfg.Path = defaultConfig().Path
	}
	if cfg.MaxSize <= 0 {
		cfg.MaxSize = defaultConfig().MaxSize
	}
	if cfg.MaxBackups <= 0 {
		cfg.MaxBackups = defaultConfig().MaxBackups
	}
	if cfg.MaxAge <= 0 {
		cfg.MaxAge = defaultConfig().MaxAge
	}
	// 解析日志级别
	logLevel, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		logLevel = zapcore.DebugLevel
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// 短文件名: foo/bar/baz.go:123
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)

	// 控制台开启颜色
	consoleEncoderConfig := encoderConfig
	consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)

	// 自动创建目录（日志文件所在目录）
	logDir := filepath.Dir(cfg.Path)
	_ = os.MkdirAll(logDir, 0755)

	fileWriter := &lumberjack.Logger{
		Filename:   cfg.Path,
		MaxSize:    cfg.MaxSize, // megabytes
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,   // days
		Compress:   cfg.Compress, // disabled by default
	}

	cores := make([]zapcore.Core, 0, 2)
	cores = append(cores, zapcore.NewCore(jsonEncoder, zapcore.AddSync(fileWriter), logLevel))
	if cfg.ToStdout {
		cores = append(cores, zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), logLevel))
	}

	teeCore := zapcore.NewTee(cores...)
	l := zap.New(
		teeCore,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	zap.ReplaceGlobals(l)
	return l
}

// Sync 刷新缓冲区。一般在进程退出前调用（defer logger.Sync()）。
func Sync() {
	_ = zap.L().Sync()
}
