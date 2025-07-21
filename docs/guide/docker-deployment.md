# 🐳 Docker 部署

## 概述

urlDB 支持使用 Docker 进行容器化部署，提供了完整的前后端分离架构。

## 系统架构

| 服务 | 端口 | 说明 |
|------|------|------|
| frontend | 3000 | Nuxt.js 前端应用 |
| backend | 8080 | Go API 后端服务 |
| postgres | 5432 | PostgreSQL 数据库 |

## 快速部署

### 1. 克隆项目

```bash
git clone https://github.com/ctwj/urldb.git
cd urldb
```

### 2. 使用启动脚本（推荐）

```bash
# 给脚本执行权限
chmod +x docker-start.sh

# 启动服务
./docker-start.sh
```

### 3. 手动启动

```bash
# 构建并启动所有服务
docker compose up --build -d

# 查看服务状态
docker compose ps
```

## 配置说明

### 环境变量

可以通过修改 `docker-compose.yml` 文件中的环境变量来配置服务：

后端 backend
```yaml
environment:
  DB_HOST: postgres
  DB_PORT: 5432
  DB_USER: postgres
  DB_PASSWORD: password
  DB_NAME: url_db
  PORT: 8080
```

前端 frontend
```yaml
environment:
  API_BASE: /api
```

### 端口映射

如果需要修改端口映射，可以编辑 `docker-compose.yml`：

```yaml
ports:
  - "3001:3000"  # 前端端口
  - "8081:8080"  # API端口
  - "5433:5432"  # 数据库端口
```

## 常用命令

### 服务管理

```bash
# 启动服务
docker compose up -d

# 停止服务
docker compose down

# 重启服务
docker compose restart

# 查看服务状态
docker compose ps

# 查看日志
docker compose logs -f [service_name]
```

### 数据管理

```bash
# 备份数据库
docker compose exec postgres pg_dump -U postgres url_db > backup.sql

# 恢复数据库
docker compose exec -T postgres psql -U postgres url_db < backup.sql

# 进入数据库
docker compose exec postgres psql -U postgres url_db
```

### 容器管理

```bash
# 进入容器
docker compose exec [service_name] sh

# 查看容器资源使用
docker stats

# 清理未使用的资源
docker system prune -a
```

## 生产环境部署

### 1. 环境准备

```bash
# 安装 Docker 和 Docker Compose
# 确保服务器有足够资源（建议 4GB+ 内存）

# 创建部署目录
mkdir -p /opt/urldb
cd /opt/urldb
```

### 2. 配置文件

创建生产环境配置文件：

```bash
# 复制项目文件
git clone https://github.com/ctwj/urldb.git .

# 创建环境变量文件
cp env.example .env.prod

# 编辑生产环境配置
vim .env.prod
```

### 3. 启动服务

```bash
# 使用生产环境配置启动
docker compose -f docker-compose.yml --env-file .env.prod up -d

# 检查服务状态
docker compose ps
```

### 4. 配置反向代理

#### Nginx 配置示例

```nginx
server {
    listen 80;
    server_name your-domain.com;

    # 前端代理
    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # API 代理
    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 5. SSL 配置

```bash
# 使用 Let's Encrypt 获取证书
sudo certbot --nginx -d your-domain.com

# 或使用自签名证书
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
    -keyout /etc/ssl/private/urldb.key \
    -out /etc/ssl/certs/urldb.crt
```

## 监控和维护

### 1. 日志管理

```bash
# 查看所有服务日志
docker compose logs -f

# 查看特定服务日志
docker compose logs -f backend

# 导出日志
docker compose logs > urldb.log
```

### 2. 性能监控

```bash
# 查看容器资源使用
docker stats

# 查看系统资源
htop
df -h
free -h
```

### 3. 备份策略

```bash
#!/bin/bash
# 创建备份脚本 backup.sh

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/backup/urldb"

# 创建备份目录
mkdir -p $BACKUP_DIR

# 备份数据库
docker compose exec -T postgres pg_dump -U postgres url_db > $BACKUP_DIR/db_$DATE.sql

# 备份上传文件
tar -czf $BACKUP_DIR/uploads_$DATE.tar.gz uploads/

# 删除7天前的备份
find $BACKUP_DIR -name "*.sql" -mtime +7 -delete
find $BACKUP_DIR -name "*.tar.gz" -mtime +7 -delete
```

### 4. 自动更新

```bash
#!/bin/bash
# 创建更新脚本 update.sh

cd /opt/urldb

# 拉取最新代码
git pull origin main

# 重新构建并启动
docker compose down
docker compose up --build -d

# 检查服务状态
docker compose ps
```

## 故障排除

### 1. 服务启动失败

```bash
# 查看详细错误信息
docker compose logs [service_name]

# 检查端口占用
netstat -tulpn | grep :3000
netstat -tulpn | grep :8080

# 检查磁盘空间
df -h
```

### 2. 数据库连接问题

```bash
# 检查数据库状态
docker compose exec postgres pg_isready -U postgres

# 检查数据库日志
docker compose logs postgres

# 重启数据库服务
docker compose restart postgres
```

### 3. 前端无法访问后端

```bash
# 检查网络连接
docker compose exec frontend ping backend

# 检查 API 配置
docker compose exec frontend env | grep API_BASE

# 测试 API 连接
curl http://localhost:8080/api/health
```

### 4. 内存不足

```bash
# 查看内存使用
free -h

# 增加 swap 空间
sudo fallocate -l 2G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile
```

## 安全建议

### 1. 网络安全

- 使用防火墙限制端口访问
- 配置 SSL/TLS 加密
- 定期更新系统和 Docker 版本

### 2. 数据安全

- 定期备份数据库
- 使用强密码
- 限制数据库访问权限

### 3. 容器安全

- 使用非 root 用户运行容器
- 定期更新镜像
- 扫描镜像漏洞

## 下一步

- [了解系统配置](../guide/configuration.md)
- [查看 API 文档](../api/overview.md)
- [学习监控和维护](../development/deployment.md) 