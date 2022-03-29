FROM node:16 AS client_builder

WORKDIR /app

COPY ./client .

RUN yarn install --immutable && \
    yarn build && \
    yarn cache clean

FROM golang:1.18-alpine AS go_builder

LABEL maintainer="Sundowndev" \
  org.label-schema.name="phoneinfoga" \
  org.label-schema.description="Advanced information gathering & OSINT tool for phone numbers." \
  org.label-schema.url="https://github.com/sundowndev/phoneinfoga" \
  org.label-schema.vcs-url="https://github.com/sundowndev/phoneinfoga" \
  org.label-schema.vendor="Sundowndev" \
  org.label-schema.schema-version="1.0"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download 
COPY . .

RUN apk add --no-cache git && \
    go get -v -t -d ./...

COPY --from=client_builder /app/dist ./client/dist

RUN go generate ./... && \
    go build -v -ldflags="-s -w \
    -X 'github.com/sundowndev/phoneinfoga/v2/config.Version=$(git describe --abbrev=0 --tags)' \
    -X 'github.com/sundowndev/phoneinfoga/v2/config.Commit=$(git rev-parse --short HEAD)'" -v -o phoneinfoga .

FROM alpine:3.15

WORKDIR /app

COPY --from=go_builder /app/phoneinfoga .

EXPOSE 5000

ENTRYPOINT ["/app/phoneinfoga"]
