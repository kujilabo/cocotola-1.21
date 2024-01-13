FROM golang:1.21 as builder

WORKDIR /go/src
COPY . .

RUN CGO_ENABLED=0 go build src/main.go 

FROM gcr.io/distroless/static-debian12

COPY --from=builder /go/src/main /
CMD ["/main"]
