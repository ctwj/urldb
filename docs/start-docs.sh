#!/bin/bash

# 启动 docsify 文档服务脚本

echo "🚀 启动 docsify 文档服务..."

# 检查是否安装了 docsify-cli
if ! command -v docsify &> /dev/null; then
    echo "❌ 未检测到 docsify-cli，正在安装..."
    npm install -g docsify-cli
    if [ $? -ne 0 ]; then
        echo "❌ docsify-cli 安装失败，请手动安装："
        echo "   npm install -g docsify-cli"
        exit 1
    fi
fi

# 获取当前脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "📖 文档目录: $SCRIPT_DIR"
echo "🌐 启动文档服务..."

# 启动 docsify 服务
docsify serve "$SCRIPT_DIR" --port 3000 --open

echo "✅ 文档服务已启动！"
echo "📱 访问地址: http://localhost:3000"
echo "🛑 按 Ctrl+C 停止服务" 