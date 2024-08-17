package http

import (
	"github.com/gofiber/fiber/v2"
	"main/internal/v1/ticket/domain/ports"
)

// MapRoutes Ticket routes
func MapRoutes(h ports.IHandlers, router fiber.Router) {

	tickets := router.Group("tickets")
	tickets.Get("/:id", h.GetTicket)
	tickets.Post("/", h.CreateTicket)
	tickets.Post("/:id/purchases", h.MakePurchase)
	tickets.Get("/:id/purchases", h.GetPurchaseLogs)
}
