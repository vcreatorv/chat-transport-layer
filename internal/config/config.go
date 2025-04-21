package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"time"
)

type HTTPConfig struct {
	Host           string        `yaml:"host"`
	Port           string        `yaml:"port"`
	ReadTimeout    time.Duration `yaml:"readTimeout"`
	WriteTimeout   time.Duration `yaml:"writeTimeout"`
	MaxHeaderBytes int           `yaml:"maxHeaderBytes"`
}

type KafkaProducerConfig struct {
	Brokers       []string `yaml:"brokers"`
	Topic         string   `yaml:"topic"`
	RequiredAcks  int      `yaml:"requiredAcks"`
	RetryMax      int      `yaml:"retryMax"`
	ReturnSuccess bool     `yaml:"returnSuccess"`
}

type KafkaConsumerConfig struct {
	Brokers         []string `yaml:"brokers"`
	Topic           string   `yaml:"topic"`
	AutoOffsetReset string   `yaml:"autoOffsetReset"`
	ReturnErrors    bool     `yaml:"returnErrors"`
}

type KafkaConfig struct {
	Producer KafkaProducerConfig `yaml:"producer"`
	Consumer KafkaConsumerConfig `yaml:"consumer"`
}

type SegmentConfig struct {
	MaxSegmentSize int           `yaml:"maxSegmentSize"`
	AssemblyPeriod time.Duration `yaml:"assemblyPeriod"`
}

type Config struct {
	HTTP    HTTPConfig    `yaml:"http"`
	Kafka   KafkaConfig   `yaml:"kafka"`
	Segment SegmentConfig `yaml:"segment"`
}

func LoadConfig(path string) (*Config, error) {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения конфигурационного файла %s: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(yamlFile, &cfg); err != nil {
		return nil, fmt.Errorf("ошибка парсинга yaml файла: %w", err)
	}

	return &cfg, nil
}
