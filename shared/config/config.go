package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

const (
	EnvLocal = "local"
	EnvProd  = "prod"
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

func IsEnv(env string) bool {
	return os.Getenv("GO_APP_ENV") == env
}

func GetAppPort() int {
	return viper.GetInt("app.port")
}

func GetDbConnectionString() string {
	return viper.GetString("database.connectionString")
}

func GetGinMode() string {
	return viper.GetString("app.ginMode")
}

func GetAppLogger() map[string]interface{} {
	return viper.GetStringMap("app.logger")
}
