// +build !integration

package config

import (
	"testing"

	"time"

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
		c := common.NewConfig()
		c.SetString("project_id", -1, tc.Project)
		c.SetString("topic", -1, tc.Topic)

		sConfig := common.NewConfig()
		sConfig.SetString("name", -1, tc.Subscription)
		sConfig.SetBool("retain_acked_messages", -1, false)
		sConfig.SetString("retention_duration", -1, "10m")
		c.SetChild("subscription", -1, sConfig)

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
		c := common.NewConfig()
		c.SetString("project_id", -1, "a-project")
		c.SetString("topic", -1, "a-topic")

		sConfig := common.NewConfig()
		sConfig.SetString("name", -1, "a-subscription")
		sConfig.SetBool("retain_acked_messages", -1, tc.RetainAckedMessages)
		sConfig.SetString("retention_duration", -1, tc.RetentionDuration)
		c.SetChild("subscription", -1, sConfig)

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
