package datasetexporter

import (
	"errors"

	"go.opentelemetry.io/collector/confmap"
)

type Config struct {
	ApiKey     string `mapstructure:"apikey"`
	DatasetUrl string `mapstructure:"dataseturl"`
}

func (c *Config) Unmarshal(conf *confmap.Conf) error {
	if err := conf.Unmarshal(c, confmap.WithErrorUnused()); err != nil {
		return nil
	}

	if len(c.DatasetUrl) > 0 && c.DatasetUrl[len(c.DatasetUrl)-1] == '/' {
		c.DatasetUrl = c.DatasetUrl[:len(c.DatasetUrl)-1]
	}

	return nil
}

func (c *Config) Validate() error {
	if c.ApiKey == "" {
		return errors.New("apikey is required")
	}
	if c.DatasetUrl == "" {
		return errors.New("dataseturl is required")
	}

	// Unfortunately cannot explicitly test the dataseturl / apikey here;
	// doing so would require writing a dummy log via uploadLogs or addEvents

	return nil
}
