package consumer

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

type Consumer[T any] struct {
	ready   chan bool
	group   sarama.ConsumerGroup
	topics  []string
	logger  *log.Logger
	handler HandlersInterface[T]
	retries int
}

func NewConsumer[T any](brokers []string, groupID string, topics []string, handler HandlersInterface[T]) (*Consumer[T], error) {
	logger := log.New(os.Stdout, fmt.Sprintf("consumer-%s: ", groupID), log.LstdFlags)

	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second
	config.Consumer.Offsets.Retry.Max = 4
	config.Consumer.Group.Session.Timeout = 20 * time.Second
	config.Consumer.Group.Heartbeat.Interval = 6 * time.Second

	config.Consumer.Retry.Backoff = 2 * time.Second

	config.Net.MaxOpenRequests = 5
	config.Net.DialTimeout = 30 * time.Second
	config.Net.ReadTimeout = 30 * time.Second
	config.Net.WriteTimeout = 30 * time.Second

	config.Consumer.Return.Errors = true

	group, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %v", err)
	}

	consumer := &Consumer[T]{
		ready:   make(chan bool),
		group:   group,
		topics:  topics,
		logger:  logger,
		handler: handler,
		retries: 0,
	}

	return consumer, nil
}
func (c *Consumer[T]) Setup(sarama.ConsumerGroupSession) error {
	c.logger.Println("Consumer Setup")
	close(c.ready)
	return nil
}
func (c *Consumer[T]) Cleanup(sarama.ConsumerGroupSession) error {
	c.logger.Println("Consumer Cleanup")
	return nil
}
func (c *Consumer[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				c.logger.Printf("Message channel closed, exiting ConsumeClaim")
				return nil
			}

			c.logger.Printf("Message received: topic=%s partition=%d offset=%d",
				message.Topic, message.Partition, message.Offset)

			var event Event[T]
			if err := json.Unmarshal(message.Value, &event); err != nil {
				c.logger.Printf("Error unmarshaling message: %v", err)
				continue
			}

			var err error

			switch event.Payload.Op {
			case "c":
				if event.Payload.After != nil {
					err = c.handler.HandleCreate(event.Payload.After)
				}
			case "u":
				if event.Payload.After != nil {
					c.logger.Printf("Before update: %+v", event.Payload.Before)
					err = c.handler.HandleUpdate(event.Payload.After)
				} else {
					c.logger.Printf("After update is nil")
				}
			case "d":
				if event.Payload.Before != nil {
					err = c.handler.HandleDelete(event.Payload.Before)
				}
			case "r":
				if event.Payload.After != nil {
					err = c.handler.HandleRead(event.Payload.After)
				}
			default:
				c.logger.Printf("Unknown operation type: %s", event.Payload.Op)
			}

			if err != nil {
				c.retries = c.retries + 1
				c.logger.Printf("Error handling operation %s: %v for the %d TH time", event.Payload.Op, err, c.retries)
				if c.retries > 3 {
					c.retries = 0
					continue
				}
				time.Sleep(time.Duration(c.retries) * time.Second)
				return err
			}

			session.MarkMessage(message, "")

		case <-session.Context().Done():
			c.logger.Println("ConsumeClaim context is done, exiting")
			return nil
		}
	}
}
func (c *Consumer[T]) Start() error {
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
		}
	}()

	<-ctx.Done()
	c.logger.Println("Terminating...")

	c.logger.Println("Waiting for consumer goroutine to finish")
	wg.Wait()
	c.logger.Println("Closing consumer group")
	err := c.group.Close()
	c.logger.Println("Consumer group closed")
	return err
}
