package services

import (
	"github.com/goccy/go-json"
	"main/config"
	"main/internal/v1/ticket/domain/entities"
	"main/internal/v1/ticket/domain/ports"
	"main/pkg/logger"
	"sync"
)

// serviceTicket Ticket Service
type serviceTicket struct {
	cfg    *config.Config
	pgRepo ports.IPostgresqlRepository
	logger logger.Logger
}

// NewTicketService service constructor
func NewTicketService(cfg *config.Config, pgRepo ports.IPostgresqlRepository, logger logger.Logger) ports.IService {
	return &serviceTicket{
		cfg:    cfg,
		pgRepo: pgRepo,
		logger: logger,
	}
}

func (s *serviceTicket) CreateTicket(req_dto entities.CreateTicketReqDto) (record int64, data entities.TicketResDto) {

	record, data = s.pgRepo.CreateTicket(req_dto)

	return
}

func (s *serviceTicket) GetTicket(ticket_id int64) (record int64, data entities.TicketResDto) {

	record, data = s.pgRepo.GetTicket(ticket_id)

	return
}

func (s *serviceTicket) MakePurchase(ticket_id int64, req_dto entities.MakePurchaseReqDto) (record int64) {
	var ticketInfo entities.TicketResDto
	mutex := sync.RWMutex{}

	record, ticketInfo = s.pgRepo.GetTicket(ticket_id)
	if record == 1 && ticketInfo.Allocation >= req_dto.Quantity {

		mutex.Lock()
		record = s.pgRepo.MakePurchase(ticket_id, req_dto)
		mutex.Unlock()
		if record == 1 {
			s.pgRepo.AddPurchaseLog(ticket_id, req_dto)
		} else {
			record = -2
		}
	} else {
		record = -3
	}

	return
}

func (s *serviceTicket) GetPurchaseLogs(ticket_id int64) (record int64, data json.RawMessage) {

	record, data = s.pgRepo.GetPurchaseLogs(ticket_id)

	return
}
