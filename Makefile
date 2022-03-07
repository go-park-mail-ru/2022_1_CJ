# path to actual config - the one that is copied to the docker container
CONFIG:=resources/config/config.yaml

# path to docker compose file
DCOMPOSE:=docker-compose.yaml

# path to external config which will copied to CONFIG
CONFIG_PATH=resources/config/config_default.yaml

# improve build time
DOCKER_BUILD_KIT:=COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1

all: down build up

debug: down build-debug up-debug

down:
	docker-compose -f ${DCOMPOSE} down

build:
	cp ${CONFIG_PATH} ${CONFIG}
	${DOCKER_BUILD_KIT} docker-compose build --no-cache --pull --build-arg CONFIG=${CONFIG}
	
build-debug:
	cp ${CONFIG_PATH} ${CONFIG}
	${DOCKER_BUILD_KIT} docker-compose build --build-arg CONFIG=${CONFIG}

up:
	docker-compose -f ${DCOMPOSE} up -d

up-debug:
	docker-compose -f ${DCOMPOSE} up

mod:
	go mod tidy -compat=1.17 && go get ./...
