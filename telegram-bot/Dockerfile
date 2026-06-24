# syntax=docker/dockerfile:1

FROM golang:1.26 AS deps

WORKDIR /src

COPY platform/go.mod platform/go.sum platform/
COPY card/go.mod card/go.sum card/
COPY telegram-bot/go.mod telegram-bot/go.sum telegram-bot/

WORKDIR /src/telegram-bot
RUN --mount=type=cache,target=/go/pkg/mod go mod download

FROM deps AS build

COPY platform/ /src/platform/
COPY card/ /src/card/
COPY telegram-bot/ /src/telegram-bot/

WORKDIR /src/telegram-bot
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" \
    -o /out/bot ./cmd/bot && \
    CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" \
    -o /out/miniapp ./cmd/miniapp

FROM golang:1.26 AS migrate

ENV GOOSE_DRIVER=postgres
RUN go install github.com/pressly/goose/v3/cmd/goose@v3.24.3

WORKDIR /app
COPY telegram-bot/migrations ./migrations

ENTRYPOINT ["goose", "-dir", "migrations", "up"]

FROM gcr.io/distroless/static-debian12:nonroot AS prod-bot

COPY --from=build /out/bot /bot

USER nonroot:nonroot

ENTRYPOINT ["/bot"]

FROM gcr.io/distroless/static-debian12:nonroot AS prod-miniapp

COPY --from=build /out/miniapp /miniapp

USER nonroot:nonroot
EXPOSE 8080

ENTRYPOINT ["/miniapp"]
