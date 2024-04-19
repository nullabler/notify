##@ —— Develop  ———————————————————————————————————————————————————————

test-api: ## to send request for test notify
	curl -X POST http://localhost:8081/send/pipeline-stage -d '{"username": "user_1", "state": "Start", "build-number": "123"}'

gateway: ## local go run gateway
	go run cmd/gateway/gateway.go --path=./config/dev.yaml
