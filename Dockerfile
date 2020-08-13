FROM node:12.18-alpine AS client_builder

WORKDIR /app

COPY ./client .

RUN yarn

RUN yarn build

FROM golang:1.15-alpine as go_builder

LABEL maintainer="Sundowndev" \
  org.label-schema.name="phoneinfoga" \
  org.label-schema.description="Advanced information gathering & OSINT tool for phone numbers." \
  org.label-schema.url="https://github.com/sundowndev/PhoneInfoga" \
  org.label-schema.vcs-url="https://github.com/sundowndev/PhoneInfoga" \
  org.label-schema.vendor="Sundowndev" \
  org.label-schema.schema-version="1.0"

WORKDIR /app

COPY . .

RUN go get -v -t -d ./...

COPY --from=client_builder /app/dist ./client/dist

RUN go get -v github.com/jessevdk/go-assets-builder

RUN go generate ./...

RUN go build -v -o phoneinfoga .

FROM golang:1.15-alpine

WORKDIR /app

COPY --from=go_builder /app/phoneinfoga .

EXPOSE 5000

ENTRYPOINT ["/app/phoneinfoga"]
