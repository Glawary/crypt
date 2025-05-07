docker-run-postgres:
	docker run --rm --name crypt \
	-p 7890:5432 \
	-e POSTGRES_USER=postgres \
	-e POSTGRES_PASSWORD=postgres \
	-e POSTGRES_DB=crypt \
	-d postgres:15

migration-up:
	migrate -path ./migrations -database "postgresql://postgres:postgres@localhost:7890/crypt?sslmode=disable" -verbose up

proto-gen:
	protoc \
     -I ./proto \
     --go_out=../ \
     --go-grpc_out=../ \
     proto/crypt.proto
.PHONY: proto-gen

docker-image:
	docker build -t goproject .

docker-run-project:
	docker run --name gpt \
		-p 8090:8090 \
		-d goproject
.PHONY: dkt