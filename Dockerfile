FROM golang:1.19 AS build
WORKDIR /go/src

COPY go.mod .
COPY go.sum .

ENV CGO_ENABLED=0
RUN go mod download

COPY main.go .
COPY hub.go .
COPY message.go .
COPY client.go .

RUN go build -a -installsuffix cgo -o service .

FROM scratch AS runtime
ENV GIN_MODE=release
EXPOSE 3000/tcp

COPY --from=build /go/src/service ./

ENTRYPOINT ["./service"]
