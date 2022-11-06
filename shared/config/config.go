package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

func Init() {
	env := os.Getenv("GO_APP_ENV")
	if env == "" {
		env = "local"
	}

	viper.SetConfigName(fmt.Sprintf("config.%s", env))
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

func GetAppPort() int {
	return viper.GetInt("app.port")
}

func GetDbConnectionString() string {
	return viper.GetString("database.connectionString")
}
