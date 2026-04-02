FROM golang:1.26

WORKDIR /api

COPY ./app/go.mod ./app/go.sum ./
RUN go mod download

COPY app/ .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /api/server ./cmd/main.go

EXPOSE 8080
CMD ["./server"]
