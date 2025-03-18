.PHONY: dev test

dev:
	sam build
	sam local start-api --parameter-overrides "StageName=dev FrontendUrl=http://localhost:5173"

test:
	clear
	go test ./... -v
