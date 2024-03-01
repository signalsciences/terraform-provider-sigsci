clean:
	rm -rf ./bin

build: clean
	go build -o bin/terraform-provider-sigsci
	@sh -c "'$(CURDIR)/scripts/generate-dev-overrides.sh'"

check:
	terraform init
	terraform plan

.PHONY: all
all: build check

lint:
	go install honnef.co/go/tools/cmd/staticcheck
	staticcheck ./...
	./scripts/gofmt.sh

testacc: ## Run acceptance tests
	TF_ACC=1 go test -v ./... $(GOTESTFLAGS)

sweep:
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	go test ./provider -v -sweep=test $(SWEEPARGS) -timeout 2m

docs:
	go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs && tfplugindocs generate

.PHONY: clean docs test
