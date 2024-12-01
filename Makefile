run:
	docker compose up --build

test:
	go test -p 1 ./... 

lint:
	golangci-lint run

fmt:
	gci write . --skip-generated -s standard -s default -s "prefix(arm-backend)" -s blank -s dot
	goimports -w .