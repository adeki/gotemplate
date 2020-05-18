FROM golang:1.14.3 AS builder

WORKDIR /go/src/migrate
RUN CGO_ENABLED=0 go get -u github.com/rubenv/sql-migrate/...


FROM alpine:3.11
COPY --from=builder /go/bin/sql-migrate /usr/local/bin

ENTRYPOINT ["sql-migrate"]
CMD ["-h"]
