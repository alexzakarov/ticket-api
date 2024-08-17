package entities

type HandlerResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type UserData struct {
	UserId int64 `json:"user_id,omitempty"`
}

type CreateTicketReqDto struct {
	Name       string `json:"name" validate:"min=1,max=75"`
	Desc       string `json:"desc" validate:"min=1,max=255"`
	Allocation uint64 `json:"allocation" validate:"required,gt=0"`
}

type TicketResDto struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	Desc       string `json:"desc"`
	Allocation uint64 `json:"allocation"`
}

type MakePurchaseReqDto struct {
	UserId   string `json:"user_id" validate:"min=1,max=200"`
	Quantity uint64 `json:"quantity" validate:"required,gt=0"`
}
