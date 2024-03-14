package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Database struct {
		DriverName string `yaml:"driver_name"`
		Host       string `yaml:"host"`
		Port       string `yaml:"port"`
		User       string `yaml:"user"`
		Password   string `yaml:"pass"`
		DbName     string `yaml:"dbname"`
		SslMode    string `yaml:"sslmode"`
	} `yaml:"database"`
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"server"`
}

func LoadConfig() *Config {
	f, err := os.Open("config.yaml")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		panic(err)
	}
	log.Println("Got configs!")
	return &cfg
}
