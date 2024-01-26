FROM golang:1.21-bullseye as builder

WORKDIR /app

COPY . .

ENV GOPROXY=https://goproxy.cn,direct
RUN make

FROM ubuntu:focal

WORKDIR /app

COPY --from=builder /app/build/monitor /app/monitor

COPY endpoint.sh /app/endpoint.sh

RUN chmod +x /app/endpoint.sh

CMD ["/app/endpoint.sh"]
