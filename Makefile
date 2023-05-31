CURDIR          := $(shell pwd)
GO_EXEC         := $(shell which go)
LOGNAME         := $(shell logname)
KUBEARMOR_PID    = $(shell pgrep kubearmor)
Dirs			 = $(shell ls)

ifeq (, $(shell which govvv))
$(shell go install github.com/ahmetb/govvv)	# This works for older go version
$(shell go install github.com/ahmetb/govvv@latest) # This works for new go version
endif

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifneq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

GIT_INFO := $(shell govvv -flags)

.PHONY: gofmt
## gofmt: Run gofmt linter
gofmt:
	cd $(CURDIR); gofmt -s -w -d $(shell find . -type f -name '*.go' -print)

.PHONY: golint
## golint: Run golint linter
golint:
ifeq (, $(shell which golint))
	@{ \
	set -e ;\
	GOLINT_TMP_DIR=$$(mktemp -d) ;\
	cd $$GOLINT_TMP_DIR ;\
	go mod init tmp ;\
	go get golang.org/x/lint/golint ;\
	go install golang.org/x/lint/golint ;\
	rm -rf $$GOLINT_TMP_DIR ;\
	}
endif
	cd $(CURDIR); golint ./...

.PHONY: install-addlicense
## install-addlicense: check license if not exist install addlicense tools
install-addlicense:
ifeq (, $(shell which addlicense))
	@{ \
	set -e ;\
	LICENSE_TMP_DIR=$$(mktemp -d) ;\
	cd $$LICENSE_TMP_DIR ;\
	go mod init tmp ;\
	go get github.com/google/addlicense ;\
	go install github.com/google/addlicense@latest ;\
	rm -rf $$LICENSE_TMP_DIR ;\
	}
ADDLICENSE_BIN=$(GOBIN)/addlicense
else
ADDLICENSE_BIN=$(shell which addlicense)
endif

.PHONY: filelicense
filelicense: SHELL:=/bin/bash
## filelicense: add license
filelicense:
	for file in ${Dirs} ; do \
		if [[  $$file != '_output' && $$file != 'docs' && $$file != 'vendor' && $$file != 'logger' && $$file != 'applications' ]]; then \
			$(ADDLICENSE_BIN)  -y $(shell date +"%Y") -c "Authors of KubeArmor" -f hack/LICENSE_TEMPLATE ./$$file ; \
		fi \
    done

.PHONY: gosec
## gosec: Run gosec linter
gosec:
ifeq (, $(shell which gosec))
	@{ \
	set -e ;\
	GOSEC_TMP_DIR=$$(mktemp -d) ;\
	cd $$GOSEC_TMP_DIR ;\
	go mod init tmp ;\
	go get github.com/securego/gosec/v2/cmd/gosec ;\
	go install github.com/securego/gosec/v2/cmd/gosec ;\
	rm -rf $$GOSEC_TMP_DIR ;\
	}
endif
	cd $(CURDIR); gosec ./...

.PHONY: scan
## scan: Run gosec linter
scan: 
	go install golang.org/x/vuln/cmd/govulncheck@latest ;\
	cd $(CURDIR);\
	govulncheck -v ./... ;\

## help: Display help information
help: Makefile
	@echo ""
	@echo "Usage:" "\n"
	@echo "  make [target]" "\n"
	@echo "Targets:" "\n" ""
	@awk -F ':|##' '/^[^\.%\t][^\t]*:.*##/{printf "  \033[36m%-20s\033[0m %s\n", $$1, $$NF}' $(MAKEFILE_LIST) | sort
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'