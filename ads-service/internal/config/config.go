package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	CfgType           string         `mapstructure:"cfg_type"`
	PgUrl             string         `mapstructure:"pg_url"`
	MigrationsDir     string         `mapstructure:"migrations_dir"`
	HTTPServerAddress string         `mapstructure:"http_server_address"`
	AuthGPRC          AuthGPRCConfig `mapstructure:"auth_grpc"`
}

type AuthGPRCConfig struct {
	Address   string `mapstructure:"address"`
	SecretKey string `mapstructure:"secret_key"`
	AppId     int32 `mapstructure:"app_id"`
}

var Cfg *Config

func MustLoad() {
	// Loading .env file which contains only CFG_TYPE
	err := godotenv.Load()
	if err != nil {
		panic("loading .env file:" + err.Error())
	}

	// Reading config from a file
	viper.SetConfigName(os.Getenv("CFG_TYPE"))
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	err = viper.ReadInConfig()
	if err != nil {
		panic("Error loading config: " + err.Error())
	}

	// Unmarshalling config to a struct
	config := &Config{}
	err = viper.UnmarshalExact(&config)
	if err != nil {
		panic("Error unmarshalling config to struct: " + err.Error())
	}
	slog.Info("cfg_type: " + config.CfgType)

	Cfg = config
}
