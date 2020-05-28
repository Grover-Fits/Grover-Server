.DEFAULT_GOAL := help

SERVER_PID := grover.pid
CLIENT_PID := client.pid
.PHONY: 

RESET  = \033[0m
PURPLE = \033[0;35m
GREEN  = \033[0;32m
LINE   = $(PURPLE)-------------------------------------------------------------------------------------------------$(RESET)

build: ## build grpc server
	go build

start-server: build ## start grpc server
	./grover & echo $$! > $(SERVER_PID);

start-client: ## start client
	python2 -m SimpleHTTPServer 8085 & echo $$! > $(CLIENT_PID);

start-all: start-server start-client ## start grpc server and client

stop-server: ## stop grpc server
	kill `cat $(SERVER_PID)` && rm -rf $(SERVER_PID)

stop-client: ## stop client
	kill `cat $(CLIENT_PID)` && rm -rf $(CLIENT_PID)

stop-all: stop-server stop-client ## stop grpc server and client

help: ## That's me!
	@echo
	@echo "#$(LINE)"
	@printf "\033[37m%-30s\033[0m %s\n" "# Makefile Help                                                                                  |"
	@echo "#$(LINE)"
	@printf "\033[37m%-30s\033[0m %s\n" "# This Makefile can be used to run, build, and tear down the atlassian suite (Confluence & Jira) |"
	@echo "#$(LINE)"
	@echo 
	@printf "\033[37m%-30s\033[0m %s\n" "#-target-----------------------description--------------------------------------------------------"
	@grep -E '^[a-zA-Z_-].+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	@echo 

print-%  : ; @echo $* = $($*)