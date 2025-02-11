FROM golang:1.23 as builder

WORKDIR /go/src/app

COPY . .

RUN go mod tidy && \
    CGO_ENABLED=0 go build -o /ya-boo

FROM gcr.io/distroless/static-debian12

COPY --from=builder /ya-boo /

CMD ["/ya-boo", "run"]
