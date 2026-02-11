package di

import "minichat/internal/config"

func provideAppConfig() config.AppConfig {
	return config.GetConfig()
}
