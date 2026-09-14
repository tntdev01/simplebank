postgres:
	docker run --name postgres-db -e POSTGRES_USER=postgres_user -e POSTGRES_PASSWORD=postgres_password -e POSTGRES_DB=postgres_db -p 5432:5432 -d postgres:latest

createdb:
	docker exec -it postgres-db createdb --username=postgres_user --owner=postgres_user simple_bank

dropdb:
	docker exec -it postgres-db dropdb -U postgres_user simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://postgres_user:postgres_password@localhost:5432/simple_bank?sslmode=disable" -verbose up 

migratedown:
	migrate -path db/migration -database "postgresql://postgres_user:postgres_password@localhost:5432/simple_bank?sslmode=disable" -verbose down 

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go simple_bank/db/sqlc Store
	
.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mock