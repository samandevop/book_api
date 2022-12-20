

go:
	go run cmd/main.go

swag-init:
	swag init -g api/api.go -o api/docs

migration-up:
	migrate -path ./migrations/postgres/ -database 'postgres://samandar:samandevop@localhost:5432/book?sslmode=disable' up

migration-down:
	migrate -path ./migrations/postgres/ -database 'postgres://samandar:samandevop@localhost:5432/book?sslmode=disable' down
