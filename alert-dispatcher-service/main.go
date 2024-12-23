package main

import (
	"alert-dispatcher-service/utils"
	"log"
	"os"

	"github.com/IBM/sarama"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer([]string{"kafka:9092"}, config)
	if err != nil {
		log.Fatalf("Error creating Kafka consumer: %v", err)
	}

	go consumeTopic(consumer, "alerts-topic")
	select {}
}

func consumeTopic(consumer sarama.Consumer, topic string) {
	hostAddress := "smtp.gmail.com"
	hostPort := "465"
	mailSubject := "ALERT SYSTEM MESSAGE"
	emailAppPassword := os.Getenv("EMAIL_PASSWORD")
	yourMail := os.Getenv("EMAIL")
	recipient := os.Getenv("EMAIL")

	if emailAppPassword == "" || yourMail == "" || recipient == "" {
		log.Fatalf("Missing required environment variables")
	}
	log.Printf("Listening for messages on %s....", topic)

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Fatalf("Error creating Kafka partition consumer: %v", err)
	}
	defer partitionConsumer.Close()

	for message := range partitionConsumer.Messages() {
		utils.EmailClient(emailAppPassword, yourMail, recipient, hostAddress, hostPort, mailSubject, string(message.Value))
		log.Printf("Received message on %s: %s", topic, string(message.Value))
	}
}
