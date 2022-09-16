FROM golang:1.19 as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CGO_ENABLED=0
RUN go build -o dist ./cmd/server

FROM gcr.io/distroless/base-debian11

COPY --from=builder /app/dist .
CMD ["./dist"]
