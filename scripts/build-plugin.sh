#!/bin/bash

# 插件构建脚本

PLUGIN_DIR="./examples/go-plugins/demo"
OUTPUT_DIR="./plugins"

# 创建输出目录
mkdir -p "$OUTPUT_DIR"

# 检查是否在Linux环境下
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    echo "在Linux环境下构建插件..."

    # 构建演示插件
    echo "构建演示插件..."
    go build -buildmode=plugin -o "$OUTPUT_DIR/demo-plugin.so" "$PLUGIN_DIR/demo_plugin.go"

    if [ $? -eq 0 ]; then
        echo "插件构建成功！"
        ls -la "$OUTPUT_DIR/demo-plugin.so"
    else
        echo "插件构建失败！"
        exit 1
    fi
else
    echo "当前环境不支持plugin构建模式"
    echo "请在Linux环境下运行此脚本"

    # 在Windows环境下，我们创建一个简单的说明文件
    echo "在Windows环境下，请使用交叉编译构建插件：" > "$OUTPUT_DIR/README.md"
    echo "GOOS=linux GOARCH=amd64 go build -buildmode=plugin -o $OUTPUT_DIR/demo-plugin.so $PLUGIN_DIR/demo_plugin.go" >> "$OUTPUT_DIR/README.md"
    echo "" >> "$OUTPUT_DIR/README.md"
    echo "或者使用Docker容器进行构建。" >> "$OUTPUT_DIR/README.md"
fi