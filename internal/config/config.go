package config

// Config 这里只做最小化：先满足 logger/db 等基础模块需求。
// 后续你可以换成读取 yaml/viper 等更完整的配置体系。
//
// LogPath: 日志文件路径（默认 ./logs/minichat.log）
// LogLevel: debug/info/warn/error（默认 debug）
type Config struct {
	LogPath  string
	LogLevel string
}

func GetConfig() Config {
	return Config{LogPath: "./logs/minichat.log", LogLevel: "debug"}
}
