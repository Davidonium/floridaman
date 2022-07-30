FROM golang:1.18 as builder
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

COPY --from=builder /tmp/go/build/floridaman /opt/floridaman/floridaman
COPY --from=builder /etc/passwd /etc/passwd
USER floridaman

CMD [ "/opt/floridaman/floridaman", "serve" ]
