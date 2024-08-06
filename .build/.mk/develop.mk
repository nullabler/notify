##@ —— Develop  ———————————————————————————————————————————————————————

gateway: ## local go run gateway
	docker-compose exec notify-gateway go run cmd/gateway/gateway.go --path=./config/dev.yaml

consumer: ## local go run consumer
	docker-compose exec notify-consumer go run cmd/consumer/consumer.go --path=./config/dev.yaml

send: ## to send request for test notify
	curl -X POST http://localhost:8081/send/pipeline-stage -d '{"username": "user_1", "state": "Start", "build-number": "123"}'

ping: ## ping/pong gateway
	curl http://localhost:8081/ping
