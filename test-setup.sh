#!/bin/bash

echo "🧪 测试项目设置..."

# 检查Go模块
echo "📦 检查Go模块..."
if [ -f "go.mod" ]; then
    echo "✅ go.mod 文件存在"
else
    echo "❌ go.mod 文件不存在"
    exit 1
fi

# 检查主要Go文件
echo "🔧 检查Go文件..."
if [ -f "main.go" ]; then
    echo "✅ main.go 文件存在"
else
    echo "❌ main.go 文件不存在"
    exit 1
fi

if [ -d "models" ]; then
    echo "✅ models 目录存在"
else
    echo "❌ models 目录不存在"
    exit 1
fi

if [ -d "handlers" ]; then
    echo "✅ handlers 目录存在"
else
    echo "❌ handlers 目录不存在"
    exit 1
fi

# 检查前端文件
echo "🎨 检查前端文件..."
if [ -f "web/package.json" ]; then
    echo "✅ package.json 文件存在"
else
    echo "❌ package.json 文件不存在"
    exit 1
fi

if [ -f "web/nuxt.config.ts" ]; then
    echo "✅ nuxt.config.ts 文件存在"
else
    echo "❌ nuxt.config.ts 文件不存在"
    exit 1
fi

if [ -d "web/pages" ]; then
    echo "✅ pages 目录存在"
else
    echo "❌ pages 目录不存在"
    exit 1
fi

if [ -d "web/components" ]; then
    echo "✅ components 目录存在"
else
    echo "❌ components 目录不存在"
    exit 1
fi

# 检查配置文件
echo "⚙️ 检查配置文件..."
if [ -f "env.example" ]; then
    echo "✅ env.example 文件存在"
else
    echo "❌ env.example 文件不存在"
    exit 1
fi

if [ -f ".gitignore" ]; then
    echo "✅ .gitignore 文件存在"
else
    echo "❌ .gitignore 文件不存在"
    exit 1
fi

if [ -f "README.md" ]; then
    echo "✅ README.md 文件存在"
else
    echo "❌ README.md 文件不存在"
    exit 1
fi

# 检查Docker文件
echo "🐳 检查Docker文件..."
if [ -f "Dockerfile" ]; then
    echo "✅ Dockerfile 文件存在"
else
    echo "❌ Dockerfile 文件不存在"
    exit 1
fi

if [ -f "docker-compose.yml" ]; then
    echo "✅ docker-compose.yml 文件存在"
else
    echo "❌ docker-compose.yml 文件不存在"
    exit 1
fi

# 检查uploads目录
echo "📁 检查uploads目录..."
if [ -d "uploads" ]; then
    echo "✅ uploads 目录存在"
else
    echo "❌ uploads 目录不存在"
    exit 1
fi

echo ""
echo "🎉 所有检查通过！项目设置正确。"
echo ""
echo "📋 下一步："
echo "1. 复制 env.example 为 .env 并配置数据库"
echo "2. 运行 ./start.sh 启动项目"
echo "3. 或者使用 docker-compose up 启动Docker版本" 