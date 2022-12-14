all: build

build: build-binaries

# Build the binaries for the service.
build-binaries:
	$(GOBUILD) -ldflags $(LDFLAGS) \
	-o $(BIN)/$(TARGET) \
	service/main.go

# Install all files in local folder
install: build
	cp -r config/config.yaml $(BIN)/

local_test_infra:
	make -C tools/compose

# Create a docker image for the service.
docker-image:
	docker build -t $(DOCKER_IMAGE) -f docker/Dockerfile .

# Start the service and its dependencies within a local Docker network.
# put perhaps a script that can start your service here. Recommended using docker compose
start:

# Publish the fs docker image to Github.
publish: docker-image
	docker tag $(DOCKER_IMAGE):latest $(GHCR_IMAGE):latest
	docker push $(GHCR_IMAGE):latest

include build/Makefile.*
