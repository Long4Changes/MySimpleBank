# Multi-Stage Build

# Build Stage
# 指定基础镜像
FROM golang:1.26.3-alpine3.23 AS builder

# 设置工作目录
WORKDIR /app

# 将本地代码复制到镜像里
COPY . .

# 构建一个二进制文件
RUN go build -o main main.go


# Run Stage
FROM alpine:3.23

WORKDIR /app

# 从 Build Stage 复制二进制文件到 Run Stage
COPY --from=builder /app/main .

# 声明程序运行的端口
EXPOSE 8080

# 设置容器启动时执行的命令
CMD [ "/app/main" ]