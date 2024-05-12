.PHONY: run/snip
run/snip:
	@go run ./cmd/web/main.go

.PHONY: run/exe
run/exe:
	@./bin/api.exe

.PHONY: build/api
build/api:
	@echo 'Building cmd/api...'
	go build -ldflags '-s' -o ./bin/api.exe ./cmd/web
	@make run/exe