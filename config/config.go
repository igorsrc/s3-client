package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	App struct {
		Port          string `yaml:"port"`
		MaxUploadSize int64  `yaml:"max-upload-size"`
	} `yaml:"app"`
	S3 struct {
		Host   string `yaml:"host"`
		Port   string `yaml:"port"`
		Access string `yaml:"access"`
		Secret string `yaml:"secret"`
		Bucket string `yaml:"bucket"`
		Region string `yaml:"region"`
	} `yaml:"s3"`
}

var cfg Config

func Get() *Config {
	return &cfg
}

func InitConfig(file string) {

	if file == "" {
		panic("Config file not found")
	}

	config, err := ioutil.ReadFile(file)
	if err != nil {
		panic("Error reading " + file)
	}

	err = yaml.Unmarshal(config, &cfg)
}
