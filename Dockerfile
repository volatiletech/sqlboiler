# builder
FROM golang:1.13-alpine AS builder

RUN apk add git build-base

WORKDIR /go/src/github.com/volatiletech/sqlboiler

COPY . .
RUN go get -v -t ./...
RUN go build -trimpath -ldflags "-w -s" . && \
    go build -trimpath -ldflags "-w -s" ./drivers/sqlboiler-psql && \
    go build -trimpath -ldflags "-w -s" ./drivers/sqlboiler-mysql && \
    go build -trimpath -ldflags "-w -s" ./drivers/sqlboiler-mssql


# sqlboiler (no drivers, just to take advantage of layer caching)
FROM alpine:3.10 AS sqlboiler

WORKDIR /sqlboiler

COPY --from=builder /go/src/github.com/volatiletech/sqlboiler/sqlboiler \
                    /usr/local/bin/

ENTRYPOINT [ "sqlboiler" ]

# sqlboiler-mssql
FROM sqlboiler AS sqlboiler-mssql

COPY --from=builder /go/src/github.com/volatiletech/sqlboiler/sqlboiler-mssql \
                    /usr/local/bin/

# sqlboiler-mysql
FROM sqlboiler AS sqlboiler-mysql

COPY --from=builder /go/src/github.com/volatiletech/sqlboiler/sqlboiler-mysql \
                    /usr/local/bin/

# sqlboiler-psql
FROM sqlboiler AS sqlboiler-psql

COPY --from=builder /go/src/github.com/volatiletech/sqlboiler/sqlboiler-psql \
                    /usr/local/bin/

# sqlboiler (all drivers included)
FROM sqlboiler

COPY --from=builder /go/src/github.com/volatiletech/sqlboiler/sqlboiler-mssql \
                    /go/src/github.com/volatiletech/sqlboiler/sqlboiler-mysql \
                    /go/src/github.com/volatiletech/sqlboiler/sqlboiler-psql \
                    /usr/local/bin/
