# GitHub版本管理指南

本项目使用GitHub进行版本管理，支持自动创建Release和标签。

## 版本管理流程

### 1. 版本号规范

遵循[语义化版本](https://semver.org/lang/zh-CN/)规范：

- **主版本号** (Major): 不兼容的API修改
- **次版本号** (Minor): 向下兼容的功能性新增  
- **修订号** (Patch): 向下兼容的问题修正

### 2. 版本管理命令

#### 显示版本信息
```bash
./scripts/version.sh show
```

#### 更新版本号
```bash
# 修订版本 (1.0.0 -> 1.0.1)
./scripts/version.sh patch

# 次版本 (1.0.0 -> 1.1.0)
./scripts/version.sh minor

# 主版本 (1.0.0 -> 2.0.0)
./scripts/version.sh major
```

#### 发布版本到GitHub
```bash
./scripts/version.sh release
```

### 3. 自动发布流程

当执行版本更新命令时，脚本会：

1. **更新版本号**: 修改 `VERSION` 文件
2. **同步文件**: 更新 `package.json`、`docker-compose.yml`、`README.md`
3. **创建Git标签**: 自动创建版本标签
4. **推送代码**: 推送代码和标签到GitHub
5. **创建Release**: 自动创建GitHub Release

### 4. 手动发布流程

如果自动发布失败，可以手动发布：

#### 步骤1: 更新版本号
```bash
./scripts/version.sh patch  # 或 minor, major
```

#### 步骤2: 提交更改
```bash
git add .
git commit -m "chore: bump version to v1.0.1"
```

#### 步骤3: 创建标签
```bash
git tag v1.0.1
```

#### 步骤4: 推送到GitHub
```bash
git push origin main
git push origin v1.0.1
```

#### 步骤5: 创建Release
在GitHub网页上：
1. 进入项目页面
2. 点击 "Releases"
3. 点击 "Create a new release"
4. 选择标签 `v1.0.1`
5. 填写Release说明
6. 发布

### 5. GitHub CLI工具

#### 安装GitHub CLI
```bash
# macOS
brew install gh

# Ubuntu/Debian
sudo apt install gh

# Windows
winget install GitHub.cli
```

#### 登录GitHub
```bash
gh auth login
```

#### 创建Release
```bash
gh release create v1.0.1 \
  --title "Release v1.0.1" \
  --notes "修复了一些bug" \
  --draft=false \
  --prerelease=false
```

### 6. 版本检查

#### API接口
- `GET /api/version/check-update` - 检查GitHub上的最新版本

#### 前端页面
- 访问 `/version` 页面查看版本信息和更新状态

### 7. 版本历史

#### 查看所有标签
```bash
git tag -l
```

#### 查看标签详情
```bash
git show v1.0.1
```

#### 查看版本历史
```bash
git log --oneline --decorate
```

### 8. 回滚版本

如果需要回滚到之前的版本：

#### 删除本地标签
```bash
git tag -d v1.0.1
```

#### 删除远程标签
```bash
git push origin :refs/tags/v1.0.1
```

#### 回滚代码
```bash
git reset --hard v1.0.0
git push --force origin main
```

### 9. 最佳实践

#### 提交信息规范
```bash
# 功能开发
git commit -m "feat: 添加新功能"

# Bug修复
git commit -m "fix: 修复某个bug"

# 文档更新
git commit -m "docs: 更新文档"

# 版本更新
git commit -m "chore: bump version to v1.0.1"
```

#### 分支管理
- `main`: 主分支，用于发布
- `develop`: 开发分支
- `feature/*`: 功能分支
- `hotfix/*`: 热修复分支

#### Release说明模板
```markdown
## Release v1.0.1

**发布日期**: 2024-01-15

### 更新内容

- 修复了某个bug
- 添加了新功能
- 优化了性能

### 下载

- [源码 (ZIP)](https://github.com/ctwj/urldb/archive/v1.0.1.zip)
- [源码 (TAR.GZ)](https://github.com/ctwj/urldb/archive/v1.0.1.tar.gz)

### 安装

```bash
# 克隆项目
git clone https://github.com/ctwj/urldb.git
cd urldb

# 切换到指定版本
git checkout v1.0.1

# 使用Docker部署
docker-compose up --build -d
```

### 更新日志

详细更新日志请查看 [CHANGELOG.md](https://github.com/ctwj/urldb/blob/v1.0.1/CHANGELOG.md)
```

### 10. 故障排除

#### 常见问题

1. **GitHub CLI未安装**
   ```bash
   # 安装GitHub CLI
   brew install gh  # macOS
   ```

2. **GitHub CLI未登录**
   ```bash
   # 登录GitHub
   gh auth login
   ```

3. **标签已存在**
   ```bash
   # 删除本地标签
   git tag -d v1.0.1
   
   # 删除远程标签
   git push origin :refs/tags/v1.0.1
   ```

4. **推送失败**
   ```bash
   # 检查远程仓库
   git remote -v
   
   # 重新设置远程仓库
   git remote set-url origin https://github.com/ctwj/urldb.git
   ```

#### 获取帮助
```bash
./scripts/version.sh help
``` 