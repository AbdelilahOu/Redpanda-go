package handlers

import (
	"log"
	"os"

	"github.com/AbdelilahOu/bahmni-sync-service/types"
)

type OrderHandler struct {
	logger *log.Logger
}

func NewOrderHandler() OrderHandler {
	return OrderHandler{
		logger: log.New(os.Stdout, "logging-handler: ", log.LstdFlags),
	}
}

func (h *OrderHandler) HandleCreate(order *types.Order) error {
	if order == nil {
		h.logger.Println("Order is nil in HandleCreate")
		return nil
	}
	h.logger.Printf("Order created: %+v", order)
	return nil
}

func (h *OrderHandler) HandleUpdate(order *types.Order) error {
	if order == nil {
		h.logger.Println("Order is nil in HandleUpdate")
		return nil
	}
	h.logger.Printf("Order updated: %+v", order)
	return nil
}

func (h *OrderHandler) HandleDelete(order *types.Order) error {
	if order == nil {
		h.logger.Println("Order is nil in HandleDelete")
		return nil
	}
	h.logger.Printf("Order deleted: %+v", order)
	return nil
}

func (h *OrderHandler) HandleRead(order *types.Order) error {
	if order == nil {
		h.logger.Println("Order is nil in HandleRead")
		return nil
	}
	h.logger.Printf("Order read: %+v", order)
	return nil
}
