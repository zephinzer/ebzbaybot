app := ebzbaybot
image_registry := docker.io
image_owner := zephinzer
image_path := ebzbaybot

-include ./Makefile.properties

image_url ?= $(image_registry)/$(image_owner)/$(image_path)
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
			-X github.com/zephinzer/ebzbaybot/internal/constants.Version=\"$(release_tag)\" \
			-extldflags 'static' -s -w \
		" \
		-o ./bin/$(app)_$$(go env GOOS)_$$(go env GOARCH)$(bin_ext) \
		.
	cd ./bin && sha256sum $(app)_$$(go env GOOS)_$$(go env GOARCH)$(bin_ext) > $(app)_$$(go env GOOS)_$$(go env GOARCH)$(bin_ext).sha256
image:
	docker build \
		--network host \
		--file ./Dockerfile \
		--build-arg release_tag=$(release_tag) \
		--tag $(image_url):latest \
		.
release: image
	docker tag $(image_url):latest $(image_url):$(release_tag)
	docker push $(image_url):$(release_tag)
install:
	go install .
