# we assume all test machines have either docker-compose or
# the compose plugin in docker command.
DOCKER_COMPOSE=docker-compose
HAS_DOCKER_COMPOSE := $(shell command -v docker-compose 2> /dev/null)
ifndef HAS_DOCKER_COMPOSE
  DOCKER_COMPOSE=docker compose
  DOCKER_COMPOSE_QUIET_PULL=--quiet-pull
endif

include .env

start:
	$(DOCKER_COMPOSE) -p$(PROJECT) up $(DOCKER_COMPOSE_QUIET_PULL) -d

stop:
	-$(DOCKER_COMPOSE) -p$(PROJECT) down

.PHONY: start stop
