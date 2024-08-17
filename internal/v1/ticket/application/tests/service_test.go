package tests

import (
	"context"
	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"log"
	main "main/config"
	"main/internal/v1/ticket/application/services"
	"main/internal/v1/ticket/domain/entities"
	"main/internal/v1/ticket/infrastructure/repository"
	"main/pkg/databases/postgresql"
	"main/pkg/logger"
	"os"
	"testing"
	"time"
)

var cfg *main.Config
var err error
var mockDb *embeddedpostgres.EmbeddedPostgres
var pgConn *pgxpool.Pool

func InitializeDbConnection() {
	if err = mockDb.Start(); err != nil {
		log.Fatal(err.Error())
	}

	// Init Clients
	pgConn, err = postgresql.NewPostgresqlDB(cfg)
	if err != nil {
		log.Fatal("Error when tyring to connect to Postgresql")
	} else {
		log.Print("Postgresql connected")
	}

	tables := `
		--region tickets
		CREATE TABLE tickets
		(
			id INT8 GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
			name VARCHAR(75) NOT NULL,
			description VARCHAR(255) NOT NULL,
			allocation INT8 NOT NULL CHECK(allocation>= 0) DEFAULT 0
		);
		COMMENT ON TABLE tickets IS 'Tickets';
		CREATE INDEX ON tickets USING btree (name);
		CREATE INDEX ON tickets USING btree (name, description);
		
		--endregion
		
		--region purchase_logs
		CREATE TABLE public.purchase_logs
		(
			log_id INT8 GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
			ticket_id INT8 NOT NULL,
			user_id VARCHAR(255) NOT NULL,
			quantity INT8 NOT NULL DEFAULT 0
		);
		COMMENT ON TABLE public.purchase_logs IS 'Purchase Logs';
		CREATE INDEX ON public.purchase_logs USING btree (ticket_id);
		
		--endregion`
	_, err = pgConn.Exec(context.Background(), tables)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func init() {
	os.Setenv("APP_ENV", "dev")

	cfg, err = main.ParseMockConfig("../../../../../../config/config.dev.yaml")
	if err != nil {
		log.Fatal(err)
	}

	mockDbConfig := embeddedpostgres.DefaultConfig().
		Username(cfg.Postgresql.USER).
		Password(cfg.Postgresql.PASS).
		Database(cfg.Postgresql.DEFAULT_DB).
		Version(embeddedpostgres.V16).
		Port(uint32(cfg.Postgresql.PORT)).
		StartTimeout(45 * time.Second)

	mockDb = embeddedpostgres.NewDatabase(mockDbConfig)

}

// TestCreateTicketSuccess Tests Can a ticket be created
func TestCreateTicketSuccess(t *testing.T) {
	var ticket entities.TicketResDto
	var recordCreate int64
	assert := assert.New(t)
	InitializeDbConnection()
	defer mockDb.Stop()

	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()

	mockPgRepo := repository.NewPostgresqlRepository(context.Background(), pgConn, appLogger)
	mockTicketService := services.NewTicketService(cfg, mockPgRepo, appLogger)
	recordCreate, ticket = mockTicketService.CreateTicket(entities.CreateTicketReqDto{
		Name:       "Test",
		Desc:       "test description",
		Allocation: 100,
	})

	assert.Equal(int64(1), recordCreate, "Record has to be 1")
	assert.Equal("Test", ticket.Name, "Name is not true")
	assert.Equal("test description", ticket.Desc, "Desc is not true")
	assert.Equal(uint64(100), ticket.Allocation, "Allocation is not true")
}

// TestGetTicketSuccess Tests Can a ticket be created and selected
func TestGetTicketSuccess(t *testing.T) {
	var ticket entities.TicketResDto
	var recordCreate, recordGet int64

	assert := assert.New(t)
	InitializeDbConnection()
	defer mockDb.Stop()

	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()

	mockPgRepo := repository.NewPostgresqlRepository(context.Background(), pgConn, appLogger)
	mockTicketService := services.NewTicketService(cfg, mockPgRepo, appLogger)
	recordCreate, ticket = mockTicketService.CreateTicket(entities.CreateTicketReqDto{
		Name:       "Test",
		Desc:       "test description",
		Allocation: 100,
	})
	recordGet, ticket = mockTicketService.GetTicket(ticket.Id)

	assert.Equal(int64(1), recordCreate, "Record has to be 1")
	assert.Equal(int64(1), recordGet, "Record has to be 1")
	assert.Equal("Test", ticket.Name, "Name is not true")
	assert.Equal("test description", ticket.Desc, "Desc is not true")
	assert.Equal(uint64(100), ticket.Allocation, "Allocation is not true")
}

// TestGetTicketFailInvalidTicketId Tests Can a ticket be created and selected
func TestGetTicketFailInvalidTicketId(t *testing.T) {
	var recordGet int64

	assert := assert.New(t)
	InitializeDbConnection()
	defer mockDb.Stop()

	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()

	mockPgRepo := repository.NewPostgresqlRepository(context.Background(), pgConn, appLogger)
	mockTicketService := services.NewTicketService(cfg, mockPgRepo, appLogger)
	recordGet, _ = mockTicketService.GetTicket(2)

	assert.Equal(int64(-1), recordGet, "Record has to be -1")
}

// TestMakePaymentSuccess Tests Can a ticket be created, selected and purchased (allocation has to be dropped by the quantity)
func TestMakePaymentSuccess(t *testing.T) {
	var ticket entities.TicketResDto
	var recordCreate, recordGet, recordPurchase int64

	assert := assert.New(t)
	InitializeDbConnection()
	defer mockDb.Stop()

	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()

	mockPgRepo := repository.NewPostgresqlRepository(context.Background(), pgConn, appLogger)
	mockTicketService := services.NewTicketService(cfg, mockPgRepo, appLogger)
	recordCreate, ticket = mockTicketService.CreateTicket(entities.CreateTicketReqDto{
		Name:       "Test",
		Desc:       "test description",
		Allocation: 100,
	})

	recordPurchase = mockTicketService.MakePurchase(ticket.Id, entities.MakePurchaseReqDto{
		Quantity: 2,
		UserId:   uuid.NewString(),
	})
	recordGet, ticket = mockTicketService.GetTicket(ticket.Id)

	assert.Equal(int64(1), recordCreate, "Record has to be 1")
	assert.Equal(int64(1), recordGet, "Record has to be 1")
	assert.Equal(int64(1), recordPurchase, "Record has to be 1")
	assert.Equal("Test", ticket.Name, "Name is not true")
	assert.Equal("test description", ticket.Desc, "Desc is not true")
	assert.Equal(uint64(98), ticket.Allocation, "Allocation is not true")
}

// TestMakePaymentFailOnInvalidAllocation Tests Can a ticket be created, selected, but purchase won't be proceeded because of the invalid quantity
func TestMakePaymentFailInvalidAllocation(t *testing.T) {
	var ticket entities.TicketResDto
	var recordCreate, recordGet, recordPurchase int64

	assert := assert.New(t)
	InitializeDbConnection()
	defer mockDb.Stop()

	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()

	mockPgRepo := repository.NewPostgresqlRepository(context.Background(), pgConn, appLogger)
	mockTicketService := services.NewTicketService(cfg, mockPgRepo, appLogger)
	recordCreate, ticket = mockTicketService.CreateTicket(entities.CreateTicketReqDto{
		Name:       "Test",
		Desc:       "test description",
		Allocation: 100,
	})

	recordPurchase = mockTicketService.MakePurchase(ticket.Id, entities.MakePurchaseReqDto{
		Quantity: 102,
		UserId:   uuid.NewString(),
	})

	recordGet, ticket = mockTicketService.GetTicket(ticket.Id)

	assert.Equal(int64(1), recordCreate, "Record has to be 1")
	assert.Equal(int64(1), recordGet, "Record has to be 1")
	assert.Equal(int64(-3), recordPurchase, "Record has to be -3")
	assert.Equal("Test", ticket.Name, "Name is not true")
	assert.Equal("test description", ticket.Desc, "Desc is not true")
	assert.Equal(uint64(100), ticket.Allocation, "Allocation is not true")
}

// TestGetPaymentLogs Tests Can a ticket will be created, selected and purchased. If there is no issue, A purchase log has to be created and can be reachable by ticket_id
func TestGetPaymentLogs(t *testing.T) {
	var ticket entities.TicketResDto
	var recordCreate, recordGet, recordPurchase int64

	assert := assert.New(t)
	InitializeDbConnection()
	defer mockDb.Stop()

	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()

	mockPgRepo := repository.NewPostgresqlRepository(context.Background(), pgConn, appLogger)
	mockTicketService := services.NewTicketService(cfg, mockPgRepo, appLogger)
	recordCreate, ticket = mockTicketService.CreateTicket(entities.CreateTicketReqDto{
		Name:       "Test",
		Desc:       "test description",
		Allocation: 100,
	})

	recordPurchase = mockTicketService.MakePurchase(ticket.Id, entities.MakePurchaseReqDto{
		Quantity: 2,
		UserId:   uuid.NewString(),
	})

	recordGet, ticket = mockTicketService.GetTicket(ticket.Id)

	recordLogs, logs := mockTicketService.GetPurchaseLogs(ticket.Id)
	assert.Equal(int64(1), recordCreate, "Record has to be 1")
	assert.Equal(int64(1), recordGet, "Record has to be 1")
	assert.Equal(int64(1), recordPurchase, "Record has to be 1")
	assert.Equal(int64(1), recordLogs, "Record has to be 1")
	assert.Equal("Test", ticket.Name, "Name is not true")
	assert.Equal("test description", ticket.Desc, "Desc is not true")
	assert.Equal(uint64(98), ticket.Allocation, "Allocation is not true")
	assert.NotEqual("[]", string(logs), "Purchase log has to be different from square brackets")
}
