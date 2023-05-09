#!/bin/sh

MAKEFILE:=$(realpath $(lastword $(MAKEFILE_LIST)))
ROOT_DIRECTORY:=$(realpath $(dir $(MAKEFILE)))

################################################################################
# Main
################################################################################

up: volumes-up docker-up lambda-up terraform-up lambda-up ## Start development environment
	@echo "Starting development environment"

down: lambda-down terraform-down docker-down ## Stop development environment
	@echo "Stopping development environment"

################################################################################
# Docker
################################################################################

docker-up: ## Start docker
	@echo "Starting docker"
	@docker-compose up -d --force-recreate --build
	@echo "Docker is running"
	@echo "To stop docker, run 'make docker-down'"

docker-down: ## Stop docker
	@echo "Stopping docker"
	@docker-compose down
	@docker-compose rm -f database localstack || :
	@echo "Docker is stopped"
	@echo "To start docker, run 'make docker-up'"

docker-clean: ## Clean docker
	@echo "Cleaning docker"
	@docker-compose rm -f
	@echo "Docker is cleaned"
	@echo "To start docker, run 'make docker-up'"

################################################################################
# Terraform
################################################################################

terraform-up: ## Start terraform
	@echo "Starting terraform"
	@cd $(ROOT_DIRECTORY) && \
		terraform init && \
		terraform plan && \
		terraform apply -auto-approve && \
		echo "Terraform is running" && \
		echo "To stop terraform, run 'make terraform-down'"

terraform-down: ## Stop terraform
	@echo "Stopping terraform"
	@cd $(ROOT_DIRECTORY) && \
 		terraform destroy -auto-approve && \
		rm -rf .terraform .terrahform.* terraform.* graph.svg && \
		echo "Terraform is stopped" && \
		echo "To start terraform, run 'make terraform-up'"

terraform-reload: terraform-down terraform-up ## Restart terraform
	@echo "Restarting terraform"

################################################################################
# Terraform tflocal
################################################################################

tflocal-up: ## Start tflocal
	@echo "Starting tflocal"
	@cd $(ROOT_DIRECTORY) && \
		tflocal init && \
		tflocal plan && \
		tflocal apply -auto-approve && \
		echo "tflocal is running" && \
		echo "To stop tflocal, run 'make tflocal-down'"

tflocal-down: ## Stop tflocal
	@echo "Stopping tflocal"
	@cd $(ROOT_DIRECTORY) && \
 		tflocal destroy -auto-approve && \
		rm -rf .terraform .terrahform.* terraform.* localstack_* graph.svg && \
		echo "tflocal is stopped" && \
		echo "To start tflocal, run 'make tflocal-up'"

tflocal-reload: tflocal-down tflocal-up ## Restart tflocal
	@echo "Restarting tflocal"

################################################################################
# Lambdas
################################################################################

lambda-up: lambda-down ## build lambdas
	@echo "Building lambda executables" && \
		for dir in "Batch"; do \
		  echo "------------------- Start $$dir ------------------- "; \
		  echo "Processing directory: $$dir"; \
		  cd $(ROOT_DIRECTORY)/lambdas/$$dir; \
		  go mod init main; \
		  go mod tidy; \
		  env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(ROOT_DIRECTORY)/bin/$$dir/main main.go; \
		  cd $(ROOT_DIRECTORY)/bin/$$dir; \
          zip main.zip main; \
          cd $(ROOT_DIRECTORY)/lambdas/$$dir; \
		  cp $(ROOT_DIRECTORY)/bin/$$dir/main .; \
		  cp $(ROOT_DIRECTORY)/bin/$$dir/main.zip .; \
		  echo "$$dir lambda is built"; \
		  echo "------------------- End $$dir ------------------- "; \
  		done
	@echo "To destroy lambda executables, run 'make lambda-down'"

lambda-down: ## destroy lambda
	@echo "Removing lambda executables" && \
		cd $(ROOT_DIRECTORY) && \
		rm -rf bin/* && \
		for dir in "Batch"; do \
			rm -rf $(ROOT_DIRECTORY)/lambdas/$$dir/go.*; \
			rm -rf $(ROOT_DIRECTORY)/lambdas/$$dir/main; \
			rm -rf $(ROOT_DIRECTORY)/lambdas/$$dir/main.zip; \
			echo "$$dir lambda executable destroyed"; \
  		done

################################################################################
# Shell
################################################################################

shell-database: ## Enter shell in database container
	@echo "Entering shell in database container"
	@docker exec -it database /bin/bash

shell-localstack: ## Enter shell in localstack container
	@echo "Entering shell in localstack container"
	@docker exec -it localstack /bin/bash

################################################################################
# Mocks
################################################################################
mock-trigger-batch-lambda: ## Trigger batch lambda
	@echo "Triggering batch lambda"
	@docker exec -it localstack awslocal lambda invoke --function-name integration-marketo-batch-lambda-local --payload '{"Records": [{"body": "{\"id\": \"1\"}"}]}' /tmp/response.json

################################################################################
# Monitoring
################################################################################

monitor-docker: ## Monitor docker
	@echo "Monitoring docker"
	@docker-compose logs -f

monitor-queues: ## Monitor sqs queue
	@echo "Monitoring sqs queue"
	@docker exec -it localstack awslocal sqs list-queues

monitor-batch-queue: ## Monitor receive queue
	@echo "Monitoring batch queue"
	@docker exec -it localstack watch -n 1 awslocal sqs get-queue-attributes --queue-url http://localhost:4566/000000000000/integration-marketo-batch-queue-local --attribute-names All

monitor-lambda: ## Monitor lambda
	@echo "Monitoring lambda"
	@docker exec -it localstack awslocal lambda list-functions

################################################################################
# Helpers
################################################################################

volumes-up:
	@cd $(ROOT_DIRECTORY) && \
		mkdir -p .volumes && \
		echo "Volumes are created" && \
		echo "To destroy volumes, run 'make volumes-down'"

volumes-down:
	@cd $(ROOT_DIRECTORY) && \
		rm -rf .volumes && \
		echo "Volumes are destroyed" && \
		echo "To create volumes, run 'make volumes-up'"
