package main

import (
	"log"
	"sync"

	"github.com/AbdelilahOu/Redpanda-go/consumer"
	"github.com/AbdelilahOu/Redpanda-go/handlers"
)

func main() {
	brokers := []string{"localhost:19092"}

	orderHandler := handlers.NewOrderHandler()
	customerHandler := handlers.NewCustomerHandler()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		consumer, err := consumer.NewConsumer(
			brokers,
			"orders-group",
			[]string{"orders.mydb.Orders"},
			&orderHandler,
		)
		if err != nil {
			log.Fatalf("Error creating orders consumer: %v", err)
		}

		if err := consumer.Start(); err != nil {
			log.Printf("Error running orders consumer: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		consumer, err := consumer.NewConsumer(
			brokers,
			"customers-group",
			[]string{"customers.mydb.Customers"},
			&customerHandler,
		)
		if err != nil {
			log.Fatalf("Error creating customers consumer: %v", err)
		}

		if err := consumer.Start(); err != nil {
			log.Printf("Error running customers consumer: %v", err)
		}
	}()

	wg.Wait()
	log.Println("All consumers stopped gracefully")
}
