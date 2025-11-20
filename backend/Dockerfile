# 后端Dockerfile
FROM golang:1.25-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的依赖
RUN apk add --no-cache gcc musl-dev sqlite-dev

# 复制go mod文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译应用
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o lottery-server .

# 运行阶段
FROM alpine:latest

# 安装运行时依赖
RUN apk --no-cache add ca-certificates sqlite

# 设置工作目录
WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/lottery-server .

# 创建数据目录
RUN mkdir -p data

# 暴露端口
EXPOSE 8080

# 启动应用
CMD ["./lottery-server"]
