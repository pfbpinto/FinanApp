package messaging

import (
	"fmt"
	"log"

	"github.com/nsqio/go-nsq"
)

type Producer struct {
	producer *nsq.Producer
}

// NewProducer creates a new instance of Producer using the connection function from connection.go
func NewProducer() (*Producer, error) {

	// Using the function you already have in connection.go to create the producer
	producer, err := CreateNSQProducer()
	if err != nil {
		return nil, fmt.Errorf("error creating NSQ producer: %v", err)
	}
	return &Producer{producer: producer}, nil
}

// UserRegistration publishes a message to the 'user_registration' topic in NSQ
func (p *Producer) UserRegistration(message []byte) error {
	topic := "user_registration" // Name of the topic where the messages will be published
	err := p.producer.Publish(topic, message)
	if err != nil {
		return fmt.Errorf("error publishing message to NSQ: %v", err)
	}
	log.Printf("Message published to topic '%s'. Message: '%s' ", topic, message)
	return nil
}
