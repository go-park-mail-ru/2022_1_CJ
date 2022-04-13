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
	docker-compose -f ${DCOMPOSE} down --remove-orphans

build:
	cp ${CONFIG_PATH} ${CONFIG}
	${DOCKER_BUILD_KIT} docker-compose build --no-cache --pull --build-arg CONFIG=${CONFIG}
	
build-debug:
	cp ${CONFIG_PATH} ${CONFIG}
	${DOCKER_BUILD_KIT} docker-compose build --build-arg CONFIG=${CONFIG}

up:
	docker-compose --compatibility -f ${DCOMPOSE} up -d --remove-orphans

up-debug:
	docker-compose --compatibility -f ${DCOMPOSE} up --remove-orphans

# Vendoring is useful for local debugging since you don't have to
# reinstall all packages again and again in docker
mod:
	go mod tidy -compat=1.18 && go install ./... && go mod vendor

tests:
	go test ./test/db/... -cover -coverpkg ./internal/db/...
	go test ./test/services/... -cover -coverpkg ./internal/service/...

mock:
	mockgen -source=internal/db/friends.go -destination=mocks/friends_db_mock.go \
	&& mockgen -source=internal/db/post.go -destination=mocks/post_db_mock.go \
	&& mockgen -source=internal/db/user.go -destination=mocks/user_db_mock.go \
	&& mockgen -source=internal/db/chat.go -destination=mocks/chat_db_mock.go
