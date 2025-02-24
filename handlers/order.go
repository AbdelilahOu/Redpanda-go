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
