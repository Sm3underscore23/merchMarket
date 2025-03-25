package config

import (
	"os"
	"path/filepath"

	"github.com/Sm3underscore23/merchStore/internal/customerrors"
	"github.com/Sm3underscore23/merchStore/internal/models"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func InitConfig(mainConfigPath, mainConfigName string, mainConfig *models.Config) error {
	configFile := filepath.Join(mainConfigPath, mainConfigName+".yaml")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return customerrors.ErrReadInConfig
	}
	viper.AddConfigPath(mainConfigPath)
	viper.SetConfigName(mainConfigName)
	err := viper.ReadInConfig()
	if err != nil {
		return customerrors.ErrReadInConfig
	}

	*mainConfig = models.Config{
		DB: models.DBConfig{
			Host:     viper.GetString("db.host"),
			Port:     viper.GetString("db.port"),
			Username: viper.GetString("db.username"),
			Password: getEnv(viper.GetString("db.password")),
			DBName:   viper.GetString("db.dbname"),
			SSLMode:  viper.GetString("db.sslmode"),
		},
		Auth: models.AuthConfig{
			Salt:      getEnv(viper.GetString("auth.salt")),
			SignedKey: getEnv(viper.GetString("auth.signedKey")),
		},
	}

	return nil
}

func getEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		logrus.Fatalf("error loading %s", key)
		return value
	}
	return value
}
