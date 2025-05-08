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

	Bucket      string `mapstructure:"BUCKET_NAME"`
	S3Key       string `mapstructure:"AWS_ACCESS_KEY"`
	S3Secret    string `mapstructure:"AWS_SECRET_KEY"`
	S3Host      string `mapstructure:"AWS_HOST"`
	S3PublicUrl string `mapstructure:"AWS_PUBLIC_URL"`
	S3Region    string `mapstructure:"AWS_REGION"`

	CasdoorUrl              string `mapstructure:"CASDOOR_URL"`
	CasdoorClientId         string `mapstructure:"CASDOOR_CLIENT_ID"`
	CasdoorClientSecret     string `mapstructure:"CASDOOR_CLIENT_SECRET"`
	CasdoorCertRaw          string `mapstructure:"CASDOOR_CERT"`
	CasdoorCert             string
	CasdoorOrganizationName string `mapstructure:"CASDOOR_ORGANIZATION_NAME"`
	CasdoorApplicationName  string `mapstructure:"CASDOOR_APPLICATION_NAME"`
	DefaultServerId         string `mapstructure:"DEFAULT_SERVER_ID"`

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
	certFormat := "-----BEGIN CERTIFICATE-----\n%s\n-----END CERTIFICATE-----"
	if env.CasdoorCertRaw != "" {
		env.CasdoorCert = fmt.Sprintf(certFormat, env.CasdoorCertRaw)
	}
	env.PublicKey = fmt.Sprintf(privateKeyFormat, env.AuthKey)

	return env
}

var Module = fx.Options(
	fx.Provide(NewEnv),
)
