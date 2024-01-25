# 使用官方 Golang 镜像作为构建环境
FROM golang:1.21-bullseye as builder

WORKDIR /app

COPY . .

RUN GOOS=linux go build -o palworld 

FROM ubuntu:focal

WORKDIR /app

COPY --from=build /app/palworld /app/palworld

CMD ["/app/palworld"]