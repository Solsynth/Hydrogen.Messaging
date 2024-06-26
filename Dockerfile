# Building Backend
FROM golang:alpine as messaging-server

WORKDIR /source
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -buildvcs -o /dist ./pkg/main.go

# Runtime
FROM golang:alpine

COPY --from=messaging-server /dist /messaging/server

EXPOSE 8447

CMD ["/messaging/server"]
