FROM golang:1.20 as builder
ARG CMD_BIN_DIR
ARG PORT
WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o myapp $CMD_BIN_DIR/main.go

FROM alpine:3.14

WORKDIR /app
COPY --from=builder /app/myapp /app/
COPY --from=builder /app/.env /app/

EXPOSE $PORT
CMD ["/app/myapp"]