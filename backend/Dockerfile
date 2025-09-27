FROM golang:1.24.5-alpine AS builder

WORKDIR /app

# install dependencies
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o backend-app ./cmd/main.go

# Stage 2: final
FROM ubuntu

WORKDIR /root

# copy compiled service from previous stage
COPY --from=builder /app/backend-app .

# copy config file
COPY ./configs ./configs

# app start
CMD [ "./backend-app" ]