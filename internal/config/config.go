package config

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string     `yaml:"env" env-default:"LOCAL"`
	HTTPServer HTTPServer `yaml:"http_server"`
	Postgres   Postgres   `yaml:"postgres"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type Postgres struct {
	Host     string `yaml:"host"`
	Database string `yaml:"database"`
	User     string `yaml:"user"`
	Password string `env:"POSTGRES_PASSWORD"`
}

type Flags struct {
	configPath string
}

func MustLoad() *Config {
	flags := parseFlags()
	if flags.configPath == "" {
		log.Fatal("no config path")
	}

	if _, err := os.Stat(flags.configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", flags.configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(flags.configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}

func parseFlags() *Flags {
	var configPath string
	flag.StringVar(&configPath, "config", "", "config path")
	flag.Parse()
	return &Flags{
		configPath: configPath,
	}
}
