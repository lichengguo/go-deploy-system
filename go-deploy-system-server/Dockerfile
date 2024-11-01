# 编译构建
FROM harbor.od.com/public/golang:1.20.14 as builder

WORKDIR /app

ENV GOPROXY https://goproxy.cn

COPY . /app/

RUN  go mod tidy && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o go-deploy-system-server-linux .

# 运行
FROM harbor.od.com/public/alpine:3.18

# 该项目依赖ssh和git
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk update && \
    apk --no-cache add sshpass rsync openssh git tzdata ca-certificates && \
    cp -f /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    rm -rf /var/cache/apk/* && \
    mkdir -p /root/.ssh && \
    echo -e "StrictHostKeyChecking no\nUserKnownHostsFile /dev/null" >/root/.ssh/config

COPY --from=builder /app/go-deploy-system-server-linux /app/go-deploy-system-server-linux

COPY --from=builder /app/config/config.ini /app/config/config.ini

EXPOSE 3000

CMD ["/app/go-deploy-system-server-linux"]

#docker build -t harbor.od.com/app/go-deploy-system-server:v0.1 .
#docker run --name go-deploy-server -d -p3000:3000 harbor.od.com/app/go-deploy-system-server:v0.1
#docker push harbor.od.com/app/go-deploy-system-server:v0.1