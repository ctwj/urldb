# 前端构建阶段
FROM node:20-slim AS frontend-builder

# 安装pnpm
WORKDIR /app/web
COPY web/ ./
RUN npm install --frozen-lockfile
ARG NUXT_PUBLIC_API_SERVER=http://backend:8080/api
ARG NUXT_PUBLIC_API_CLIENT=/api
RUN npm run build

# 前端运行阶段
FROM node:20-alpine AS frontend

# RUN npm install -g pnpm
ENV NODE_ENV=production

WORKDIR /app
COPY --from=frontend-builder /app/web/.output ./.output
COPY --from=frontend-builder /app/web/package*.json ./
EXPOSE 3000
CMD ["node", ".output/server/index.mjs"]

# 后端构建阶段
FROM golang:1.24.5-alpine AS backend-builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# 复制所有源代码
COPY . .

# 定义构建参数
ARG VERSION
ARG GIT_COMMIT
ARG GIT_BRANCH
ARG BUILD_TIME

# 获取版本信息并编译
RUN VERSION=${VERSION:-$(cat VERSION)} && \
    GIT_COMMIT=${GIT_COMMIT:-$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")} && \
    GIT_BRANCH=${GIT_BRANCH:-$(git branch --show-current 2>/dev/null || echo "unknown")} && \
    BUILD_TIME=${BUILD_TIME:-$(date '+%Y-%m-%d %H:%M:%S')} && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -ldflags "-X 'github.com/ctwj/urldb/utils.Version=${VERSION}' \
              -X 'github.com/ctwj/urldb/utils.BuildTime=${BUILD_TIME}' \
              -X 'github.com/ctwj/urldb/utils.GitCommit=${GIT_COMMIT}' \
              -X 'github.com/ctwj/urldb/utils.GitBranch=${GIT_BRANCH}'" \
    -o main .

# 后端运行阶段
FROM alpine:latest AS backend

# 安装时区数据
RUN apk add --no-cache tzdata

WORKDIR /root/

# 复制后端二进制文件
COPY --from=backend-builder /app/main .

# 创建uploads目录
RUN mkdir -p uploads

# 设置环境变量
ENV GIN_MODE=release
ENV TZ=Asia/Shanghai

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./main"] 