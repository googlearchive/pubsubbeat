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

// +build !integration

package config

import (
	"testing"

	"time"

	"path/filepath"

	"github.com/elastic/beats/libbeat/common"
)

func TestGetAndValidateConfigMissingRequiredFields(t *testing.T) {
	cases := map[string]struct {
		Project, Topic, Subscription string
		ExpectError                  bool
	}{
		"missing project": {
			Topic:        "a-topic",
			Subscription: "a-subscription",
			ExpectError:  true,
		},
		"missing topic": {
			Topic:        "a-topic",
			Subscription: "a-subscription",
			ExpectError:  true,
		},
		"missing subscription": {
			Project:     "a-project",
			Topic:       "a-topic",
			ExpectError: true,
		},
		"no missing field": {
			Project:      "a-project",
			Topic:        "a-topic",
			Subscription: "a-subscription",
			ExpectError:  false,
		},
	}

	for tn, tc := range cases {
		c := createDefaultTestConfig()

		c.SetString("project_id", -1, tc.Project)
		c.SetString("topic", -1, tc.Topic)

		sConfig, _ := c.Child("subscription", -1)
		sConfig.SetString("name", -1, tc.Subscription)

		_, err := GetAndValidateConfig(c)

		if tc.ExpectError {
			if err == nil {
				t.Errorf("%s: expected to fail", tn)
			}
		} else {
			if err != nil {
				t.Errorf("%s: not expected to fail; %s", tn, err)
			}
		}
	}
}

func TestGetAndValidateConfigSubscriptionConfig(t *testing.T) {
	cases := map[string]struct {
		RetainAckedMessages bool
		RetentionDuration   string
		ExpectError         bool
	}{
		"168h retention and keep acked messages": {
			RetainAckedMessages: true,
			RetentionDuration:   "168h",
			ExpectError:         false,
		},
		"10m retention and don't keep acked messages": {
			RetainAckedMessages: false,
			RetentionDuration:   "10m",
			ExpectError:         false,
		},
		"retention period invalid format": {
			RetainAckedMessages: true,
			RetentionDuration:   "1d", // Duration should be in hours or minutes.
			ExpectError:         true,
		},
		"retention period too short": {
			RetainAckedMessages: true,
			RetentionDuration:   "9m",
			ExpectError:         true,
		},
		"retention period too long": {
			RetainAckedMessages: true,
			RetentionDuration:   "168h1m",
			ExpectError:         true,
		},
	}

	for tn, tc := range cases {
		c := createDefaultTestConfig()

		sConfig, _ := c.Child("subscription", -1)
		sConfig.SetBool("retain_acked_messages", -1, tc.RetainAckedMessages)
		sConfig.SetString("retention_duration", -1, tc.RetentionDuration)

		conf, err := GetAndValidateConfig(c)

		if tc.ExpectError {
			if err == nil {
				t.Errorf("%s: expected to fail", tn)
			}
		} else {
			if err != nil {
				t.Errorf("%s: not expected to fail; %s", tn, err)
			}

			if conf.Subscription.RetainAckedMessages != tc.RetainAckedMessages {
				t.Errorf("%s: expected retain_acked_messages to be %t", tn, tc.RetainAckedMessages)
			}

			if d, _ := time.ParseDuration(tc.RetentionDuration); conf.Subscription.RetentionDuration != d {
				t.Errorf("%s: expected retention_duration to be %s", tn, tc.RetentionDuration)
			}
		}
	}
}

func TestGetAndValidateConfigCredentialsFile(t *testing.T) {
	cases := map[string]struct {
		CredentialsFilePath string
		ExpectError         bool
	}{
		"credentials file exists": {
			CredentialsFilePath: filepath.Join("testdata", "fake-creds.json"),
			ExpectError:         false,
		},
		"credentials file does not exist": {
			CredentialsFilePath: filepath.Join("testdata", "missing.json"),
			ExpectError:         true,
		},
	}

	for tn, tc := range cases {
		c := createDefaultTestConfig()
		c.SetString("credentials_file", -1, tc.CredentialsFilePath)

		_, err := GetAndValidateConfig(c)

		if tc.ExpectError {
			if err == nil {
				t.Errorf("%s: expected to fail", tn)
			}
		} else {
			if err != nil {
				t.Errorf("%s: not expected to fail; %s", tn, err)
			}
		}
	}
}

// Creates a config with valid default values
func createDefaultTestConfig() *common.Config {
	c := common.NewConfig()
	c.SetString("project_id", -1, "a-project")
	c.SetString("topic", -1, "a-topic")

	sConfig := common.NewConfig()
	sConfig.SetString("name", -1, "a-subscription")
	sConfig.SetBool("retain_acked_messages", -1, false)
	sConfig.SetString("retention_duration", -1, "10m")
	c.SetChild("subscription", -1, sConfig)

	return c
}
