# 编译说明

## 方案1：使用编译脚本（推荐）

### 在Git Bash中执行：

```bash
# 给脚本添加执行权限（首次使用）
chmod +x scripts/build.sh

# 编译Linux版本（推荐，用于服务器部署）
./scripts/build.sh

# 或者明确指定编译Linux版本
./scripts/build.sh build-linux

# 或者指定目标文件名
./scripts/build.sh build-linux myapp

# 编译当前平台版本（用于本地测试）
./scripts/build.sh build
```

### 编译脚本功能：
- 自动读取 `VERSION` 文件中的版本号
- 自动获取Git提交信息和分支信息
- 自动获取构建时间
- 将版本信息编译到可执行文件中
- 支持跨平台编译（默认编译Linux版本）
- 使用静态链接，适合服务器部署

## 方案2：手动编译

### Linux版本（推荐）：

```bash
# 获取版本信息
VERSION=$(cat VERSION)
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GIT_BRANCH=$(git branch --show-current 2>/dev/null || echo "unknown")
BUILD_TIME=$(date '+%Y-%m-%d %H:%M:%S')

# 编译Linux版本
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-X 'github.com/ctwj/urldb/utils.Version=${VERSION}' -X 'github.com/ctwj/urldb/utils.BuildTime=${BUILD_TIME}' -X 'github.com/ctwj/urldb/utils.GitCommit=${GIT_COMMIT}' -X 'github.com/ctwj/urldb/utils.GitBranch=${GIT_BRANCH}'" -o main .
```

### 当前平台版本：

```bash
# 获取版本信息
VERSION=$(cat VERSION)
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GIT_BRANCH=$(git branch --show-current 2>/dev/null || echo "unknown")
BUILD_TIME=$(date '+%Y-%m-%d %H:%M:%S')

# 编译当前平台版本
go build -ldflags "-X 'github.com/ctwj/urldb/utils.Version=${VERSION}' -X 'github.com/ctwj/urldb/utils.BuildTime=${BUILD_TIME}' -X 'github.com/ctwj/urldb/utils.GitCommit=${GIT_COMMIT}' -X 'github.com/ctwj/urldb/utils.GitBranch=${GIT_BRANCH}'" -o main .
```

## 验证版本信息

编译完成后，可以通过以下方式验证版本信息：

```bash
# 命令行验证
./main version

# 启动服务器后通过API验证
curl http://localhost:8080/api/version
```

## 部署说明

使用方案1编译后，部署时只需要：

1. 复制可执行文件到服务器
2. 启动程序

**不再需要复制 `VERSION` 文件**，因为版本信息已经编译到程序中。

### 使用部署脚本（可选）

```bash
# 给部署脚本添加执行权限
chmod +x scripts/deploy-example.sh

# 部署到服务器
./scripts/deploy-example.sh root example.com /opt/urldb
```

### 使用Docker构建脚本：

```bash
# 给脚本添加执行权限
chmod +x scripts/docker-build.sh

# 构建Docker镜像
./scripts/docker-build.sh build

# 构建指定版本镜像
./scripts/docker-build.sh build 1.2.4

# 推送镜像到Docker Hub
./scripts/docker-build.sh push 1.2.4
```

### 手动Docker构建：

```bash
# 构建镜像
docker build --target backend -t ctwj/urldb-backend:1.2.3 .
docker build --target frontend -t ctwj/urldb-frontend:1.2.3 .
```

## 版本管理

更新版本号：

```bash
# 更新版本号
./scripts/version.sh patch  # 修订版本
./scripts/version.sh minor  # 次版本
./scripts/version.sh major  # 主版本

# 然后重新编译
./scripts/build.sh

# 或者构建Docker镜像
./scripts/docker-build.sh build
``` 