# 多阶段构建 - 最小化最终镜像大小
FROM golang:1.24-alpine AS builder

WORKDIR /app

# 复制依赖文件
COPY app-server/go.mod app-server/go.sum ./
RUN go mod download

# 复制源代码
COPY app-server/ ./
COPY shared/ ../shared/

# 构建静态二进制文件
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 最终阶段 - 最小镜像
FROM alpine:latest

# 安装CA证书（HTTPS需要）
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .

# 复制静态文件
COPY --from=builder /shared ./shared

# 暴露端口
EXPOSE 8080

# 设置环境变量
ENV PORT=8080
ENV AWS_REGION=ap-northeast-1
ENV S3_BUCKET_NAME=raymond-go-s3-uploader-dev-2025

# 启动应用
CMD ["./main"]
