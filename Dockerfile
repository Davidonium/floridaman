FROM golang:1.17 as builder
RUN useradd -u 10001 floridaman

WORKDIR /tmp/go

COPY Makefile ./
COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY cmd/ ./cmd
COPY *.go ./

RUN make build

FROM scratch
USER floridaman
COPY --from=builder /tmp/go/build/floridaman /usr/local/bin/floridaman
COPY --from=builder /etc/passwd /etc/passwd

CMD [ "/usr/local/bin/floridaman", "serve" ]
