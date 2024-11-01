## 容器化发布项目

1.Dockerfile文件内容  
```dockerfile
# 编译构建
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
```

2.yaml文件  
#cm.yaml  
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: go-deploy-system-server-cm
  namespace: test
data:
  config.ini: |
    [server]
    # debug 开发模式; release 生产模式
    AppMode = release
    # 监听地址和端口
    HttpPort = :3000
    # JWT加盐字符串
    JwtKey = 89js82js72@a=KCAFJWQER012
    # 登录密码加盐字符串
    PwdKey = aoefqCINAETCA
    # 服务器、Git密码模式下的key
    ServerGitKey = a&D*71&FBA12-9P*
    # 秘钥文件存储目录
    KeyFilePath = data/go_deployment_system/upload/key
    # 代码存放目录
    CodePath = data/go_deployment_system/git
    
    [database]
    # 数据库类型需要为MySQL
    # 数据库连接地址
    DbHost = 10.4.7.11
    # 数据库连接端口
    DbPort = 3306
    # 数据库连接账号
    DbUser = root
    # 数据库连接密码
    DbPassWord = 123456
    # 数据库名称
    DbName = go_deployment_system
    
    [log]
    # 日志文件存储目录
    LogPath = data/go_deployment_system/log
    # 日志文件名称
    LogFileName = ops.log
    # 日志最大保存时间 单位:天
    LogSaveTime = 10
    # 日志切割大小 单位:MB
    LogSplitSize = 10


```

#dp.yaml  
```yaml
kind: Deployment
apiVersion: extensions/v1beta1
metadata:
  name: go-deploy-system-server
  namespace: test
  labels:
    name: go-deploy-system-server
spec:
  replicas: 1
  selector:
    matchLabels:
      name: go-deploy-system-server
  template:
    metadata:
      labels:
        app: go-deploy-system-server
        name: go-deploy-system-server
    spec:
      volumes:
      - name: configmap-volume
        configMap:
          name: go-deploy-system-server-cm
      containers:
      - name: go-deploy-system-server
        image: harbor.od.com/app/go-deploy-system-server:v0.1
        ports:
        - containerPort: 3000
          protocol: TCP
        volumeMounts:
        - name: configmap-volume
          mountPath: app/config
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
        imagePullPolicy: IfNotPresent
      imagePullSecrets:
      - name: harbor
      restartPolicy: Always
      terminationGracePeriodSeconds: 30
      securityContext:
        runAsUser: 0
      schedulerName: default-scheduler
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 1
      maxSurge: 1
  revisionHistoryLimit: 7
  progressDeadlineSeconds: 600

```

#svc.yaml  
```yaml
kind: Service
apiVersion: v1
metadata:
  name: go-deploy-system-server
  namespace: test
spec:
  ports:
  - protocol: TCP
    port: 8888
    targetPort: 3000
  selector:
    app: go-deploy-system-server
```
   
#ingress.yaml(注意是为了测试，否则不需要暴露ingress，前端可以直接使用服务名来访问)  
```yaml
kind: Ingress
apiVersion: extensions/v1beta1
metadata:
  name: go-deploy-system-server
  namespace: test
spec:
  rules:
  - host: go-deploy-server-test.od.com
    http:
      paths:
      - path: /
        backend:
          serviceName: go-deploy-system-server
          servicePort: 8888


```

3.域名解析  
```text
#增加域名解析
#serial需要前滚一位
#vi /var/named/od.com.zone
$ORIGIN od.com.
$TTL 600        ; 10 minutes
@               IN SOA  dns.od.com. dnsadmin.od.com. (
                                2024070518 ; serial
                                10800      ; refresh (3 hours)
                                900        ; retry (15 minutes)
                                604800     ; expire (1 week)
                                86400      ; minimum (1 day)
                                )
                                NS   dns.od.com.
$TTL 60 ; 1 minute
dns                A    10.4.7.11
go-deploy-server-test   A    10.4.7.10

#重启域名服务
#systemctl restart named
#验证
#dig  -t A go-deploy-server-test.od.com @10.4.7.11 +short
10.4.7.10
```

4.k8s上应用
```shell
# kubectl apply -f cm.yaml
# kubectl apply -f dp.yaml
# kubectl apply -f svc.yaml

#注意是为了测试，否则不需要暴露ingress
# kubectl apply -f ingress.yaml
```
