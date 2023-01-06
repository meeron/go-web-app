FROM golang:1.20rc2-alpine AS source
WORKDIR /source
COPY . .
RUN go test -v ./...
RUN go build .

FROM alpine:3.17 as app
WORKDIR /app
COPY --from=source /source/web-app ./web-app
COPY --from=source /source/config.prod.json ./config.prod.json
ENV GO_APP_ENV=prod
CMD ["/app/web-app"]
