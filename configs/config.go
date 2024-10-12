package configs

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	APIs     APIsConfig     `mapstructure:"apis"`
	Redis    CacheConfig    `mapstructure:"redis"`
}

// Server
type ServerConfig struct {
	Port string `mapstructure:"port"`
}

// Database
type DatabaseConfig struct {
	URL                string `mapstructure:"url"`
	MaxIdleConnections int    `mapstructure:"max_idle_connections"`
	MaxOpenConnections int    `mapstructure:"max_open_connections"`
	ConnMaxLifetime    int    `mapstructure:"conn_max_lifetime"`
}

// Cache
type CacheType string

const (
	CacheTypeInMemory CacheType = "in_memory"
	CacheTypeRedis    CacheType = "redis"
)

type CacheConfig struct {
	Type     CacheType `mapstructure:"type"`
	Address  string    `mapstructure:"address"`
	Username string    `mapstructure:"username"`
	Password string    `mapstructure:"password"`
}

// External APIs
type APIsConfig struct {
	APIKeys map[string]string `mapstructure:"apikeys"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
