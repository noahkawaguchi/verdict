.PHONY: dev

dev:
	sam build && sam local start-api --parameter-overrides "StageName=dev FrontendUrl=http://localhost:5173"
