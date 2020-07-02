
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
