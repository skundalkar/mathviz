# MathViz — developer + build-loop entry points.
# The autonomous cycle uses `make check` as its commit gate and `make digest`
# for the daily summary. You use `make serve` to look at the gallery.

GOROOT := $(shell go env GOROOT)

.PHONY: check test vet build serve digest tidy clean

## check: the gate the build loop must pass before committing (vet + tests).
check: vet test

test:
	go test ./...

vet:
	go vet ./...

## build: compile the WASM front-end and stage the loader into web/.
build:
	GOOS=js GOARCH=wasm go build -o web/main.wasm ./cmd/wasm
	cp "$(GOROOT)/lib/wasm/wasm_exec.js" web/wasm_exec.js
	@echo "Built web/main.wasm"

## serve: build, then serve the gallery at http://localhost:8080
serve: build
	@echo "Serving http://localhost:8080  (Ctrl-C to stop)"
	cd web && python3 -m http.server 8080

## digest: print the daily lesson summary (what the scheduled task sends you).
digest:
	@go run ./cmd/digest

tidy:
	go mod tidy

clean:
	rm -f web/main.wasm
