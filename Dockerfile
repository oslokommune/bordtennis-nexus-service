FROM golang:1.19 AS build
ENV CGO_ENABLED=0

WORKDIR /go/src

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY main.go .
COPY pkg/ pkg/

RUN go build -a -installsuffix cgo -o service .

FROM scratch AS runtime
ENV GIN_MODE=release

EXPOSE 3000/tcp

COPY --from=build /go/src/service ./

ENTRYPOINT ["./service"]
