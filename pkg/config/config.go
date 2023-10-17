package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	AppName string   `yaml:"appName"`
	Api     *Api     `yaml:"api"`
	Logging *Logging `yaml:"logging"`
}
type Api struct {
	Port int  `yaml:"port"`
	Cors Cors `yaml:"cors"`
}

type Cors struct {
	AllowedOrigins   []string `yaml:"allowedOrigins"`
	AllowedMethods   []string `yaml:"allowedMethods"`
	AllowedHeaders   []string `yaml:"allowedHeaders"`
	ExposedHeaders   []string `yaml:"exposedHeaders"`
	AllowCredentials bool     `yaml:"allowCredentials"`
	MaxAge           int      `yaml:"maxAge"`
}

type Logging struct {
	Level string `yaml:"level"`
}

func Read(file string) *Config {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}

	d := yaml.NewDecoder(f)

	var cfg Config

	err = d.Decode(&cfg)
	if err != nil {
		_ = f.Close()
		log.Fatal(err)
	}

	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}

	return &cfg
}
