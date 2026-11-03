PACK            := maple
PROJECT         := github.com/MapleTechLabs/pulumi-maple
VERSION         ?= 0.0.1-dev
LDFLAGS         := -s -w -X $(PROJECT)/provider/pkg/version.Version=$(VERSION)
TFGEN           := bin/pulumi-tfgen-$(PACK)
PROVIDER        := bin/pulumi-resource-$(PACK)
SCHEMA_DIR      := provider/cmd/pulumi-resource-$(PACK)
# Keep codegen hermetic: no plugin downloads, no example cache in $HOME.
GEN_ENV         := PULUMI_HOME=$(CURDIR)/.pulumi PULUMI_DISABLE_AUTOMATIC_PLUGIN_ACQUISITION=true PULUMI_CONVERT=1
# Converts the Terraform docs' HCL examples into each SDK language.
CONVERTER       := 1.4.0
PYTHON          ?= python3

.PHONY: all tfgen provider schema generate sdk_nodejs sdk_python sdk_go build_nodejs build_python build_go test lint clean

all: generate provider

.pulumi/.converter-$(CONVERTER):
	PULUMI_HOME=$(CURDIR)/.pulumi pulumi plugin install converter terraform $(CONVERTER)
	touch $@

tfgen: .pulumi/.converter-$(CONVERTER)
	cd provider && go build -ldflags "$(LDFLAGS)" -o ../$(TFGEN) ./cmd/pulumi-tfgen-$(PACK)

provider: schema
	cd provider && go build -ldflags "$(LDFLAGS)" -o ../$(PROVIDER) ./cmd/pulumi-resource-$(PACK)

schema: tfgen
	$(GEN_ENV) $(TFGEN) schema --out $(SCHEMA_DIR)

# Everything generated is committed at the dev version; CI regenerates it and fails on
# a diff. Releases regenerate with VERSION=x.y.z, which tfgen writes into each SDK.
# Each SDK directory is cleared first so a removed resource leaves no stale file.
generate: schema sdk_nodejs sdk_python sdk_go

sdk_nodejs: tfgen
	rm -rf sdk/nodejs
	$(GEN_ENV) $(TFGEN) nodejs --out sdk/nodejs/

sdk_python: tfgen
	rm -rf sdk/python
	$(GEN_ENV) $(TFGEN) python --out sdk/python/

sdk_go: tfgen
	rm -rf sdk/go
	$(GEN_ENV) $(TFGEN) go --out sdk/go/

build_nodejs: sdk_nodejs
	cd sdk/nodejs && npm install --no-audit --no-fund && npx tsc
	cp sdk/nodejs/package.json sdk/nodejs/bin/

build_python: sdk_python
	cp README.md sdk/python/
	cd sdk/python && $(PYTHON) -m venv .venv && .venv/bin/pip install -q build && .venv/bin/python -m build --outdir bin .

build_go: sdk_go
	cd sdk && go build ./go/...

test:
	cd provider && go test -race -count=1 ./...

lint:
	cd provider && test -z "$$(gofmt -l .)" && go vet ./...

clean:
	rm -rf bin .pulumi sdk/nodejs/bin sdk/nodejs/node_modules sdk/python/bin sdk/python/.venv
