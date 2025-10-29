#!/bin/bash

echo "========================================="
echo "  启动彩票管理系统后端服务"
echo "========================================="
echo ""

cd backend

# 检查是否安装了Go
if ! command -v go &> /dev/null; then
    echo "❌ 错误: 未检测到Go环境，请先安装Go 1.25或更高版本"
    exit 1
fi

echo "✅ Go版本:"
go version
echo ""

# 检查go.mod是否存在
if [ ! -f "go.mod" ]; then
    echo "❌ 错误: go.mod文件不存在"
    exit 1
fi

# 下载依赖
echo "📦 正在下载依赖..."
go mod download
go mod tidy

echo ""
echo "🚀 启动后端服务..."
echo "   服务地址: http://localhost:8080"
echo "   API文档: http://localhost:8080/api"
echo ""

# 运行服务
go run main.go

