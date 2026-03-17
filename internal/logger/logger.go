package logger

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func InitLogger() error {
	logrus.New()

	logLevel := viper.GetString("logger.logger_lvl")
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.Error("logging: failed parse lvl ", err)
	}
	logrus.SetLevel(level)

	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:            true,
		FullTimestamp:          true,
		TimestampFormat:        "2006-01-02 15:04:05",
		DisableLevelTruncation: true,
	})

	logrus.SetOutput(os.Stdout)

	return nil
}
