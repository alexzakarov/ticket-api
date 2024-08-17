package repository

import (
	"context"
	"github.com/goccy/go-json"
	"github.com/jackc/pgx/v4/pgxpool"
	"main/internal/v1/ticket/domain/entities"
	"main/internal/v1/ticket/domain/ports"
	"main/pkg/logger"
	cm "main/pkg/utils/common"
)

// postgresqlRepo Struct
type postgresqlRepo struct {
	ctx    context.Context
	db     *pgxpool.Pool
	logger *logger.ApiLogger
}

// NewPostgresqlRepository Ticket postgresql repository constructor
func NewPostgresqlRepository(ctx context.Context, db *pgxpool.Pool, logger *logger.ApiLogger) ports.IPostgresqlRepository {
	return &postgresqlRepo{
		ctx:    ctx,
		db:     db,
		logger: logger,
	}
}

func (p *postgresqlRepo) CreateTicket(req_dto entities.CreateTicketReqDto) (record int64, data entities.TicketResDto) {
	var errDb error
	var query string
	record = 1

	query = `INSERT INTO public.tickets (name, description, allocation) VALUES (trim($1),$2,$3) RETURNING id, name, description, allocation`
	errDb = p.db.QueryRow(p.ctx, query, req_dto.Name, req_dto.Desc, req_dto.Allocation).Scan(&data.Id, &data.Name, &data.Desc, &data.Allocation)
	if errDb != nil && cm.CheckStringIfContains(errDb.Error(), "duplicate key value") == false {
		record = -1
		p.logger.Error("CreateTicket: ", errDb.Error())
	} else if errDb != nil && cm.CheckStringIfContains(errDb.Error(), "duplicate key value") == true {
		record = -2
		p.logger.Error("CreateTicket: ", errDb.Error())
	}

	return
}

func (p *postgresqlRepo) GetTicket(ticket_id int64) (record int64, data entities.TicketResDto) {
	var errDb error
	var query string
	record = 1

	query = `SELECT id, name, description, allocation FROM public.tickets WHERE id=$1`
	errDb = p.db.QueryRow(p.ctx, query, ticket_id).Scan(&data.Id, &data.Name, &data.Desc, &data.Allocation)
	if errDb != nil && cm.CheckStringIfContains(errDb.Error(), "no rows in result set") == false {
		record = -1
		p.logger.Error("GetTicket: ", errDb.Error())
	} else if errDb != nil && cm.CheckStringIfContains(errDb.Error(), "no rows in result set") == true {
		record = 0
		p.logger.Error("GetTicket: ", errDb.Error())
	}

	return
}

func (p *postgresqlRepo) MakePurchase(ticket_id int64, req_dto entities.MakePurchaseReqDto) (record int64) {
	var errDb error
	var query string
	record = 1

	query = `UPDATE public.tickets SET allocation=(allocation-$1) WHERE id=$2`
	_, errDb = p.db.Exec(p.ctx, query, req_dto.Quantity, ticket_id)
	if errDb != nil && cm.CheckStringIfContains(errDb.Error(), "no rows in result set") == false {
		record = -1
		p.logger.Error("MakePurchase: ", errDb.Error())
	} else if errDb != nil && cm.CheckStringIfContains(errDb.Error(), "no rows in result set") == true {
		record = 0
		p.logger.Error("MakePurchase: ", errDb.Error())
	}

	return
}

func (p *postgresqlRepo) AddPurchaseLog(ticket_id int64, req_dto entities.MakePurchaseReqDto) (record int64) {
	var errDb error
	var query string
	record = 1

	query = `INSERT INTO public.purchase_logs (ticket_id, user_id, quantity) VALUES ($1,$2,$3)`
	_, errDb = p.db.Exec(p.ctx, query, ticket_id, req_dto.UserId, req_dto.Quantity)
	if errDb != nil && cm.CheckStringIfContains(errDb.Error(), "duplicate key value") == false {
		record = -1
		p.logger.Error("AddPurchaseLog: ", errDb.Error())
	} else if errDb != nil && cm.CheckStringIfContains(errDb.Error(), "duplicate key value") == true {
		record = -2
		p.logger.Error("AddPurchaseLog: ", errDb.Error())
	}

	return
}

func (p *postgresqlRepo) GetPurchaseLogs(ticket_id int64) (record int64, data json.RawMessage) {
	var errDb error
	var query string
	record = 1

	query = `SELECT COALESCE(JSON_AGG(LOGS),'[]') FROM (SELECT log_id, ticket_id, user_id, quantity FROM public.purchase_logs WHERE ticket_id=$1) LOGS`
	errDb = p.db.QueryRow(p.ctx, query, ticket_id).Scan(&data)
	if errDb != nil && cm.CheckStringIfContains(errDb.Error(), "no rows in result set") == false {
		record = -1
		p.logger.Error("GetPurchaseLogs: ", errDb.Error())
	} else if errDb != nil && cm.CheckStringIfContains(errDb.Error(), "no rows in result set") == true {
		record = 0
		p.logger.Error("GetPurchaseLogs: ", errDb.Error())
	}

	return
}
