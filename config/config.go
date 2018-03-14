// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import (
	"fmt"

	"time"

	"os"

	"github.com/elastic/beats/libbeat/common"
)

type Config struct {
	Project         string `config:"project_id" validate:"required"`
	Topic           string `config:"topic" validate:"required"`
	CredentialsFile string `config:"credentials_file"`
	Subscription    struct {
		Name                string        `config:"name" validate:"required"`
		RetainAckedMessages bool          `config:"retain_acked_messages"`
		RetentionDuration   time.Duration `config:"retention_duration"`
	}
	Json struct {
		Enabled     bool `config:"enabled"`
		AddErrorKey bool `config:"add_error_key"`
	}
}

var DefaultConfig = Config{}

func GetAndValidateConfig(cfg *common.Config) (*Config, error) {
	c := DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("error in config file: %v", err)
	}

	if d, _ := time.ParseDuration("10m"); c.Subscription.RetentionDuration < d {
		return nil, fmt.Errorf("retention_duration cannot be shorter than 10 minutes")
	}

	if d, _ := time.ParseDuration("168h"); c.Subscription.RetentionDuration > d {
		return nil, fmt.Errorf("retention_duration cannot be longer than 7 days")
	}

	if c.CredentialsFile != "" {
		if _, err := os.Stat(c.CredentialsFile); os.IsNotExist(err) {
			return nil, fmt.Errorf("cannot find the credentials_file %q", c.CredentialsFile)
		}
	}

	return &c, nil
}
