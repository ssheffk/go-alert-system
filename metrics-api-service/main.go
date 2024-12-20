package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type MetricPayload struct {
	MetricName string    `json:"metricName" validate:"required"`
	Value      float64   `json:"value" validate:"required"`
	Timestamp  time.Time `json:"timestamp" validate:"required"`
}

var validate = validator.New()

func main() {

	// Kafka producer
	producer, err := sarama.NewSyncProducer([]string{"kafka:9092"}, nil)
	if err != nil {
		log.Fatalf("Error creating Kafka producer: %v", err)
	}
	defer producer.Close()

	// Gin router
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Metrics API service is running"})
	})

	r.POST("/api/metrics", func(c *gin.Context) {
		var payload MetricPayload

		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
			return
		}

		if err := validate.Struct(payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		metricsData, err := json.Marshal(payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal metrics data"})
			return
		}

		message := &sarama.ProducerMessage{
			Topic: "metrics-topic",
			Value: sarama.StringEncoder(metricsData),
		}

		_, _, err = producer.SendMessage(message)
		if err != nil {
			log.Printf("Failled to send message: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send metrics data"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Metrics data sent to Kafka"})

	})
	r.Run(":8001")
}
