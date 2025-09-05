#!/bin/bash

# Docker构建脚本
# 用法: ./scripts/docker-build.sh [version]

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 获取版本号
get_version() {
    if [ -n "$1" ]; then
        echo "$1"
    else
        cat VERSION
    fi
}

# 获取Git信息
get_git_commit() {
    git rev-parse --short HEAD 2>/dev/null || echo "unknown"
}

get_git_branch() {
    git branch --show-current 2>/dev/null || echo "unknown"
}

# 构建Docker镜像
build_docker() {
    local version=$(get_version $1)
    local skip_frontend=$2
    local git_commit=$(get_git_commit)
    local git_branch=$(get_git_branch)
    local build_time=$(date '+%Y-%m-%d %H:%M:%S')
    
    echo -e "${BLUE}开始Docker构建...${NC}"
    echo -e "版本: ${GREEN}${version}${NC}"
    echo -e "Git提交: ${GREEN}${git_commit}${NC}"
    echo -e "Git分支: ${GREEN}${git_branch}${NC}"
    echo -e "构建时间: ${GREEN}${build_time}${NC}"
    if [ "$skip_frontend" = "true" ]; then
        echo -e "跳过前端构建: ${GREEN}是${NC}"
    fi
    
    # 直接使用 docker build，避免 buildx 的复杂性
    BUILD_CMD="docker build"
    echo -e "${BLUE}使用构建命令: ${BUILD_CMD}${NC}"

    # 构建前端镜像（可选）
    if [ "$skip_frontend" != "true" ]; then
        echo -e "${YELLOW}构建前端镜像...${NC}"
        FRONTEND_CMD="${BUILD_CMD} --build-arg VERSION=${version} --build-arg GIT_COMMIT=${git_commit} --build-arg GIT_BRANCH=${git_branch} --build-arg BUILD_TIME=${build_time} --target frontend -t ctwj/urldb-frontend:${version} ."
        echo -e "${BLUE}执行命令: ${FRONTEND_CMD}${NC}"
        ${BUILD_CMD} \
            --build-arg VERSION=${version} \
            --build-arg GIT_COMMIT=${git_commit} \
            --build-arg GIT_BRANCH=${git_branch} \
            --build-arg "BUILD_TIME=${build_time}" \
            --target frontend \
            -t ctwj/urldb-frontend:${version} \
            .
        
        if [ $? -ne 0 ]; then
            echo -e "${RED}前端构建失败!${NC}"
            exit 1
        fi
    else
        echo -e "${YELLOW}跳过前端构建${NC}"
    fi
    
    # 构建后端镜像
    echo -e "${YELLOW}构建后端镜像...${NC}"
    BACKEND_CMD="${BUILD_CMD} --build-arg VERSION=${version} --build-arg GIT_COMMIT=${git_commit} --build-arg GIT_BRANCH=${git_branch} --build-arg BUILD_TIME=${build_time} --target backend -t ctwj/urldb-backend:${version} ."
    echo -e "${BLUE}执行命令: ${BACKEND_CMD}${NC}"
    ${BUILD_CMD} \
        --build-arg VERSION=${version} \
        --build-arg GIT_COMMIT=${git_commit} \
        --build-arg GIT_BRANCH=${git_branch} \
        --build-arg BUILD_TIME="${build_time}" \
        --target backend \
        -t ctwj/urldb-backend:${version} \
        .
    
    if [ $? -ne 0 ]; then
        echo -e "${RED}后端构建失败!${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}Docker构建完成!${NC}"
    echo -e "镜像标签:"
    echo -e "  ${GREEN}ctwj/urldb-backend:${version}${NC}"
    if [ "$skip_frontend" != "true" ]; then
        echo -e "  ${GREEN}ctwj/urldb-frontend:${version}${NC}"
    fi
}

# 推送镜像
push_images() {
    local version=$(get_version $1)
    
    echo -e "${YELLOW}推送镜像到Docker Hub...${NC}"
    
    # 推送后端镜像
    docker push ctwj/urldb-backend:${version}
    
    # 推送前端镜像
    docker push ctwj/urldb-frontend:${version}
    
    echo -e "${GREEN}镜像推送完成!${NC}"
}

# 清理镜像
clean_images() {
    local version=$(get_version $1)
    
    echo -e "${YELLOW}清理Docker镜像...${NC}"
    docker rmi ctwj/urldb-backend:${version} 2>/dev/null || true
    docker rmi ctwj/urldb-frontend:${version} 2>/dev/null || true
    
    echo -e "${GREEN}镜像清理完成${NC}"
}

# 显示帮助
show_help() {
    echo -e "${BLUE}Docker构建脚本${NC}"
    echo ""
    echo "用法: $0 [命令] [版本] [选项]"
    echo ""
    echo "命令:"
    echo "  build [version] [--skip-frontend]  构建Docker镜像"
    echo "  push [version]                     推送镜像到Docker Hub"
    echo "  clean [version]                    清理Docker镜像"
    echo "  help                               显示此帮助信息"
    echo ""
    echo "选项:"
    echo "  --skip-frontend                    跳过前端构建"
    echo ""
    echo "示例:"
    echo "  $0 build                           # 构建当前版本镜像"
    echo "  $0 build 1.2.4                     # 构建指定版本镜像"
    echo "  $0 build 1.2.4 --skip-frontend     # 构建指定版本镜像，跳过前端"
    echo "  $0 push 1.2.4                      # 推送指定版本镜像"
    echo "  $0 clean                           # 清理当前版本镜像"
}

# 主函数
main() {
    case $1 in
        "build")
            # 检查是否有 --skip-frontend 选项
            local skip_frontend="false"
            if [ "$3" = "--skip-frontend" ]; then
                skip_frontend="true"
            fi
            build_docker $2 $skip_frontend
            ;;
        "push")
            push_images $2
            ;;
        "clean")
            clean_images $2
            ;;
        "help"|"-h"|"--help")
            show_help
            ;;
        "")
            show_help
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