APP_NAME := server

run-serve:
	go run ./cmd/${APP_NAME}

.PHONY: run

.PHONY: air