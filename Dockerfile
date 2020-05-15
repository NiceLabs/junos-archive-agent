FROM golang:alpine AS builder
LABEL stage=builder
WORKDIR /workspace
COPY . .
RUN CGO_ENABLED=0 go build -o agent -v .

FROM alpine:latest AS release
WORKDIR /
COPY --from=builder /workspace/agent .
EXPOSE 20
EXPOSE 21
CMD ["./agent"]