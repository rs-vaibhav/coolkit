package config

import (
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	Port        string `mapstructure:"PORT"`
	GinMode     string `mapstructure:"GIN_MODE"`
	DatabaseURL string `mapstructure:"DATABASE_URL"`
	DBHost      string `mapstructure:"DB_HOST"`
	DBPort      string `mapstructure:"DB_PORT"`
	DBUser      string `mapstructure:"DB_USER"`
	DBPassword  string `mapstructure:"DB_PASSWORD"`
	DBName      string `mapstructure:"DB_NAME"`
	DBSSLMode   string `mapstructure:"DB_SSLMODE"`
	JWTSecret   string `mapstructure:"JWT_SECRET"`
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

		viper.SetDefault("PORT", "8080")
		viper.SetDefault("GIN_MODE", "debug")
		viper.SetDefault("DB_HOST", "localhost")
		viper.SetDefault("DB_PORT", "5432")
		viper.SetDefault("DB_USER", "coolkit")
		viper.SetDefault("DB_PASSWORD", "coolkit")
		viper.SetDefault("DB_NAME", "coolkit")
		viper.SetDefault("DB_SSLMODE", "disable")
		viper.SetDefault("JWT_SECRET", "change-me-in-production")

		_ = viper.BindEnv("PORT")
		_ = viper.BindEnv("GIN_MODE")
		_ = viper.BindEnv("DATABASE_URL")
		_ = viper.BindEnv("DB_HOST")
		_ = viper.BindEnv("DB_PORT")
		_ = viper.BindEnv("DB_USER")
		_ = viper.BindEnv("DB_PASSWORD")
		_ = viper.BindEnv("DB_NAME")
		_ = viper.BindEnv("DB_SSLMODE")
		_ = viper.BindEnv("JWT_SECRET")

		_ = viper.ReadInConfig() // Ignore file not found error

		cfg = &Config{}
		_ = viper.Unmarshal(cfg)
	})
}

func Get() *Config {
	return cfg
}
