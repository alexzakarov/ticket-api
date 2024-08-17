package ports

import (
	"github.com/goccy/go-json"
	"main/internal/v1/ticket/domain/entities"
)

// IService service interface
type IService interface {
	//Tickets
	CreateTicket(entities.CreateTicketReqDto) (int64, entities.TicketResDto)
	GetTicket(int64) (int64, entities.TicketResDto)

	//Purchase
	MakePurchase(int64, entities.MakePurchaseReqDto) int64
	GetPurchaseLogs(int64) (int64, json.RawMessage)
}
