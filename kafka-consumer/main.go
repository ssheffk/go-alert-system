package main

import (
	"log"

	"github.com/IBM/sarama"
)

func main() {

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer([]string{"kafka:9092"}, config)
	if err != nil {
		log.Fatalf("Error creating Kafka consumer: %v", err)
	}

	defer consumer.Close()

	go consumeTopic(consumer, "metrics-topic")

	go consumeTopic(consumer, "alerts-topic")

	// Prevent the main function from exiting immediately
	select {}
}

func consumeTopic(consumer sarama.Consumer, topic string) {
	log.Printf("Listening for messages on %s....", topic)

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Fatalf("Error creating Kafka partition consumer: %v", err)
	}
	defer partitionConsumer.Close()

	for message := range partitionConsumer.Messages() {
		log.Printf("Received message on %s: %s", topic, string(message.Value))
	}
}
