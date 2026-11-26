# pokemon-generator: one Go binary serving the page, the API and sprite images, with the page's
# files embedded.
FROM artifacts.home.romaine.life/dockerhub/library/golang:1.26-alpine@sha256:30b1e546a9c7fa3b747c2bdba3a0c7bdf3a02ebb4c5b8405d2f948d7e8aac4ec AS build
WORKDIR /src
ENV GOPROXY=https://artifacts.home.romaine.life/go/go,direct
ENV GOMAXPROCS=4 GOFLAGS=-p=4
# The release this build is (the web-app suite passes the tag), shown in the page's footer and
# put in the page's addresses so browsers never keep an older release's answers.
ARG APP_VERSION=dev

COPY go.mod go.sum ./
RUN --mount=type=cache,id=pokemon-generator-go-mod,target=/go/pkg/mod go mod download
COPY cmd ./cmd
COPY gen ./gen
COPY internal ./internal
COPY web ./web
RUN --mount=type=cache,id=pokemon-generator-go-mod,target=/go/pkg/mod \
    --mount=type=cache,id=pokemon-generator-go-build,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -X main.version=${APP_VERSION}" -o /out/pokemon-generator ./cmd/pokemon-generator

FROM artifacts.home.romaine.life/gcr/distroless/static-debian12:nonroot@sha256:afa5c872c891853ca7fcf1f12c3edb23f7eeef36189728842dd51042ff57f7ab
COPY --from=build /out/pokemon-generator /pokemon-generator
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/pokemon-generator"]
