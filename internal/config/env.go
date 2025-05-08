package config

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type Env struct {
	AppEnv        string `mapstructure:"APP_ENV"`
	ServerAddress string `mapstructure:"SERVER_HOST"`
	Port          string `mapstructure:"PORT"`

	DBHost string `mapstructure:"DB_HOST"`
	DBPort string `mapstructure:"DB_PORT"`
	DBUser string `mapstructure:"DB_USER"`
	DBPass string `mapstructure:"DB_PASS"`
	DBName string `mapstructure:"DB_NAME"`

	RedisHost string `mapstructure:"REDIS_HOST"`
	RedisPort string `mapstructure:"REDIS_PORT"`
	RedisPass string `mapstructure:"REDIS_PASS"`

	AuthKey   string `mapstructure:"AUTH_KEY"`
	PublicKey string
}

func NewEnv() Env {
	env := Env{}

	viper.SetConfigFile(".env")

	_, err := os.Stat(".env")
	useEnvFile := !os.IsNotExist(err)

	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("REDIS_PASS", "")

	if useEnvFile {
		viper.SetConfigType("env")
		viper.SetConfigName(".env")
		viper.AddConfigPath(".")

		err := viper.ReadInConfig()
		if err != nil {
			log.Fatal("Can't read the .env file: ", err)
		}
	} else {
		viper.AutomaticEnv()
		val := reflect.ValueOf(&env).Elem()
		for i := 0; i < val.NumField(); i++ {
			field := val.Type().Field(i)
			tag := field.Tag.Get("mapstructure")
			if tag != "" {
				err = viper.BindEnv(tag)
				if err != nil {
					log.Fatal(err)
				}
				if value := viper.GetString(tag); value != "" {
					val.Field(i).SetString(value)
				}
			}
		}
	}

	viper.AutomaticEnv()
	err = viper.Unmarshal(&env)

	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	if env.AppEnv != "production" {
		log.Println("The App is running in development env")
	}

	privateKeyFormat := "-----BEGIN PUBLIC KEY-----\n%s\n-----END PUBLIC KEY-----"
	env.PublicKey = fmt.Sprintf(privateKeyFormat, env.AuthKey)

	return env
}

var Module = fx.Options(
	fx.Provide(NewEnv),
)
