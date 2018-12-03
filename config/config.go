// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
		Create              bool          `config:"create"`
	}
	Json struct {
		Enabled               bool   `config:"enabled"`
		AddErrorKey           bool   `config:"add_error_key"`
		FieldsUnderRoot       bool   `config:"fields_under_root"`
		FieldsUseTimestamp    bool   `config:"fields_use_timestamp"`
		FieldsTimestampName   string `config:"fields_timestamp_name"`
		FieldsTimestampFormat string `config:"fields_timestamp_format"`
	}
}

func GetDefaultConfig() Config {
	config := Config{}
	config.Subscription.Create = true
	config.Json.FieldsTimestampName = "@timestamp"
	return config
}

var DefaultConfig = GetDefaultConfig()

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
