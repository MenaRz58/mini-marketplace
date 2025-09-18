.PHONY: run-products run-users run-orders all

run-products:
	go run ./products/cmd/main

run-users:
	go run ./users/cmd/main

run-orders:
	go run ./orders/cmd/main

all:
	@echo "Start in separate terminals:"
	@echo " make run-products"
	@echo " make run-users"
	@echo " make run-orders"
