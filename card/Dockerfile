# syntax=docker/dockerfile:1

FROM golang:1.26 AS deps

WORKDIR /src

COPY platform/go.mod platform/go.sum platform/
COPY card/go.mod card/go.sum card/

WORKDIR /src/card
RUN --mount=type=cache,target=/go/pkg/mod go mod download

FROM deps AS build

COPY platform/ /src/platform/
COPY card/ /src/card/

WORKDIR /src/card
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" \
    -o /out/card ./cmd/grpc_server

FROM golang:1.26 AS migrate

ENV GOOSE_DRIVER=postgres
RUN go install github.com/pressly/goose/v3/cmd/goose@v3.24.3

WORKDIR /app
COPY card/migrations ./migrations

ENTRYPOINT ["goose", "-dir", "migrations", "up"]

FROM gcr.io/distroless/static-debian12:nonroot AS prod

COPY --from=build /out/card /card

USER nonroot:nonroot
EXPOSE 50051

ENTRYPOINT ["/card"]
