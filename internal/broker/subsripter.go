package broker

import (
	"context"
	"log"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
)

type EventSubsripter struct {
	subscriber *amqp.Subscriber
}

func NewEventSubsripter(amqpURI, queueSuffix string) (*EventSubsripter, error) {
	amqpConfig := amqp.NewDurablePubSubConfig(amqpURI, amqp.GenerateQueueNameTopicNameWithSuffix(queueSuffix))
	subscriber, err := amqp.NewSubscriber(amqpConfig, watermill.NewStdLogger(true, true))
	if err != nil {
		log.Fatalf("Connection to amqp failed: %v", err)
	}
	return &EventSubsripter{subscriber}, nil
}

func (eventSubsripter *EventSubsripter) SubscribeMessange(ctx context.Context, topic string, data chan []byte, end chan bool) error {
	message, err := eventSubsripter.subscriber.Subscribe(ctx, topic)
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case msg := <-message:
				data <- msg.Payload
				msg.Ack()
				break
			case <-end:
				break
			}
		}
	}()
	return nil
}
