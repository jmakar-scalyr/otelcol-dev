package datasetexporter

import (
	"errors"

	"go.opentelemetry.io/collector/confmap"
)

type Config struct {
	ApiKey     string `mapstructure:"apikey"`
	DataSetUrl string `mapstructure:"dataseturl"`
}

func (c *Config) Unmarshal(conf *confmap.Conf) error {
	if err := conf.Unmarshal(c, confmap.WithErrorUnused()); err != nil {
		return nil
	}
	return nil
}

func (c *Config) Validate() error {
	if c.ApiKey == "" {
		return errors.New("apikey is required")
	}
	if c.DataSetUrl == "" {
		return errors.New("dataseturl is required")
	}
	return nil
}
