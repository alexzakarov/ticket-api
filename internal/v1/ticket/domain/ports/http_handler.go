package ports

import (
	"github.com/gofiber/fiber/v2"
)

// IHandlers HTTP handler interface
type IHandlers interface {
	//Tickets
	CreateTicket(c *fiber.Ctx) error
	GetTicket(c *fiber.Ctx) error

	//Purchase
	MakePurchase(c *fiber.Ctx) error
	GetPurchaseLogs(c *fiber.Ctx) error
}
