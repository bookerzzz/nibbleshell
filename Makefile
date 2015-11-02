OK_COLOR=\033[32;01m
NO_COLOR=\033[0m

build:
	@echo "$(OK_COLOR)==> Compiling binary$(NO_COLOR)"
	mkdir -p bin
	GOBIN=bin/ go install

clean:
	@rm -rf bin/
	@rm -rf result/

deps:
	@echo "$(OK_COLOR)==> Installing dependencies$(NO_COLOR)"
	glide up

format:
	go fmt ./...

.PHONY: clean format deps build
