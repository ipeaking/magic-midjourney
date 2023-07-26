PROJECT="wrap-midjourney"
BUILDTIME=`date '+%Y%m%d%H%M'`

linux: ## Build for linux
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/${PROJECT}.${BUILDTIME} ./main.go

help: ## Display this help message
	@cat $(MAKEFILE_LIST) | grep -e "^[a-zA-Z_\-]*: *.*## *" | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

deploy: ## Deploy
	@echo "Deploying ${PROJECT}..."
	@echo "Building for linux..."
	@make linux
	@echo "Deploying to server..."
	@mv /opt/${PROJECT}/${PROJECT} /opt/${PROJECT}/${PROJECT}.bak
	@cp bin/${PROJECT}.${BUILDTIME} /opt/${PROJECT}/${PROJECT}
	@echo "Restarting service..."
	@supervisorctl restart ${PROJECT}
	@echo "Done!"

.DEFAULT_GOAL := help

run: ## Run
	go run cmd/${PROJECT}/main.go

.SILENT: linux

.PHONY: all linux