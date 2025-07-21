#!/bin/bash

echo "🚀 启动网盘资源管理系统..."

# 检查Docker是否运行
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker未运行，请先启动Docker"
    exit 1
fi

# 检测Docker Compose命令
if command -v docker-compose &> /dev/null; then
    DOCKER_COMPOSE="docker-compose"
elif docker compose version &> /dev/null; then
    DOCKER_COMPOSE="docker compose"
else
    echo "❌ 未找到Docker Compose，请安装Docker Compose"
    exit 1
fi

echo "📦 使用Docker Compose命令: $DOCKER_COMPOSE"

# 停止并删除现有容器
echo "🔄 清理现有容器..."
$DOCKER_COMPOSE down

# 构建并启动服务
echo "🔨 构建并启动服务..."
$DOCKER_COMPOSE up --build -d

# 等待服务启动
echo "⏳ 等待服务启动..."
sleep 10

# 检查服务状态
echo "📊 服务状态："
$DOCKER_COMPOSE ps

echo ""
echo "✅ 系统启动完成！"
echo "🌐 前端访问地址: http://localhost:3030"
echo "🔧 后端API地址: http://localhost:8080"
echo "🗄️  数据库地址: localhost:5432"
echo ""
echo "📝 查看日志: $DOCKER_COMPOSE logs -f"
echo "🛑 停止服务: $DOCKER_COMPOSE down" 