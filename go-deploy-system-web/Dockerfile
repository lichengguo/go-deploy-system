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