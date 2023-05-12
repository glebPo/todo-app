package main

import (
	"fmt"
	"os"

	"github.com/glebpo/todo-app"
	"github.com/glebpo/todo-app/pkg/handler"
	"github.com/glebpo/todo-app/pkg/repository"
	"github.com/glebpo/todo-app/pkg/service"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	_ "github.com/lib/pq"
)

func main() {
	if err := initConfig(); err != nil {
		logrus.Fatalf("ERROR init config:  %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("ERROR loading env  variables %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handler := handler.NewHandler(services)

	logrus.Info(fmt.Sprintf("server listen on port %s", viper.GetString("port")))
	srv := new(todo.Server)
	if err := srv.Run(viper.GetString("port"), handler.InitRoutes()); err != nil {
		logrus.Fatalf("ERROR occured while running http server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
