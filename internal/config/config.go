package config

import (
	"strings"

	"github.com/gofiber/fiber/v3/log"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	ServerHost    string
	ServerPort    int
	RedisHost     string
	RedisPort     int
	RedisPassword string
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) InitConfig() error {
	err := godotenv.Load()
	if err != nil {
		logrus.Warn("config: failed load env %w ", err)
	}

	viper.SetEnvPrefix("WN")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	if err = viper.ReadInConfig(); err != nil {
		log.Fatal("config: failed read in config")
	}

	c.ServerHost = viper.GetString("config.server.host")
	c.ServerPort = viper.GetInt("config.server.port")
	c.RedisHost = viper.GetString("config.redis.host")
	c.RedisPassword = viper.GetString("config.redis.password")
	c.RedisPort = viper.GetInt("config.redis.port")

	return nil
}
