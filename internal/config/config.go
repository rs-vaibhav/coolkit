package config

import (
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort string `mapstructure:"SERVER_PORT"`
	GinMode    string `mapstructure:"GIN_MODE"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBSSLMode  string `mapstructure:"DB_SSLMODE"`
	JWTSecret  string `mapstructure:"JWT_SECRET"`
}

var (
	cfg  *Config
	once sync.Once
)

func Load() {
	once.Do(func() {
		viper.SetConfigName(".env")
		viper.SetConfigType("env")
		viper.AddConfigPath(".")

		viper.SetDefault("SERVER_PORT", "8080")
		viper.SetDefault("GIN_MODE", "debug")
		viper.SetDefault("DB_HOST", "localhost")
		viper.SetDefault("DB_PORT", "5432")
		viper.SetDefault("DB_USER", "coolkit")
		viper.SetDefault("DB_PASSWORD", "coolkit")
		viper.SetDefault("DB_NAME", "coolkit")
		viper.SetDefault("DB_SSLMODE", "disable")
		viper.SetDefault("JWT_SECRET", "change-me-in-production")

		viper.BindEnv("SERVER_PORT")
		viper.BindEnv("GIN_MODE")
		viper.BindEnv("DB_HOST")
		viper.BindEnv("DB_PORT")
		viper.BindEnv("DB_USER")
		viper.BindEnv("DB_PASSWORD")
		viper.BindEnv("DB_NAME")
		viper.BindEnv("DB_SSLMODE")
		viper.BindEnv("JWT_SECRET")

		_ = viper.ReadInConfig() // Ignore file not found error

		cfg = &Config{}
		_ = viper.Unmarshal(cfg)
	})
}

func Get() *Config {
	return cfg
}
