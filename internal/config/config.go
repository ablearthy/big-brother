package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type DbConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
	Sslmode  string `yaml:"sslmode"`
}

type ServerConfig struct {
	Port          int    `yaml:"port"`
	CookiesSecret string `yaml:"cookiesSecret"`
}

type Config struct {
	Db     DbConfig     `yaml:"db"`
	Server ServerConfig `yaml:"web"`
}

func NewConfig(configPath string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(configPath)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	d := yaml.NewDecoder(file)

	if err = d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
