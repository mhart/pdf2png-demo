FROM golang:1-alpine3.12 AS build-image

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
RUN go build -tags lambda.norpc -ldflags="-s -w" main.go


FROM alpine:3.12

WORKDIR /app

# pdf2cairo is part of the poppler-utils package
# ttf-liberation includes common fonts we might need
RUN apk add --no-cache poppler-utils ttf-liberation

COPY --from=build-image /app/main ./

ENTRYPOINT ["/app/main"]
