# 多阶段构建
FROM node:18-alpine AS frontend-builder

WORKDIR /app/web
COPY web/package*.json ./
RUN npm ci --only=production

COPY web/ ./
RUN npm run build

FROM golang:1.21-alpine AS backend-builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 复制后端二进制文件
COPY --from=backend-builder /app/main .

# 复制前端构建文件
COPY --from=frontend-builder /app/web/.output /root/web

# 创建uploads目录
RUN mkdir -p uploads

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./main"] 