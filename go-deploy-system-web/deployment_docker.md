# 容器化发布项目-前端

1.Dockerfile文件内容  
``` Dockerfile
# 编译构建
FROM harbor.od.com/public/node:15.14.0 as builder

WORKDIR /app

COPY . /app/


# RUN npm i -g @vue/cli-init && \
#     npm i -g @vue/cli && \
#     npm install --save-dev webpack && \
#     npm install --save axios@^0.27.2  --save && \
#     npm install --save element-ui &&
#     npm run build

# 把整个项目包括依赖都上传，只要执行编译即可
RUN npm run build

# 运行容器
FROM harbor.od.com/public/nginx:v1.7.9

COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]

#docker build -t harbor.od.com/app/go-deploy-system-web:v0.1 .
#docker run --name go-deploy-web -d -p8889:80 harbor.od.com/app/go-deploy-system-web:v0.1
#docker push harbor.od.com/app/go-deploy-system-web:v0.1

```

2.yaml文件  
#dp.yaml  
```yaml
kind: Deployment
apiVersion: extensions/v1beta1
metadata:
  name: go-deploy-system-web
  namespace: test
  labels:
    name: go-deploy-system-web
spec:
  replicas: 1
  selector:
    matchLabels:
      name: go-deploy-system-web
  template:
    metadata:
      labels:
        app: go-deploy-system-web
        name: go-deploy-system-web
    spec:
      containers:
      - name: go-deploy-system-web
        image: harbor.od.com/app/go-deploy-system-web:v0.1
        ports:
        - containerPort: 80
          protocol: TCP
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
  name: go-deploy-system-web
  namespace: test
spec:
  ports:
  - protocol: TCP
    port: 80
    targetPort: 80
  selector:
    app: go-deploy-system-web
```

#ingress.yaml  
```yaml
kind: Ingress
apiVersion: extensions/v1beta1
metadata:
  name: go-deploy-system-web
  namespace: test
spec:
  rules:
  - host: go-deploy-system-web-test.od.com
    http:
      paths:
      - path: /
        backend:
          serviceName: go-deploy-system-web
          servicePort: 80

```

3.域名解析  
```shell
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
go-deploy-system-web-test   A    10.4.7.10

#重启域名服务
#systemctl restart named
#验证
#dig  -t A go-deploy-system-web-test.od.com @10.4.7.11 +short
10.4.7.10
```

4.k8s应用  
```shell
# kubectl apply -f dp.yaml
# kubectl apply -f svc.yaml
# kubectl apply -f ingress.yaml

```