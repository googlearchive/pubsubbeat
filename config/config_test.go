// +build !integration

package config

import (
	"testing"

	"github.com/elastic/beats/libbeat/common"
)

func TestGetAndValidateConfigMissingRequiredFields(t *testing.T) {
	cases := map[string]struct {
		Project, Topic, Subscription string
	}{
		"missing project": {
			Topic:        "a-topic",
			Subscription: "a-subscription",
		},
		"missing topic": {
			Topic:        "a-topic",
			Subscription: "a-subscription",
		},
		"missing subscription": {
			Project: "a-project",
			Topic:   "a-topic",
		},
	}

	for tn, tc := range cases {
		c := common.NewConfig()
		c.SetString("project", -1, tc.Project)
		c.SetString("topic", -1, tc.Topic)
		c.SetString("subscription", -1, tc.Subscription)
		_, err := GetAndValidateConfig(c)

		if err == nil {
			t.Errorf("%s: expected to fail", tn)
		}
	}
}
