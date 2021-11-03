FROM golang:alpine3.13 as builder

RUN apk update && apk upgrade && apk add build-base git make sed
RUN go get github.com/silenceper/gowatch

WORKDIR /go/src/github.com/eduardohoraciosanto/blockchain-experiment
COPY . .

RUN GIT_COMMIT=$(git rev-parse --short HEAD) && \
  go build -o service -ldflags "-X 'github.com/eduardohoraciosanto/blockchain-experiment/config.serviceVersion=$GIT_COMMIT'"

FROM alpine:3.13 

COPY --from=builder /go/src/github.com/eduardohoraciosanto/blockchain-experiment/service /

ENTRYPOINT [ "./service" ]