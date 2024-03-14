.PHONY: help

.DEFAULT_GOAL := gflow

.PHONY: help
help:
	@echo ""
	@echo "gflow: Build gflow binary"
	@echo "help:  Show this help"
	@echo ""

gflow:
	go mod download
	go build

.PHONY: clean
clean:
	rm gflow
