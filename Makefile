default: help


swagger/doc: ## Generate Swagger specs. Requires swag.
	@echo "Generating specs..."
	@swag init -d ./api/ -s ./api/docs/swagger

go/test: ## Run the tests
	@cd api && go test -v

docker/build: ## Build the docker images
	@echo "Building using docker-compose..."
	@docker-compose build

docker/deploy: ## Deploy the containers using docker-compose
	@docker-compose up -d

docker/destroy: ## Stop and remove the containers
	@docker-compose down

help: ## This is default and it helps
	@echo "\n  ## GOLANG-REDIS-COMPOSE\n"
	@awk 'BEGIN {FS = ":.*?## "} /^[\/a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "  \033[36mgolang-redis-compose =>    make \033[0m%-25s %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo
