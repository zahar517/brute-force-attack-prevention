package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Logger   LoggerConf   `json:"logger"`
	Database DatabaseConf `json:"database"`
	Server   ServerConf   `json:"server"`
}

type LoggerConf struct {
	Level string `json:"level"`
	File  string `json:"file"`
}

type DatabaseConf struct {
	Name           string `json:"name"`
	Host           string `json:"host"`
	Port           string `json:"port"`
	User           string `json:"user"`
	Password       string `json:"password"`
	MigrationsPath string `json:"migrationsPath"`
	Dsn            string
}

type ServerConf struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	GrpcHost string `json:"grpcHost"`
	GrpcPort string `json:"grpcPort"`
}

func NewConfig(path string) (*Config, error) {
	var config Config

	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	config.Database.Dsn = fmt.Sprintf(
		"host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		config.Database.Password,
		config.Database.Name,
	)

	return &config, nil
}
