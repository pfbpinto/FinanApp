package messaging

import (
	"log"

	"github.com/nsqio/go-nsq"
)

// CreateNSQProducer establishes the producer connection to NSQ
func CreateNSQProducer() (*nsq.Producer, error) {
	producer, err := nsq.NewProducer("localhost:4150", nsq.NewConfig())
	if err != nil {
		log.Println("Error creating NSQ Producer:", err)
		return nil, err
	}
	return producer, nil
}

// CreateNSQConsumer establishes the consumer connection to NSQ
func CreateNSQConsumer(topic string, channel string) (*nsq.Consumer, error) {
	consumer, err := nsq.NewConsumer(topic, channel, nsq.NewConfig())
	if err != nil {
		log.Println("Error creating NSQ Consumer:", err)
		return nil, err
	}
	return consumer, nil
}
