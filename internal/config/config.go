package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// struct to hold the configuration values
type Config struct {
	Env         string `yaml:"env env:"ENV" env-required:"true"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	//struct embedding
	HTTPServer `yaml:"http_server"`
}

// different struct for http server cause it is nested
type HTTPServer struct {
	Addr string `yaml:"addr" env-required:"true"`
}

// here must be a function which executed otherwise the program should not run
func MustLoad() *Config {
	var configPath string

	//check if the config path is set in the environment variable
	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" { // is the the env path is an empty string
		//if the config path is not set in the env variable, check if it passed as command line argument
		flags := flag.String("config", "", "path to the config file")
		flag.Parse()

		configPath = *flags

		//if you still dont find the path then throw an error
		if configPath == "" {
			log.Fatal("config path is not set in the enviroment variable ")
		}
	}

	//check if any file is available at the config path or throw an error
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("config file not found")
	}

	var cfg Config

	//this is return an error
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("cannot read config file: %s", err.Error())
	}

	return &cfg
}
