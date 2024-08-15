package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env            string        `yaml:"env" env-required:"true"`
	GRPC           GRPCConfig    `yaml:"grpc" env-required:"true"`
	PGUrl          string        `yaml:"pg_url" env-required:"true"`
	MigrationsPath string        `yaml:"migrations_path" env-required:"true"`
	TokenTTL       time.Duration `yaml:"token_ttl" env-required:"true"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

var Cfg *Config

func MustLoad() {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}
	// проверка существования файла при наличии пути
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	// распаковываем прочитанный конфиг в структуру
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("config path is empty: " + err.Error())
	}

	Cfg = &cfg
}

// fetchConfigPath запрашивает путь к конфигу через командную строку или
// переменную окружения CONFIG_PATH.
// Приоритет: flag > env.
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		err := godotenv.Load()
		if err != nil {
			panic("loading .env file: " + err.Error())
		}
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
