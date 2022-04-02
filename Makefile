app := ebzbaybot
image_registry := docker.io
image_owner := zephinzer
image_path := ebzbaybot

-include ./Makefile.properties

image_url ?= $(image_registry)/$(image_owner)/$(image_path)
release_version ?= v0.0.1
release_tag ?= latest

# this adjusts the binary extension for when building with windows
bin_ext := 
ifeq "${GOOS}" "windows"
bin_ext := .exe
endif

env:
	@docker-compose up -d
deps:
	@go mod vendor
start:
	@go run . start $(app)
test:
	mkdir -p ./coverage
	go test -v -mod=mod -cover -covermode=atomic -coverpkg=./... -coverprofile=./coverage/golang.all.out ./...
	grep -v ".pb.go" ./coverage/golang.all.out > ./coverage/golang.out
	go tool cover -func ./coverage/golang.out
lint:
	go vet ./...
build:
	CGO_ENABLED=0 go build \
		-ldflags " \
			-X github.com/zephinzer/ebzbaybot/internal/constants.Version=\$(release_version)-$(release_tag)\" \
			-extldflags 'static' -s -w \
		" \
		-o ./bin/$(app)_$$(go env GOOS)_$$(go env GOARCH)$(bin_ext) \
		.
	cd ./bin && sha256sum $(app)_$$(go env GOOS)_$$(go env GOARCH)$(bin_ext) > $(app)_$$(go env GOOS)_$$(go env GOARCH)$(bin_ext).sha256
image:
	docker build \
		--network host \
		--file ./Dockerfile \
		--build-arg RELEASE_TAG=$(release_tag) \
		--build-arg RELEASE_VERSION=$(release_version) \
		--tag $(image_url):latest \
		.
release: image
	docker tag $(image_url):latest $(image_url):$(release_tag)
	docker push $(image_url):latest
	docker push $(image_url):$(release_tag)
ifdef release_version
	docker tag $(image_url):latest $(image_url):$(release_version)
	docker push $(image_url):$(release_version)
endif
release-git:
	@$(MAKE) release release_tag=$$(git rev-parse HEAD | head -c 8) release_version=$$(git describe --tags)
install:
	go install .
