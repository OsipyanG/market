## Build
FROM golang:1.22.5-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/app/main.go

## Deploy
FROM scratch as final

COPY --from=build /app/main .
COPY --from=build /app/.env .

ENTRYPOINT ["./main"]
