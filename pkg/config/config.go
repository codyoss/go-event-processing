package config

import (
	"io/ioutil"
	"os"

	"github.com/go-yaml/yaml"
)

// Config represents valid fields for a yaml configuration file
type Config struct {
	App struct {
		Parallelism        int `yaml:"parallelism"`
		ShutdownTimeoutSec int `yaml:"shutdownTimeoutSeconds"`
	} `yaml:"app"`
	Kafka struct {
		Brokers     []string `yaml:"brokers"`
		Group       string   `yaml:"group"`
		InputTopic  string   `yaml:"input-topic"`
		OutputTopic string   `yaml:"output-topic"`
	} `yaml:"kafka"`
	Nats struct {
		InputChannel  string `yaml:"input-channel"`
		OutputChannel string `yaml:"output-channel"`
		ClusterID     string `yaml:"cluster-id"`
		Group         string `yaml:"group"`
	} `yaml:"nats"`
}

// Parse parses a file name `config.yml`. This file should be in the same directory as the running binary.
func Parse(config *Config) {
	f, err := os.Open("config.yml")
	if err != nil {
		panic("Could not open config file")
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic("Could not read config file")
	}

	err = yaml.Unmarshal(b, config)
	if err != nil {
		panic("Could not parse config")
	}
}
