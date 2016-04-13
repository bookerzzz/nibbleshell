## Nibbleshell build targets
## See also README.md

## list of buildable Go packages
PKGS := $(shell glide nv)

OK_COLOR=\033[32;01m
NO_COLOR=\033[0m

build:
	@echo "$(OK_COLOR)==> Compiling binary$(NO_COLOR)"
	mkdir -p dist/bin
	GOBIN="$(CURDIR)/dist/bin" go install $(PKGS)

clean:
	@rm -rf dist/bin/

vendor:
	@echo "$(OK_COLOR)==> Installing dependencies$(NO_COLOR)"
	glide up --no-recursive

.PHONY: clean vendor build
