package config

import (
	"strings"

	"github.com/gofiber/fiber/v3/log"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func InitConfig() error {
	err := godotenv.Load()
	if err != nil {
		logrus.Error("config: failed load env %w ", err)
		return err
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

	return nil
}
