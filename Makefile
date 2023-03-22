include .env
export $(shell sed 's/=.*//' .env)

default:

run:
	@go run cmd/main.go

tidy:
	@go mod tidy
