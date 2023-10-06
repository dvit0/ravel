run-worker: build-drivers build-init
	sudo go run cmd/worker/main.go

build-drivers:
	go build -o bin/drivers/firecracker-driver drivers/firecracker-driver/*.go
build-init:
	CGO_ENABLED=0 go build -o bin/init cmd/init/main.go

driver-protoc:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    pkg/driver/proto/driver.proto