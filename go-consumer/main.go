package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/IBM/sarama"
)

type OrderEvent struct {
	Schema  Schema  `json:"schema"`
	Payload Payload `json:"payload"`
}

type Schema struct {
	Type   string `json:"type"`
	Fields []struct {
		Type     string `json:"type"`
		Field    string `json:"field"`
		Optional bool   `json:"optional"`
	} `json:"fields"`
	Optional bool   `json:"optional"`
	Name    string `json:"name"`
}

type Payload struct {
	Before      *Order      `json:"before"`
	After       *Order      `json:"after"`
	Source      Source      `json:"source"`
	Op          string      `json:"op"`
	TsMs        int64       `json:"ts_ms"`
	Transaction interface{} `json:"transaction"`
}

type Source struct {
	Version   string `json:"version"`
	Connector string `json:"connector"`
	Name      string `json:"name"`
	TsMs      int64  `json:"ts_ms"`
	Snapshot  string  `json:"snapshot"`
	Db        string `json:"db"`
	Table     string `json:"table"`
	ServerID  int    `json:"server_id"`
	GTID      string `json:"gtid"`
	File      string `json:"file"`
	Pos       int    `json:"pos"`
	Row       int    `json:"row"`
	Thread    int    `json:"thread"`
	Query     string `json:"query"`
}

type Order struct {
	OrderID     int         `json:"OrderID"`
	CustomerID  int         `json:"CustomerID"`
	OrderDate   interface{} `json:"OrderDate"`
	TotalAmount float64     `json:"TotalAmount"`
	Status      string      `json:"Status"`
}

type Consumer struct {
	ready    chan bool
	group    sarama.ConsumerGroup
	topics   []string
	logger   *log.Logger
	handlers map[string]func(order *Order) error
}

func NewConsumer(brokers []string, groupID string, topics []string) (*Consumer, error) {
	logger := log.New(os.Stdout, "kafka-consumer: ", log.LstdFlags)

	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Group.Session.Timeout = 20 * time.Second
	config.Consumer.Group.Heartbeat.Interval = 6 * time.Second

	group, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %v", err)
	}

	consumer := &Consumer{
		ready:    make(chan bool),
		group:    group,
		topics:   topics,
		logger:   logger,
		handlers: make(map[string]func(order *Order) error),
	}

	consumer.handlers["c"] = consumer.handleCreate
	consumer.handlers["u"] = consumer.handleUpdate
	consumer.handlers["d"] = consumer.handleDelete
	consumer.handlers["r"] = consumer.handleRead

	return consumer, nil
}

func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	c.logger.Println("Consumer Setup")
	close(c.ready)
	return nil
}

func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	c.logger.Println("Consumer Cleanup")
	return nil
}

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				c.logger.Printf("Message channel closed, exiting ConsumeClaim")
				return nil
			}

			c.logger.Printf("Message received: topic=%s partition=%d offset=%d\n",
				message.Topic, message.Partition, message.Offset)

			var event OrderEvent
			if err := json.Unmarshal(message.Value, &event); err != nil {
				c.logger.Printf("Error unmarshaling message: %v", err)
				session.MarkMessage(message, "")
				continue
			}

			// Handle the message based on operation type
			if handler, ok := c.handlers[event.Payload.Op]; ok {
				var order *Order
				switch event.Payload.Op {
				case "c", "r":
					order = event.Payload.After
				case "u":
					if event.Payload.After != nil {
						order = event.Payload.After
					} else {
						c.logger.Printf("After update is nil")
						session.MarkMessage(message, "")
						continue
					}
					c.logger.Printf("Before update: %+v", event.Payload.Before)
				case "d":
					order = event.Payload.Before
				}

				if err := handler(order); err != nil {
					c.logger.Printf("Error handling operation %s: %v", event.Payload.Op, err)
				}
			} else {
				c.logger.Printf("Unknown operation type: %s", event.Payload.Op)
			}

			session.MarkMessage(message, "")

		case <-session.Context().Done():
			c.logger.Println("ConsumeClaim context is done, exiting")
			return nil
		}
	}
}

func (c *Consumer) handleCreate(order *Order) error {
	if order == nil {
		c.logger.Println("Order is nil in handleCreate")
		return nil
	}
	c.logger.Printf("Order created: %+v", order)
	return nil
}

func (c *Consumer) handleUpdate(order *Order) error {
	if order == nil {
		c.logger.Println("Order is nil in handleUpdate")
		return nil
	}
	c.logger.Printf("Order updated: %+v", order)
	return nil
}

func (c *Consumer) handleDelete(order *Order) error {
	if order == nil {
		c.logger.Println("Order is nil in handleDelete")
		return nil
	}
	c.logger.Printf("Order deleted: %+v", order)
	return nil
}

func (c *Consumer) handleRead(order *Order) error {
	if order == nil {
		c.logger.Println("Order is nil in handleRead")
		return nil
	}
	c.logger.Printf("Order read: %+v", order)
	return nil
}

func (c *Consumer) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			c.logger.Println("Starting a new consumer session")
			if err := c.group.Consume(ctx, c.topics, c); err != nil {
				c.logger.Printf("Error from consumer: %v", err)
				if ctx.Err() != nil {
					c.logger.Println("Context is done, exiting consumer loop")
					return
				}
			} else {
				c.logger.Println("Consumer session ended gracefully")
			}
			if ctx.Err() != nil {
				c.logger.Println("Context is done, exiting consumer loop")
				return
			}
			c.ready = make(chan bool)
		}
	}()

	select {
	case <-c.ready:
		c.logger.Println("Consumer is ready")
	case <-time.After(10 * time.Second):
		c.logger.Println("Timeout waiting for consumer to become ready")
	}

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case <-sigterm:
			c.logger.Println("Terminating: via signal")
			cancel()
		case <-ctx.Done():
			// Context is already cancelled
		}
	}()

	<-ctx.Done() // Wait for either signal or internal termination
	c.logger.Println("Terminating...")

	c.logger.Println("Waiting for consumer goroutine to finish")
	wg.Wait()
	c.logger.Println("Closing consumer group")
	err := c.group.Close()
	c.logger.Println("Consumer group closed")
	return err
}

func main() {
	brokers := []string{"localhost:9092"}
	topics := []string{"orders.mydb.Orders"}
	groupID := "orders"

	consumer, err := NewConsumer(brokers, groupID, topics)
	if err != nil {
		log.Fatalf("Error creating consumer: %v", err)
	}

	if err := consumer.Start(); err != nil {
		log.Fatalf("Error running consumer: %v", err)
	}
	log.Println("Consumer stopped gracefully")
}
