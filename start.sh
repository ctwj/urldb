#!/bin/bash

echo "🚀 启动资源管理系统..."

# 检查Go是否安装
if ! command -v go &> /dev/null; then
    echo "❌ Go未安装，请先安装Go"
    exit 1
fi

# 检查Node.js是否安装
if ! command -v node &> /dev/null; then
    echo "❌ Node.js未安装，请先安装Node.js"
    exit 1
fi

# 检查PostgreSQL是否运行
if ! pg_isready -q; then
    echo "⚠️  PostgreSQL未运行，请确保PostgreSQL服务已启动"
fi

echo "📦 安装Go依赖..."
go mod tidy

echo "🌐 启动后端服务器..."
go run main.go &
BACKEND_PID=$!

echo "⏳ 等待后端启动..."
sleep 3

echo "📦 安装前端依赖..."
cd web
npm install

echo "🎨 启动前端开发服务器..."
npm run dev &
FRONTEND_PID=$!

echo "✅ 系统启动完成！"
echo "📱 前端地址: http://localhost:3000"
echo "🔧 后端地址: http://localhost:8080"
echo ""
echo "按 Ctrl+C 停止服务"

# 等待用户中断
trap "echo '🛑 正在停止服务...'; kill $BACKEND_PID $FRONTEND_PID; exit" INT
wait 