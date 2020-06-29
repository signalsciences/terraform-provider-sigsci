
build:
	go build -o terraform-provider-sigsci

check:
	terraform init
	terraform plan

all: build check

lint:
	golint . ./provider