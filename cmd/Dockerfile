FROM golang:1.22-alpine AS pre-builder

WORKDIR /build

RUN apk update && apk upgrade && apk add --no-cache ca-certificates upx
RUN update-ca-certificates

COPY .. ./

FROM pre-builder AS prod-builder

ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" \
    CGO_ENABLE=0 \
    GOOS=linux \
    go build -a -ldflags '-s -w' -mod vendor -o /app/cmd ./cmd

RUN upx /app/cmd

FROM scratch AS prod

WORKDIR /app
COPY --from=prod-builder /app/cmd ./
COPY --from=prod-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8000

ENV APP__CONFIG__PATH="/app/config.json"

CMD ["/app/cmd"]

FROM pre-builder AS dev-builder

ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" \
    GOOS=linux \
    go build -mod vendor -o /app/cmd ./cmd

FROM scratch AS dev

WORKDIR /app
COPY --from=dev-builder /app/cmd ./
COPY --from=dev-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8000

ENV APP__CONFIG__PATH="/app/config.json"

CMD ["/app/cmd"]