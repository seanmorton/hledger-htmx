.PHONY: build
build:
	docker build --tag smorton517/hledger-htmx --platform linux/amd64 .

.PHONY: push
push:
	docker image push smorton517/hledger-htmx:latest

