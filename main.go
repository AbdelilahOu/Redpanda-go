package main

import (
	"log"
	"sync"

	"github.com/AbdelilahOu/bahmni-sync-service/consumer"
	"github.com/AbdelilahOu/bahmni-sync-service/handlers"
	"github.com/AbdelilahOu/bahmni-sync-service/types"
)

func main() {
	brokers := []string{"localhost:9092"}
	topic := "orders_topic_Orders"

	orderHandler := handlers.NewOrderHandler()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		consumer, err := consumer.NewConsumer[types.Order](
			brokers,
			"orders-group",
			[]string{topic},
			&orderHandler,
		)
		if err != nil {
			log.Fatalf("Error creating logging consumer: %v", err)
		}

		if err := consumer.Start(); err != nil {
			log.Printf("Error running logging consumer: %v", err)
		}
	}()

	wg.Wait()
	log.Println("All consumers stopped gracefully")
}
