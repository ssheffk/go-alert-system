package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
)

type MetricPayload struct {
	MetricName string    `json:"metricName"`
	Value      float64   `json:"value"`
	Timestamp  time.Time `json:"timestamp"`
}

type Rule struct {
	MetricName string  `json:"metricName"`
	Threshold  float64 `json:"threshold"`
	Operator   string  `json:"operator"` // // e.g., ">", "<", ">=", "<=", "=="
}

type Alert struct {
	MetricName string    `json:"metricName"`
	Value      float64   `json:"value"`
	Timestamp  time.Time `json:"timestamp"`
	Message    string    `json:"message"`
}

func main() {

	// Kafka consumer
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer([]string{"kafka:9092"}, config)
	if err != nil {
		log.Fatalf("Error creating Kafka consumer: %v", err)
	}

	// Kafka producer for alerts
	producer, err := sarama.NewSyncProducer([]string{"kafka:9092"}, nil)
	if err != nil {
		log.Fatalf("Error creating Kafka producer: %v", err)
	}
	defer producer.Close()

	// Start consuming metrics
	go consumeMetrics(consumer, producer)

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Rule engine service is running"})
	})

	r.Run(":8003")
}

func consumeMetrics(consumer sarama.Consumer, producer sarama.SyncProducer) {
	partitionConsumer, err := consumer.ConsumePartition("metrics-topic", 0, sarama.OffsetOldest)
	if err != nil {
		log.Fatalf("Error creating Kafka partition consumer: %v", err)
	}
	defer partitionConsumer.Close()

	for message := range partitionConsumer.Messages() {
		var metric MetricPayload
		if err := json.Unmarshal(message.Value, &metric); err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			continue
		}

		// Fetch rules from the rule-management-service
		rules, err := fetchRules()
		if err != nil {
			log.Printf("Error fetching rules: %v", err)
			continue
		}

		// evaluate metric against rules
		for _, rule := range rules {
			if rule.MetricName == metric.MetricName && evaluateRule(metric.Value, rule) {
				alert := Alert{
					MetricName: metric.MetricName,
					Value:      metric.Value,
					Timestamp:  metric.Timestamp,
					Message:    generateAlertMessage(metric, rule),
				}
				sendAlert(producer, alert)
			}
		}
	}
}

func evaluateRule(metricValue float64, rule Rule) bool {
	switch rule.Operator {
	case ">":
		return metricValue > rule.Threshold
	case "<":
		return metricValue < rule.Threshold
	case ">=":
		return metricValue >= rule.Threshold
	case "<=":
		return metricValue <= rule.Threshold
	case "==":
		return metricValue == rule.Threshold
	default:
		log.Printf("Unsupported operator: %s", rule.Operator)
		return false
	}
}

func generateAlertMessage(metric MetricPayload, rule Rule) string {
	return "ALERT: " + metric.MetricName + " value " + fmt.Sprintf("%.2f", metric.Value) + " breached threshold " + fmt.Sprintf("%.2f", rule.Threshold)
}

func sendAlert(producer sarama.SyncProducer, alert Alert) {
	alertData, err := json.Marshal(alert)
	if err != nil {
		log.Printf("Error marshalling alert: %v", err)
		return
	}

	message := &sarama.ProducerMessage{
		Topic: "alerts-topic",
		Value: sarama.StringEncoder(alertData),
	}

	_, _, err = producer.SendMessage(message)
	if err != nil {
		log.Printf("Error sending alert: %v", err)
	}
}

func fetchRules() ([]Rule, error) {

	// ruleServiceURL := "http://go-alert-system-rule-management-service.default.svc.cluster.local:8004/api/rules"

	// Define the URL of the rule-management-service
	ruleServiceURL := "http://rule-management-service:8004/api/rules"

	// Make the HTTP GET request
	resp, err := http.Get(ruleServiceURL)
	if err != nil {
		return nil, fmt.Errorf("error fetching rules: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK status code: %d", resp.StatusCode)
	}

	// Parse the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var rules []Rule
	if err := json.Unmarshal(body, &rules); err != nil {
		return nil, fmt.Errorf("error unmarshalling rules: %w", err)
	}

	return rules, nil
}
