package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	AwsRegion       string `env:"AWS_REGION,required"`
	KnowledgeBaseID string `env:"KNOWLEDGE_BASE_ID,required"`
	DataSourceID    string `env:"DATA_SOURCE_ID,required"`
	Port            int    `env:"PORT" envDefault:"8080"`
}

func NewConfig() (*Config, error) {
	cfg, err := env.ParseAs[Config]()
	return &cfg, err
}

func NewConfigMust() *Config {
	cfg, err := NewConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}
	return cfg
}

func (c *Config) GetAddress() string {
	return fmt.Sprintf(":%d", c.Port)
}
