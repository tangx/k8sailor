
WORKDIR ?= cmd/k8sailor

up:
	cd $(WORKDIR) && go run .

