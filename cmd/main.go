package main

import (
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	merchstore "github.com/Sm3underscore23/merchStore"
	"github.com/Sm3underscore23/merchStore/internal/config"
	"github.com/Sm3underscore23/merchStore/internal/models"
	"github.com/Sm3underscore23/merchStore/pkg/handler"
	"github.com/Sm3underscore23/merchStore/pkg/repository"
	"github.com/Sm3underscore23/merchStore/pkg/service"
	"github.com/spf13/viper"
)

const (
	mainConfigPath = "configs"
	mainConfigName = "config"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := godotenv.Load(); err != nil {
		logrus.Fatal("error file .env not found")
	}

	var mainConfig models.Config

	if err := config.InitConfig(
		mainConfigPath,
		mainConfigName,
		&mainConfig,
	); err != nil {
		logrus.Fatalf("error initializing config: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(mainConfig.DB)

	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos, mainConfig.Auth)
	handlers := handler.NewHandler(services)

	srv := new(merchstore.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error occured while running http server: %s", err.Error())
	}
}
