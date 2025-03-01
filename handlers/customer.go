package handlers

import (
	"log"
	"os"

	"github.com/AbdelilahOu/Redpanda-go/types"
)

type CustomerHandler struct {
	logger *log.Logger
}

func NewCustomerHandler() CustomerHandler {
	return CustomerHandler{
		logger: log.New(os.Stdout, "customer-handler: ", log.LstdFlags),
	}
}

func (h *CustomerHandler) HandleCreate(customer *types.Customer) error {
	if customer == nil {
		h.logger.Println("Customer is nil in HandleCreate")
		return nil
	}
	h.logger.Printf("Customer created: %+v", customer)
	return nil
}

func (h *CustomerHandler) HandleUpdate(customer *types.Customer) error {
	if customer == nil {
		h.logger.Println("Customer is nil in HandleUpdate")
		return nil
	}
	h.logger.Printf("Customer updated: %+v", customer)
	return nil
}

func (h *CustomerHandler) HandleDelete(customer *types.Customer) error {
	if customer == nil {
		h.logger.Println("Customer is nil in HandleDelete")
		return nil
	}
	h.logger.Printf("Customer deleted: %+v", customer)
	return nil
}

func (h *CustomerHandler) HandleRead(customer *types.Customer) error {
	if customer == nil {
		h.logger.Println("Customer is nil in HandleRead")
		return nil
	}
	h.logger.Printf("Customer read: %+v", customer)
	return nil
}
