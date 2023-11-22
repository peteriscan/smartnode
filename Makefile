VERSION := $(shell grep RocketPoolVersion shared/version.go | grep -Eo '".+?"' | sed 's/"//g')

all:
	docker run -it --rm -v $(shell pwd):/smartnode -v $(shell pwd)/tmp/go:/go -w /smartnode/rocketpool-cli rocketpool/smartnode-builder:latest go build -buildvcs=false
	docker run -it --rm -v $(shell pwd):/smartnode -v $(shell pwd)/tmp/go:/go rocketpool/smartnode-builder:latest /smartnode/rocketpool/build.sh
	docker build -f docker/rocketpool-dockerfile -t rocketpool/smartnode:v$(VERSION) .
