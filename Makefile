postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=Lc@200409200034 -d postgres
createdb:
	docker exec -it postgres createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres dropdb simple_bank
migrateup:
	migrate -path db/migration -database "postgresql://root:Lc@200409200034@localhost:5432/simple_bank?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://root:Lc@200409200034@localhost:5432/simple_bank?sslmode=disable" -verbose down
# 只向上或向下回滚一个版本
migrateup1:
	migrate -path db/migration -database "postgresql://root:Lc@200409200034@localhost:5432/simple_bank?sslmode=disable" -verbose up 1
migratedown1:
	migrate -path db/migration -database "postgresql://root:Lc@200409200034@localhost:5432/simple_bank?sslmode=disable" -verbose down 1
sqlc: 
	sqlc generate
test:
	go test -v -cover ./...
server:
	go run main.go
mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/Long4Changes/MySimpleBank/db/sqlc Store 
# 每次添加好一个指令都要加到下面去
.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mock migrateup1 migratedown1