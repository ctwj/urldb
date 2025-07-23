# 前端构建阶段
FROM node:20-slim AS frontend-builder

# 安装pnpm
WORKDIR /app/web
COPY web/ ./
RUN npm install --frozen-lockfile
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

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 后端运行阶段
FROM alpine:latest AS backend

WORKDIR /root/

# 复制后端二进制文件
COPY --from=backend-builder /app/main .

# 创建uploads目录
RUN mkdir -p uploads

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./main"] 