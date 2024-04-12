##@ —— System 🐳 ——————————————————————————————————————————————————————

build: ## docker build
	docker build --platform=linux/amd64 -t geticit/notify:0.2 -f Dockerfile.builder

up: ## up project [make up] [make up build=1 watch=1]
ifdef build
	$(eval OPTS=${OPTS} --build)
	sudo rm -rf .build/kafka/data
endif
ifdef watch
else
	$(eval OPTS=${OPTS} -d)
endif
	docker-compose up ${OPTS} --remove-orphans

down: ## down project [make down]
	docker-compose down --remove-orphans


start: ## start service [make start notify]
	docker-compose start $(ARGS)

stop: ## stop service [make stop notify]
	docker-compose stop $(ARGS)

exec: ## call exec service [make exec notify go version]
ifdef su
	$(eval OPTS=${OPTS} -u root)
endif
	docker-compose exec ${OPTS} $(ARGS)

run: ## call run service [make run notify go version]
	docker-compose run $(ARGS)

ps: ## show ps [make ps]
	docker-compose ps

logs: ## show last 100 lines of logs [make logs notify]
	docker-compose logs --tail=100 $(ARGS)

restart: ## restart service [make restart notify]
	docker-compose restart $(ARGS)

test-api: ## to send request for test notify
	curl -X POST http://localhost:8081/send/pipeline-stage -d '{"username": "user_1", "state": "Start", "build-number": "123"}'

