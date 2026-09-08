FROM golang:1.25-alpine
ENV GO111MODULE=off
ENV GOCACHE=/root/.cache/go-build
WORKDIR /opt/warmup
COPY warmup.go .
RUN go build -o /tmp/warmup warmup.go && rm -rf /opt/warmup
WORKDIR /code
