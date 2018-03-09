// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import (
	"fmt"

	"github.com/elastic/beats/libbeat/common"
)

type Config struct {
	Project      string `config:"project_id" validate:"required"`
	Topic        string `config:"topic" validate:"required"`
	Subscription string `config:"subscription" validate:"required"`
}

var DefaultConfig = Config{}

func GetAndValidateConfig(cfg *common.Config) (*Config, error) {
	c := DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("error in config file: %v", err)
	}
	return &c, nil
}
