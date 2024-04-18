#target: ghcr.io/kispaljr/update-configmap:v1

FROM golang:1.21-alpine
ENV CGO_ENABLED=0
WORKDIR /go/src/
COPY go.* ./
RUN go mod download
COPY ./ ./
RUN go build -o /usr/local/bin/function ./

FROM gcr.io/distroless/static:latest
COPY --from=0 /usr/local/bin/function /usr/local/bin/function
ENTRYPOINT ["function"]
