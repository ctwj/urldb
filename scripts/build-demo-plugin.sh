#!/bin/bash

# 构建演示插件的脚本

echo "构建演示插件..."

# 创建插件输出目录
mkdir -p ./plugins

# 检查是否在Linux环境下
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    echo "在Linux环境下构建插件..."

    # 构建演示插件
    echo "编译演示插件..."
    go build -buildmode=plugin -o ./plugins/demo-plugin.so ./plugin/demo/demo_plugin.go

    if [ $? -eq 0 ]; then
        echo "演示插件构建成功！"
        ls -la ./plugins/demo-plugin.so
    else
        echo "演示插件构建失败！"
        exit 1
    fi
else
    echo "当前环境是Windows，不支持plugin构建模式"
    echo "请在Linux环境下运行此脚本，或者使用Docker容器进行构建"

    # 提供Docker构建选项的说明
    echo ""
    echo "使用Docker构建插件的命令："
    echo "docker run --rm -v \${PWD}:/usr/src/myapp -w /usr/src/myapp golang:1.19 go build -buildmode=plugin -o ./plugins/demo-plugin.so ./plugin/demo/demo_plugin.go"
fi