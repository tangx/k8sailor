
WORKDIR ?= cmd/k8sailor

up:
	cd $(WORKDIR) && go run .
httpserver:
	cd $(WORKDIR) && go run . httpserver

help:
	cd $(WORKDIR) && go run . help
