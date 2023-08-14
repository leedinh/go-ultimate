package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type OrderPlacer struct {
	producer  *kafka.Producer
	topic     string
	deliverch chan kafka.Event
}

func NewOrderPlacer(p *kafka.Producer, topic string) *OrderPlacer {
	return &OrderPlacer{
		producer:  p,
		topic:     topic,
		deliverch: make(chan kafka.Event, 10000),
	}
}

func (op *OrderPlacer) placeOrder(orderType string, size int) error {
	var (
		format  = fmt.Sprintf("%s - %d", orderType, size)
		payload = []byte(format)
	)

	err := op.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &op.topic, Partition: kafka.PartitionAny},
		Value:          payload},
		op.deliverch,
	)

	if err != nil {
		log.Fatal(err)
	}
	<-op.deliverch
	return nil
}

func main() {
	topic := "GGG"
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"client.id":         "supersecret",
		"acks":              "all"})

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}

	go func() {
		consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
			"bootstrap.servers": "localhost:9092",
			"group.id":          "supersecret",
			"auto.offset.reset": "smallest"})
		if err != nil {
			log.Fatal(err)
		}

		err = consumer.Subscribe(topic, nil)
		if err != nil {
			log.Fatal(err)
		}

		for {
			ev := consumer.Poll(100)
			switch e := ev.(type) {
			case *kafka.Message:
				fmt.Printf("consumed message from the queue: %s\n", string(e.Value))
			case kafka.Error:
				fmt.Printf("%% Error: %v\n", e)
			}
		}

	}()

	op := NewOrderPlacer(p, topic)
	for i := 0; i < 1000; i++ {
		if err := op.placeOrder("marker", i); err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second * 3)
	}
}
