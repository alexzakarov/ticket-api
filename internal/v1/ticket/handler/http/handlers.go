package http

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"main/config"
	"main/internal/v1/ticket/domain/entities"
	"main/internal/v1/ticket/domain/ports"
	"main/pkg/logger"
	cm "main/pkg/utils/common"
	"main/pkg/utils/typeconv"
	"main/pkg/utils/validator"
)

// handlerHttp Tickets handlers
type handlerHttp struct {
	ctx     context.Context
	cfg     *config.Config
	service ports.IService
	logger  logger.Logger
}

// NewHttpHandler Tickets HTTP handlers constructor
func NewHttpHandler(ctx context.Context, cfg *config.Config, service ports.IService, logger logger.Logger) ports.IHandlers {
	return &handlerHttp{
		ctx:     ctx,
		cfg:     cfg,
		service: service,
		logger:  logger,
	}
}

// CreateTicket function Creates an event with an allocation of tickets available to purchase.
// @Description "CreateTicket Function". Creates an allocation of tickets available to purchase.
// @Summary CreateTicket Function Creates a ticket
// @Tags Tickets
// @Accept json
// @Produce json
// @Param CreateTicket body entities.CreateTicketReqDto true "CreateTicket"
// @Success 200 {object} entities.HandlerResponse
// @Failure 400 {object} entities.HandlerResponse
// @Router /tickets [post]
func (h handlerHttp) CreateTicket(c *fiber.Ctx) error {
	var responser fiber.Map
	var statusCode int
	var record int64
	var errValidate error
	var data interface{}
	dat := entities.CreateTicketReqDto{}

	errParser := c.BodyParser(&dat)
	if errParser != nil {
		statusCode, responser = fiber.StatusBadRequest, cm.HTTPResponser(nil, true, errParser.Error())
	}

	errValidate = validator.Validate.Struct(dat)
	if errValidate != nil {
		statusCode, responser = fiber.StatusBadRequest, cm.HTTPResponser(nil, true, errValidate.Error())
	}

	if errParser == nil && errValidate == nil {
		record, data = h.service.CreateTicket(dat)
		if record == -2 {
			statusCode, responser = fiber.StatusConflict, cm.HTTPResponser(nil, true, cm.Translate(c, "General.ERROR_CONFLICT"))
		} else if record == -1 {
			statusCode, responser = fiber.StatusInternalServerError, cm.HTTPResponser(nil, true, cm.Translate(c, "General.ERROR_DB"))
		} else {
			statusCode, responser = fiber.StatusOK, cm.HTTPResponser(data, false, cm.Translate(c, "General.OK"))
		}
	}

	return c.Status(statusCode).JSON(responser)
}

// GetTicket function Gets an event with an allocation of tickets available to purchase.
// @Description "GetTicket Function". Gets an allocation of tickets available to purchase.
// @Summary GetTicket Function Gets a ticket
// @Tags Tickets
// @Accept json
// @Produce json
// @Param id path int true "GetTicket"
// @Success 200 {object} entities.HandlerResponse
// @Failure 400 {object} entities.HandlerResponse
// @Router /tickets/{id} [post]
func (h handlerHttp) GetTicket(c *fiber.Ctx) error {
	var responser fiber.Map
	var statusCode int
	var record int64
	var data interface{}

	id := typeconv.StrToInt64(c.Params("id"))

	record, data = h.service.GetTicket(id)
	if record == -1 {
		statusCode, responser = fiber.StatusConflict, cm.HTTPResponser(nil, true, cm.Translate(c, "General.ERROR_DB"))
	} else if record == 0 {
		statusCode, responser = fiber.StatusInternalServerError, cm.HTTPResponser(nil, true, cm.Translate(c, "General.ERROR_NO_ROW"))
	} else {
		statusCode, responser = fiber.StatusOK, cm.HTTPResponser(data, false, cm.Translate(c, "General.OK"))
	}

	return c.Status(statusCode).JSON(responser)
}

// MakePurchase function Creates a purchase call to reduce allocation from ticket and create its log
// @Description "CreateTicket Function". Creates a purchase call to reduce allocation from ticket and create its log
// @Summary MakePurchase Function Creates a purchase call
// @Tags Purchases
// @Accept json
// @Produce json
// @Param id path int true "MakePurchase"
// @Param MakePurchase body entities.MakePurchaseReqDto true "MakePurchase"
// @Success 200 {object} entities.HandlerResponse
// @Failure 400 {object} entities.HandlerResponse
// @Router /tickets/{id}/purchases [post]
func (h handlerHttp) MakePurchase(c *fiber.Ctx) error {
	var responser fiber.Map
	var statusCode int
	var record int64
	var errValidate error
	dat := entities.MakePurchaseReqDto{}

	id := typeconv.StrToInt64(c.Params("id"))

	errParser := c.BodyParser(&dat)
	if errParser != nil {
		statusCode, responser = fiber.StatusBadRequest, cm.HTTPResponser(nil, true, errParser.Error())
	}

	errValidate = validator.Validate.Struct(dat)
	if errValidate != nil {
		statusCode, responser = fiber.StatusBadRequest, cm.HTTPResponser(nil, true, errValidate.Error())
	}

	if errParser == nil && errValidate == nil {
		record = h.service.MakePurchase(id, dat)
		if record == -3 {
			statusCode, responser = fiber.StatusConflict, cm.HTTPResponser(nil, true, cm.Translate(c, "Purchase.ERROR_INVALID_AMOUNT"))
		} else if record == -2 {
			statusCode, responser = fiber.StatusConflict, cm.HTTPResponser(nil, true, cm.Translate(c, "General.ERROR_CONFLICT"))
		} else if record == -1 {
			statusCode, responser = fiber.StatusInternalServerError, cm.HTTPResponser(nil, true, cm.Translate(c, "General.ERROR_DB"))
		} else {
			statusCode, responser = fiber.StatusOK, cm.HTTPResponser(nil, false, cm.Translate(c, "General.OK"))
		}
	}

	return c.Status(statusCode).JSON(responser)
}

// GetPurchaseLogs function Gets all proceeded purchases of a ticket
// @Description "GetPurchaseLogs Function". Gets all proceeded purchases of a ticket
// @Summary GetPurchaseLogs Function Gets all purchases
// @Tags Purchases
// @Accept json
// @Produce json
// @Param id path int true "GetPurchaseLogs"
// @Success 200 {object} entities.HandlerResponse
// @Failure 400 {object} entities.HandlerResponse
// @Router /tickets/{id}/purchases [get]
func (h handlerHttp) GetPurchaseLogs(c *fiber.Ctx) error {
	var responser fiber.Map
	var statusCode int
	var record int64
	var data interface{}

	id := typeconv.StrToInt64(c.Params("id"))

	record, data = h.service.GetPurchaseLogs(id)
	if record == -1 {
		statusCode, responser = fiber.StatusConflict, cm.HTTPResponser(nil, true, cm.Translate(c, "General.ERROR_DB"))
	} else if record == 0 {
		statusCode, responser = fiber.StatusInternalServerError, cm.HTTPResponser(nil, true, cm.Translate(c, "General.ERROR_NO_ROW"))
	} else {
		statusCode, responser = fiber.StatusOK, cm.HTTPResponser(data, false, cm.Translate(c, "General.OK"))
	}

	return c.Status(statusCode).JSON(responser)
}
