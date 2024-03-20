package infrastructure

import (
	"context"
	"log"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
)

type EventSubsripter struct {
	subscriber *amqp.Subscriber
}

func NewEventSubsripter(amqpURI string) (*EventSubsripter, error) {
	amqpConfig := amqp.NewDurableQueueConfig(amqpURI)
	subscriber, err := amqp.NewSubscriber(amqpConfig, watermill.NewStdLogger(true, true))
	if err != nil {
		log.Fatalf("Connection to amqp failed: %v", err)
	}
	return &EventSubsripter{subscriber}, nil
}

func (eventSubsripter *EventSubsripter) SubscribeMessange(ctx context.Context, topic string, data chan []byte) error {
	message, err := eventSubsripter.subscriber.Subscribe(ctx, topic)
	if err != nil {
		return err
	}

	for msg := range message {
		data <- msg.Payload
		msg.Ack()
	}
	return nil
}
