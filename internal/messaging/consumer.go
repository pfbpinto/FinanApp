package messaging

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nsqio/go-nsq"
)

type Consumer struct {
	consumer *nsq.Consumer
}

// NewConsumer creates a new instance of Consumer using the connection function from connection.go
func NewConsumer() (*Consumer, error) {

	// Tópico e canal para o consumidor
	topic := "user_registration"
	channel := "user_channel"

	// Criar o consumidor
	consumer, err := nsq.NewConsumer(topic, channel, nsq.NewConfig())
	if err != nil {
		return nil, fmt.Errorf("error creating NSQ consumer: %v", err)
	}

	// Retornar o consumidor
	return &Consumer{consumer: consumer}, nil
}

// HandleMessage processes the messages received from NSQ
func (c *Consumer) HandleMessage(message *nsq.Message) error {

	// Decode the message data
	var userData map[string]interface{}
	if err := json.Unmarshal(message.Body, &userData); err != nil {
		log.Printf("Error decoding the message: %v", err)
		return err
	}

	// Logic to process the user data
	username := userData["username"].(string)
	email := userData["email"].(string)
	password := userData["password"].(string)
	role := userData["role"].(string)

	// Placeholder for user creation logic
	err := createUser(username, email, password, role)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return err
	}

	log.Printf("User successfully registered: %s", username)
	return nil
}

// StartConsumer starts the consumer and begins processing messages
func (c *Consumer) StartConsumer() {

	// Define how the messages will be handled
	c.consumer.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		return c.HandleMessage(message)
	}))

	// Start the consumer and connect to the NSQD instance
	if err := c.consumer.ConnectToNSQD("localhost:4150"); err != nil {
		log.Printf("Error connecting to NSQD: %v", err)
	}
}

// Placeholder function for creating a user in the database
func createUser(username, email, password, role string) error {	
	log.Printf("Creating user: %s, Email: %s, Role: %s , Pass: %s", username, email, role,password)	

	// TODO save to database	
	return nil
}
