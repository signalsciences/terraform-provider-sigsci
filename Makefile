
build:
	go build -o terraform-provider-sigsci

check:
	terraform init
	terraform plan

all: build check

lint:
	golint ./...

testacc: ## Run acceptance tests
	TF_ACC=1 go test -v ./...

sweep:
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	go test ./provider -v -sweep=test $(SWEEPARGS) -timeout 2m