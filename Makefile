.DEFAULT_GOAL := help

SERVER_PID := grover.pid
.PHONY: 

RESET  = \033[0m
PURPLE = \033[0;35m
GREEN  = \033[0;32m
LINE   = $(PURPLE)-------------------------------------------------------------------------------------------------$(RESET)

build: ## build grpc server
	go build -o grover

start: build ## start grpc server
	./grover & echo $$! > $(SERVER_PID);


stop: ## stop grpc server
	kill `cat $(SERVER_PID)` && rm -rf $(SERVER_PID)


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