#!make

swag:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init -g cmd/ticket_service/main.go