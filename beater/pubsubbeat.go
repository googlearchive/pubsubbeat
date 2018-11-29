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

package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"context"

	"runtime"

	"encoding/json"

	"cloud.google.com/go/pubsub"
	"github.com/GoogleCloudPlatform/pubsubbeat/config"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Pubsubbeat struct {
	done         chan struct{}
	config       *config.Config
	client       beat.Client
	pubsubClient *pubsub.Client
	subscription *pubsub.Subscription
	logger       *logp.Logger
}

func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config, err := config.GetAndValidateConfig(cfg)
	if err != nil {
		return nil, err
	}

	logger := logp.NewLogger(fmt.Sprintf("PubSub: %s/%s/%s", config.Project, config.Topic, config.Subscription.Name))
	logger.Infof("config retrieved: %+v", config)

	client, err := createPubsubClient(config)
	if err != nil {
		return nil, err
	}

	subscription, err := getOrCreateSubscription(client, config)
	if err != nil {
		return nil, err
	}

	bt := &Pubsubbeat{
		done:         make(chan struct{}),
		config:       config,
		pubsubClient: client,
		subscription: subscription,
		logger:       logger,
	}
	return bt, nil
}

func (bt *Pubsubbeat) Run(b *beat.Beat) error {
	bt.logger.Info("pubsubbeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-bt.done
		// The beat is stopping...
		bt.logger.Info("cancelling PubSub receive context...")
		cancel()
		bt.logger.Info("closing PubSub client...")
		bt.pubsubClient.Close()
	}()

	err = bt.subscription.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		// This callback is invoked concurrently by multiple goroutines
		var datetime time.Time
		eventMap := common.MapStr{
			"type":         b.Info.Name,
			"message_id":   m.ID,
			"publish_time": m.PublishTime,
			"message":      string(m.Data),
		}

		if len(m.Attributes) > 0 {
			eventMap["attributes"] = m.Attributes
		}

		if bt.config.Json.Enabled {
			var unmarshalErr error
			if bt.config.Json.FieldsUnderRoot {
				unmarshalErr = json.Unmarshal(m.Data, &eventMap)
				if unmarshalErr == nil && bt.config.Json.FieldsUseTimestamp {
					var timeErr error
					timestamp := eventMap[bt.config.Json.FieldsTimestampName]
					delete(eventMap, bt.config.Json.FieldsTimestampName)
					datetime, timeErr = time.Parse(bt.config.Json.FieldsTimestampFormat, timestamp.(string))
					if timeErr != nil {
						bt.logger.Errorf("Failed to format timestamp string as time. Using time.Now(): %s", timeErr)
					}
				}
			} else {
				var jsonData interface{}
				unmarshalErr = json.Unmarshal(m.Data, &jsonData)
				if unmarshalErr == nil {
					eventMap["json"] = jsonData
				}
			}

			if unmarshalErr != nil {
				bt.logger.Warnf("failed to decode json message: %s", unmarshalErr)
				if bt.config.Json.AddErrorKey {
					eventMap["error"] = common.MapStr{
						"key":     "json",
						"message": fmt.Sprintf("failed to decode json message: %s", unmarshalErr),
					}
				}
			}
		}

		if datetime.IsZero() {
			datetime = time.Now()
		}
		bt.client.Publish(beat.Event{
			Timestamp: datetime,
			Fields:    eventMap,
		})

		// TODO: Evaluate using AckHandler.
		m.Ack()
	})

	if err != nil {
		return fmt.Errorf("fail to receive message from subscription %q: %v", bt.subscription.String(), err)
	}

	return nil
}

func (bt *Pubsubbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}

func createPubsubClient(config *config.Config) (*pubsub.Client, error) {
	ctx := context.Background()
	userAgent := fmt.Sprintf(
		"Elastic/Pubsubbeat (%s; %s)", runtime.GOOS, runtime.GOARCH)
	options := []option.ClientOption{option.WithUserAgent(userAgent)}
	if config.CredentialsFile != "" {
		options = append(options, option.WithCredentialsFile(config.CredentialsFile))
	}

	client, err := pubsub.NewClient(ctx, config.Project, options...)
	if err != nil {
		return nil, fmt.Errorf("fail to create pubsub client: %v", err)
	}
	return client, nil
}

func getOrCreateSubscription(client *pubsub.Client, config *config.Config) (*pubsub.Subscription, error) {
	if !config.Subscription.Create {
		subscription := client.Subscription(config.Subscription.Name)
		return subscription, nil
	}

	topic := client.Topic(config.Topic)
	ctx := context.Background()

	subscription, err := client.CreateSubscription(ctx, config.Subscription.Name, pubsub.SubscriptionConfig{
		Topic:               topic,
		RetainAckedMessages: config.Subscription.RetainAckedMessages,
		RetentionDuration:   config.Subscription.RetentionDuration,
	})

	if st, ok := status.FromError(err); ok && st.Code() == codes.AlreadyExists {
		// The subscription already exists.
		subscription = client.Subscription(config.Subscription.Name)
	} else if ok && st.Code() == codes.NotFound {
		return nil, fmt.Errorf("topic %q does not exists", config.Topic)
	} else {
		return nil, fmt.Errorf("fail to create subscription: %v", err)
	}

	return subscription, nil
}
