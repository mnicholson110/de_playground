package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type orderStatus int

const (
	created orderStatus = iota
	confirmed
	paid
	shipped
	delivered
)

type Order struct {
	OrderId     int         `json:"orderId"`
	OrderStatus orderStatus `json:"orderStatus"`
}

func main() {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9094"})
	if err != nil {
		panic(err)
	}
	defer p.Close()

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	topic := "orderTopic"
	var maxOrderID int
	orderMap := make(map[int]orderStatus)

	for maxOrderID < 1000 {
		// Randomly choose an existing order and increment the status, or create a new order
		a := rand.Intn(3)
		if a != 0 {
			// Create a new order
			maxOrderID++
			orderMap[maxOrderID] = created
			order := Order{OrderId: maxOrderID, OrderStatus: orderMap[maxOrderID]}
			orderJSON, err := json.Marshal(order)
			if err != nil {
				panic(err)
			}
			// Produce the new order to the topic with orderId as the key
			p.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Key:            []byte(strconv.Itoa(maxOrderID)),
				Value:          orderJSON,
			}, nil)
		} else {
			// Update an existing order that is not yet delivered
			if maxOrderID == 0 {
				continue
			}
			orderID := rand.Intn(maxOrderID)
			if orderMap[orderID] != delivered {
				orderMap[orderID]++
				order := Order{OrderId: orderID, OrderStatus: orderMap[orderID]}
				orderJSON, err := json.Marshal(order)
				if err != nil {
					panic(err)
				}
				// Produce the updated order to the topic with orderId as the key
				p.Produce(&kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
					Key:            []byte(strconv.Itoa(orderID)),
					Value:          orderJSON,
				}, nil)
			}
		}

		p.Flush(15 * 1000)
	}
}
