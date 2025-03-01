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

func (h *CustomerHandler) HandleCreate(order *types.Customer) error {
	if order == nil {
		h.logger.Println("Customer is nil in HandleCreate")
		return nil
	}
	h.logger.Printf("Customer created: %+v", order)
	return nil
}

func (h *CustomerHandler) HandleUpdate(order *types.Customer) error {
	if order == nil {
		h.logger.Println("Customer is nil in HandleUpdate")
		return nil
	}
	h.logger.Printf("Customer updated: %+v", order)
	return nil
}

func (h *CustomerHandler) HandleDelete(order *types.Customer) error {
	if order == nil {
		h.logger.Println("Customer is nil in HandleDelete")
		return nil
	}
	h.logger.Printf("Customer deleted: %+v", order)
	return nil
}

func (h *CustomerHandler) HandleRead(order *types.Customer) error {
	if order == nil {
		h.logger.Println("Customer is nil in HandleRead")
		return nil
	}
	h.logger.Printf("Customer read: %+v", order)
	return nil
}
