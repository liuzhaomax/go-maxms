FROM golang:1.22 as builder

ENV GO111MODULE on
ENV CGO_ENABLED 1
ENV GOOS linux
ENV GOARCH amd64

WORKDIR /workspace

COPY . .

RUN go mod tidy

RUN go build -o bin/main main/main.go

FROM ubuntu:latest

ENV GO111MODULE on
ENV CGO_ENABLED 1
ENV GOOS linux
ENV GOARCH amd64

WORKDIR /workspace

COPY --from=builder /workspace/bin/main /workspace/bin/main
COPY --from=builder /workspace/script /workspace/script
COPY --from=builder /workspace/environment /workspace/environment

CMD ["./bin/main"]
