package config

import "github.com/caarlos0/env/v11"

type Config struct {
	AwsRegion       string `env:"AWS_REGION,required"`
	KnowledgeBaseID string `env:"KNOWLEDGE_BASE_ID,required"`
	DataSourceID    string `env:"DATA_SOURCE_ID,required"`
}

func NewConfig() (*Config, error) {
	return env.ParseAs[*Config]()
}
