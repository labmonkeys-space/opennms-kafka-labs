.PHONY: help

.DEFAULT_GOAL := gevent

.PHONY: help
help:
	@echo ""
	@echo "gevent: Build gevent binary"
	@echo "help:  Show this help"
	@echo ""

gevent:
	go mod download
	go build

.PHONY: clean
clean:
	rm gevent
