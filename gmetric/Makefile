.PHONY: help

.DEFAULT_GOAL := gmetric

.PHONY: help
help:
	@echo ""
	@echo "gmetric: Build gmetric binary"
	@echo "help:  Show this help"
	@echo ""

gmetric:
	go mod download
	go build

.PHONY: clean
clean:
	rm gmetric
