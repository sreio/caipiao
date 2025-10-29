#!/bin/bash

echo "========================================="
echo "  启动彩票管理系统前端服务"
echo "========================================="
echo ""

cd frontend

# 检查是否安装了Node.js
if ! command -v node &> /dev/null; then
    echo "❌ 错误: 未检测到Node.js环境，请先安装Node.js"
    exit 1
fi

echo "✅ Node.js版本:"
node -v
echo ""

echo "✅ npm版本:"
npm -v
echo ""

# 检查package.json是否存在
if [ ! -f "package.json" ]; then
    echo "❌ 错误: package.json文件不存在"
    exit 1
fi

# 检查node_modules是否存在，不存在则安装依赖
if [ ! -d "node_modules" ]; then
    echo "📦 正在安装依赖..."
    npm install
else
    echo "✅ 依赖已安装"
fi

echo ""
echo "🚀 启动前端开发服务器..."
echo "   访问地址: http://localhost:5173"
echo ""

# 运行开发服务器
npm run dev

