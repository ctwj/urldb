#!/bin/bash

# 编译脚本 - 自动注入版本信息
# 用法: ./scripts/build.sh [target]

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 获取当前版本
get_current_version() {
    cat VERSION
}

# 获取Git信息
get_git_commit() {
    git rev-parse --short HEAD 2>/dev/null || echo "unknown"
}

get_git_branch() {
    git branch --show-current 2>/dev/null || echo "unknown"
}

# 获取构建时间
get_build_time() {
    date '+%Y-%m-%d %H:%M:%S'
}

# 编译函数
build() {
    local target=${1:-"main"}
    local version=$(get_current_version)
    local git_commit=$(get_git_commit)
    local git_branch=$(get_git_branch)
    local build_time=$(get_build_time)
    
    echo -e "${BLUE}开始编译...${NC}"
    echo -e "版本: ${GREEN}${version}${NC}"
    echo -e "Git提交: ${GREEN}${git_commit}${NC}"
    echo -e "Git分支: ${GREEN}${git_branch}${NC}"
    echo -e "构建时间: ${GREEN}${build_time}${NC}"
    
    # 构建 ldflags
    local ldflags="-X 'github.com/ctwj/urldb/utils.Version=${version}'"
    ldflags="${ldflags} -X 'github.com/ctwj/urldb/utils.BuildTime=${build_time}'"
    ldflags="${ldflags} -X 'github.com/ctwj/urldb/utils.GitCommit=${git_commit}'"
    ldflags="${ldflags} -X 'github.com/ctwj/urldb/utils.GitBranch=${git_branch}'"
    
    # 编译 - 使用跨平台编译设置
    echo -e "${YELLOW}编译中...${NC}"
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "${ldflags}" -o "${target}" .
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}编译成功!${NC}"
        echo -e "可执行文件: ${GREEN}${target}${NC}"
        echo -e "目标平台: ${GREEN}Linux${NC}"
        
        # 显示版本信息（在Linux环境下）
        echo -e "${BLUE}版本信息验证:${NC}"
        if [[ "$OSTYPE" == "linux-gnu"* ]]; then
            ./${target} version 2>/dev/null || echo "无法验证版本信息"
        else
            echo "当前非Linux环境，无法直接验证版本信息"
            echo "请将编译后的文件复制到Linux服务器上验证"
        fi
    else
        echo -e "${RED}编译失败!${NC}"
        exit 1
    fi
}

# 清理函数
clean() {
    echo -e "${YELLOW}清理编译文件...${NC}"
    rm -f main
    echo -e "${GREEN}清理完成${NC}"
}

# 显示帮助
show_help() {
    echo -e "${BLUE}编译脚本${NC}"
    echo ""
    echo "用法: $0 [命令]"
    echo ""
    echo "命令:"
    echo "  build [target]     编译程序 (当前平台)"
    echo "  build-linux [target] 编译Linux版本 (推荐)"
    echo "  clean              清理编译文件"
    echo "  help               显示此帮助信息"
    echo ""
    echo "示例:"
    echo "  $0                  # 编译Linux版本 (默认)"
    echo "  $0 build-linux      # 编译Linux版本"
    echo "  $0 build-linux app  # 编译Linux版本为 app"
    echo "  $0 build            # 编译当前平台版本"
    echo "  $0 clean            # 清理编译文件"
    echo ""
    echo "注意:"
    echo "  - Linux版本使用静态链接，适合部署到服务器"
    echo "  - 默认编译Linux版本，无需复制VERSION文件"
}

# Linux编译函数
build_linux() {
    local target=${1:-"main"}
    local version=$(get_current_version)
    local git_commit=$(get_git_commit)
    local git_branch=$(get_git_branch)
    local build_time=$(get_build_time)
    
    echo -e "${BLUE}开始Linux编译...${NC}"
    echo -e "版本: ${GREEN}${version}${NC}"
    echo -e "Git提交: ${GREEN}${git_commit}${NC}"
    echo -e "Git分支: ${GREEN}${git_branch}${NC}"
    echo -e "构建时间: ${GREEN}${build_time}${NC}"
    
    # 构建 ldflags
    local ldflags="-X 'github.com/ctwj/urldb/utils.Version=${version}'"
    ldflags="${ldflags} -X 'github.com/ctwj/urldb/utils.BuildTime=${build_time}'"
    ldflags="${ldflags} -X 'github.com/ctwj/urldb/utils.GitCommit=${git_commit}'"
    ldflags="${ldflags} -X 'github.com/ctwj/urldb/utils.GitBranch=${git_branch}'"
    
    # Linux编译
    echo -e "${YELLOW}编译中...${NC}"
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "${ldflags}" -o "${target}" .
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}Linux编译成功!${NC}"
        echo -e "可执行文件: ${GREEN}${target}${NC}"
        echo -e "目标平台: ${GREEN}Linux${NC}"
        echo -e "静态链接: ${GREEN}是${NC}"
        
        # 显示文件信息
        if command -v file >/dev/null 2>&1; then
            echo -e "${BLUE}文件信息:${NC}"
            file "${target}"
        fi
        
        echo -e "${BLUE}注意: 请在Linux服务器上验证版本信息${NC}"
    else
        echo -e "${RED}Linux编译失败!${NC}"
        exit 1
    fi
}

# 主函数
main() {
    case $1 in
        "build")
            build $2
            ;;
        "build-linux")
            build_linux $2
            ;;
        "clean")
            clean
            ;;
        "help"|"-h"|"--help")
            show_help
            ;;
        "")
            build_linux
            ;;
        *)
            echo -e "${RED}错误: 未知命令 '$1'${NC}"
            echo "使用 '$0 help' 查看帮助信息"
            exit 1
            ;;
    esac
}

# 运行主函数
main "$@" 