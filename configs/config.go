package configs

import (
	"errors"
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	JWT    *JWTConfig
	Server *ServerConfig
	DB     *DBConfig
	Admin  *AdminConfig
}

type ServerConfig struct {
	Host string
	Port int
}

type JWTConfig struct {
	Secret            string
	Expiration        int
	RefreshSecret     string
	RefreshExpiration int
}

type DBConfig struct {
	Dialect  string `mapstructure:"dialect"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

type AdminConfig struct {
	Tokens []string
}

func loadConfig(filePath string) (*Config, error) {
	if filePath == "" {
		return nil, errors.New("config file path is required")
	}
	v := viper.New()
	v.SetConfigFile(filePath)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return &config, nil
}

type ConfigManager struct {
	config *Config
	mu     sync.RWMutex
}

func (manager *ConfigManager) LoadConfig(filepath string) error {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	config, err := loadConfig(filepath)
	if err != nil {
		return err
	}

	manager.config = config
	return nil
}

func (manager *ConfigManager) GetConfig(filepath string) (*Config, error) {

	if manager.config == nil {
		err := manager.LoadConfig(filepath)

		if err != nil {
			return nil, err
		}
	}

	return manager.config, nil
}
