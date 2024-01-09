# 使用官方的Go镜像作为基础镜像
FROM golang:1.21.4

# 设置工作目录
WORKDIR /app

# 将本地的main.go文件复制到容器中的/app目录
COPY . /app/

# 使用go mod下载依赖
# RUN go mod tidy
RUN ls
# 设置环境变量
ENV GOARCH=amd64
ENV GOOS=linux

# 编译Go应用程序
RUN go build main.go




# 启动应用程序
CMD ["/app/main"]