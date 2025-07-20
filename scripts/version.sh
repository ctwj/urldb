#!/bin/bash

# 版本管理脚本
# 用法: ./scripts/version.sh [major|minor|patch|show|update]

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

# 显示版本信息
show_version() {
    echo -e "${BLUE}当前版本信息:${NC}"
    echo -e "版本号: ${GREEN}$(get_current_version)${NC}"
    echo -e "构建时间: ${GREEN}$(date '+%Y-%m-%d %H:%M:%S')${NC}"
    echo -e "Git提交: ${GREEN}$(git rev-parse --short HEAD 2>/dev/null || echo 'N/A')${NC}"
    echo -e "Git分支: ${GREEN}$(git branch --show-current 2>/dev/null || echo 'N/A')${NC}"
}

# 更新版本号
update_version() {
    local version_type=$1
    local current_version=$(get_current_version)
    local major minor patch
    
    # 解析版本号
    IFS='.' read -r major minor patch <<< "$current_version"
    
    case $version_type in
        "major")
            major=$((major + 1))
            minor=0
            patch=0
            ;;
        "minor")
            minor=$((minor + 1))
            patch=0
            ;;
        "patch")
            patch=$((patch + 1))
            ;;
        *)
            echo -e "${RED}错误: 无效的版本类型${NC}"
            echo "用法: $0 [major|minor|patch|show|update|release]"
            exit 1
            ;;
    esac
    
    local new_version="$major.$minor.$patch"
    
    # 更新版本文件
    echo "$new_version" > VERSION
    
    echo -e "${GREEN}版本已更新: ${current_version} -> ${new_version}${NC}"
    
    # 更新其他文件中的版本信息
    update_version_in_files "$new_version"
    
    # 创建Git标签和发布
    if git rev-parse --git-dir > /dev/null 2>&1; then
        echo -e "${YELLOW}创建Git标签 v${new_version}...${NC}"
        git add VERSION
        git commit -m "chore: bump version to v${new_version}" || true
        git tag "v${new_version}" || true
        echo -e "${GREEN}Git标签已创建: v${new_version}${NC}"
        
        # 询问是否推送到GitHub
        read -p "是否推送到GitHub并创建Release? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            push_to_github "$new_version"
        fi
    fi
}

# 更新文件中的版本信息
update_version_in_files() {
    local new_version=$1
    
    echo -e "${YELLOW}更新文件中的版本信息...${NC}"
    
    # 更新前端package.json
    if [ -f "web/package.json" ]; then
        sed -i.bak "s/\"version\": \"[^\"]*\"/\"version\": \"${new_version}\"/" web/package.json
        rm -f web/package.json.bak
        echo -e "  ✅ 更新 web/package.json"
    fi
    
    # 更新useVersion.ts中的默认版本
    if [ -f "web/composables/useVersion.ts" ]; then
        # 使用更简单的模式匹配，先获取当前版本号
        current_use_version=$(grep -o "version: '[0-9]\+\.[0-9]\+\.[0-9]\+'" web/composables/useVersion.ts | head -1)
        if [ -n "$current_use_version" ]; then
            sed -i.bak "s/$current_use_version/version: '${new_version}'/" web/composables/useVersion.ts
            rm -f web/composables/useVersion.ts.bak
            echo -e "  ✅ 更新 web/composables/useVersion.ts"
        else
            echo -e "  ⚠️  未找到useVersion.ts中的版本号"
        fi
    fi
    
    # 更新Docker镜像标签
    if [ -f "docker-compose.yml" ]; then
        # 获取当前镜像版本
        current_backend_version=$(grep -o "ctwj/urldb-backend:[0-9]\+\.[0-9]\+\.[0-9]\+" docker-compose.yml | head -1)
        current_frontend_version=$(grep -o "ctwj/urldb-frontend:[0-9]\+\.[0-9]\+\.[0-9]\+" docker-compose.yml | head -1)
        
        if [ -n "$current_backend_version" ]; then
            sed -i.bak "s|$current_backend_version|ctwj/urldb-backend:${new_version}|" docker-compose.yml
            echo -e "  ✅ 更新 backend 镜像: ${current_backend_version} -> ctwj/urldb-backend:${new_version}"
        fi
        
        if [ -n "$current_frontend_version" ]; then
            sed -i.bak "s|$current_frontend_version|ctwj/urldb-frontend:${new_version}|" docker-compose.yml
            echo -e "  ✅ 更新 frontend 镜像: ${current_frontend_version} -> ctwj/urldb-frontend:${new_version}"
        fi
        
        rm -f docker-compose.yml.bak
        echo -e "  ✅ 更新 docker-compose.yml 完成"
    fi
    
    # 更新README中的版本信息
    if [ -f "README.md" ]; then
        sed -i.bak "s/版本.*[0-9]\+\.[0-9]\+\.[0-9]\+/版本 ${new_version}/" README.md
        rm -f README.md.bak
        echo -e "  ✅ 更新 README.md"
    fi
}

# 推送到GitHub并创建Release
push_to_github() {
    local version=$1
    local release_notes=$(generate_release_notes "$version")
    
    echo -e "${YELLOW}推送到GitHub...${NC}"
    
    # 推送代码和标签
    git push origin main || git push origin master || true
    git push origin "v${version}" || true
    
    echo -e "${GREEN}代码和标签已推送到GitHub${NC}"
    
    # 创建GitHub Release
    create_github_release "$version" "$release_notes"
}

# 创建GitHub Release
create_github_release() {
    local version=$1
    local release_notes=$2
    
    # 检查是否安装了gh CLI
    if ! command -v gh &> /dev/null; then
        echo -e "${YELLOW}未安装GitHub CLI (gh)，跳过自动创建Release${NC}"
        echo -e "${BLUE}请手动在GitHub上创建Release: v${version}${NC}"
        return
    fi
    
    # 检查是否已登录
    if ! gh auth status &> /dev/null; then
        echo -e "${YELLOW}GitHub CLI未登录，跳过自动创建Release${NC}"
        echo -e "${BLUE}请运行 'gh auth login' 登录后重试${NC}"
        return
    fi
    
    echo -e "${YELLOW}创建GitHub Release...${NC}"
    
    # 创建Release
    gh release create "v${version}" \
        --title "Release v${version}" \
        --notes "$release_notes" \
        --draft=false \
        --prerelease=false || {
        echo -e "${RED}创建Release失败，请手动创建${NC}"
        return
    }
    
    echo -e "${GREEN}GitHub Release已创建: v${version}${NC}"
}

# 生成Release说明
generate_release_notes() {
    local version=$1
    local current_date=$(date '+%Y-%m-%d')
    
    # 获取上次版本
    local previous_version=$(git describe --tags --abbrev=0 2>/dev/null | sed 's/v//' || echo "0.0.0")
    
    # 获取提交历史
    local commits=$(git log --oneline "v${previous_version}..HEAD" 2>/dev/null || echo "Initial release")
    
    cat << EOF
## Release v${version}

**发布日期**: ${current_date}

### 更新内容

${commits}

### 下载

- [源码 (ZIP)](https://github.com/ctwj/urldb/archive/v${version}.zip)
- [源码 (TAR.GZ)](https://github.com/ctwj/urldb/archive/v${version}.tar.gz)

### 安装

\`\`\`bash
# 克隆项目
git clone https://github.com/ctwj/urldb.git
cd urldb

# 切换到指定版本
git checkout v${version}

# 使用Docker部署
docker-compose up --build -d
\`\`\`

### 更新日志

详细更新日志请查看 [CHANGELOG.md](https://github.com/ctwj/urldb/blob/v${version}/CHANGELOG.md)
EOF
}

# 生成版本信息文件
generate_version_info() {
    local version=$(get_current_version)
    local build_time=$(date '+%Y-%m-%d %H:%M:%S')
    local git_commit=$(git rev-parse --short HEAD 2>/dev/null || echo 'N/A')
    local git_branch=$(git branch --show-current 2>/dev/null || echo 'N/A')
    
    cat > version_info.json << EOF
{
  "version": "${version}",
  "build_time": "${build_time}",
  "git_commit": "${git_commit}",
  "git_branch": "${git_branch}",
  "go_version": "$(go version 2>/dev/null | cut -d' ' -f3 || echo 'N/A')",
  "node_version": "$(node --version 2>/dev/null || echo 'N/A')"
}
EOF
    
    echo -e "${GREEN}版本信息文件已生成: version_info.json${NC}"
}

# 主函数
main() {
    case $1 in
        "show")
            show_version
            ;;
        "major"|"minor"|"patch")
            update_version $1
            ;;
        "release")
            release_version
            ;;
        "update")
            generate_version_info
            ;;
        "help"|"-h"|"--help")
            echo -e "${BLUE}版本管理脚本${NC}"
            echo ""
            echo "用法: $0 [命令]"
            echo ""
            echo "命令:"
            echo "  show    显示当前版本信息"
            echo "  major   主版本号更新 (x.0.0)"
            echo "  minor   次版本号更新 (0.x.0)"
            echo "  patch   修订版本号更新 (0.0.x)"
            echo "  release 发布当前版本到GitHub"
            echo "  update  生成版本信息文件"
            echo "  help    显示此帮助信息"
            echo ""
            echo "示例:"
            echo "  $0 show          # 显示版本信息"
            echo "  $0 patch         # 更新修订版本号"
            echo "  $0 minor         # 更新次版本号"
            echo "  $0 major         # 更新主版本号"
            echo "  $0 release       # 发布版本到GitHub"
            ;;
        *)
            echo -e "${RED}错误: 未知命令 '$1'${NC}"
            echo "使用 '$0 help' 查看帮助信息"
            exit 1
            ;;
    esac
}

# 发布版本
release_version() {
    local current_version=$(get_current_version)
    
    echo -e "${BLUE}准备发布版本 v${current_version}${NC}"
    
    # 检查是否有未提交的更改
    if ! git diff-index --quiet HEAD --; then
        echo -e "${YELLOW}检测到未提交的更改，请先提交更改${NC}"
        git status --short
        return 1
    fi
    
    # 检查是否已存在该版本的标签
    if git tag -l "v${current_version}" | grep -q "v${current_version}"; then
        echo -e "${YELLOW}版本 v${current_version} 的标签已存在${NC}"
        read -p "是否继续发布? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            return 1
        fi
    fi
    
    # 创建标签
    echo -e "${YELLOW}创建Git标签 v${current_version}...${NC}"
    git tag "v${current_version}" || {
        echo -e "${RED}创建标签失败${NC}"
        return 1
    }
    
    # 推送到GitHub
    push_to_github "$current_version"
}

# 运行主函数
main "$@" 