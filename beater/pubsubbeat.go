package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"context"

	"cloud.google.com/go/pubsub"
	"github.com/rosbo/pubsubbeat/config"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Pubsubbeat struct {
	done         chan struct{}
	config       *config.Config
	client       beat.Client
	pubsubClient *pubsub.Client
	subscription *pubsub.Subscription
}

func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config, err := config.GetAndValidateConfig(cfg)
	if err != nil {
		return nil, err
	}

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
	}
	return bt, nil
}

func (bt *Pubsubbeat) Run(b *beat.Beat) error {
	logp.Info("pubsubbeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-bt.done
		// The beat is stopping...
		// Cancel the pubsub message receiver first.
		cancel()
		// Then close the pubsub client
		bt.pubsubClient.Close()
	}()

	err = bt.subscription.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		// This callback is invoked concurrently by multiple goroutines
		messageMap := common.MapStr{
			"id":           m.ID,
			"publish_time": m.PublishTime,
			"data":         string(m.Data),
		}

		if len(m.Attributes) > 0 {
			messageMap["attributes"] = m.Attributes
		}

		bt.client.Publish(beat.Event{
			Timestamp: time.Now(),
			Fields: common.MapStr{
				"type":    b.Info.Name,
				"message": messageMap,
			},
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
	// TODO: Add user-agent to track our impact
	ctx := context.Background()

	// TODO: Support different auth mecanism
	client, err := pubsub.NewClient(ctx, config.Project)
	if err != nil {
		return nil, fmt.Errorf("fail to create pubsub client: %v", err)
	}
	return client, nil
}

func getOrCreateSubscription(client *pubsub.Client, config *config.Config) (*pubsub.Subscription, error) {
	topic := client.Topic(config.Topic)
	ctx := context.Background()

	subscription, err := client.CreateSubscription(ctx, config.Subscription, pubsub.SubscriptionConfig{
		Topic: topic,
		// TODO: Allow setting other parameters such as AckDeadline
	})

	if st, ok := status.FromError(err); ok && st.Code() == codes.AlreadyExists {
		// The subscription already exists.
		subscription = client.Subscription(config.Subscription)
	} else if ok && st.Code() == codes.NotFound {
		return nil, fmt.Errorf("topic %q does not exists", config.Topic)
	} else {
		return nil, fmt.Errorf("fail to create subscription: %v", err)
	}

	return subscription, nil
}
