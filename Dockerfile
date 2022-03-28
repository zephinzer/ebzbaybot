ARG GO_VERSION=1.18
FROM golang:${GO_VERSION}-alpine AS build
RUN apk add --no-cache make g++ ca-certificates
WORKDIR /go/src/app
COPY ./go.mod ./go.sum ./
RUN go mod download -x
COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./pkg ./pkg
COPY ./Makefile ./
COPY ./*.go ./
# this should be set by the build recipe in the Makefile
ARG RELEASE_TAG=latest
ENV RELEASE_TAG=${RELEASE_TAG}
RUN make release_tag=${RELEASE_TAG} build
RUN mv ./bin/ebzbaybot_$(go env GOOS)_$(go env GOARCH) ./bin/ebzbaybot
RUN mv ./bin/ebzbaybot_$(go env GOOS)_$(go env GOARCH).sha256 ./bin/ebzbaybot.sha256

FROM scratch AS final
COPY --from=build /go/src/app/bin/ebzbaybot /ebzbaybot
COPY --from=build /go/src/app/bin/ebzbaybot.sha256 /ebzbaybot.sha256
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT ["/ebzbaybot"]
LABEL repo_url https://github.com/zephinzer/ebzbaybot
LABEL maintainer ebzbaybot
