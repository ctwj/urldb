# Product Overview

urlDB (老九网盘资源数据库) is a modern cloud storage resource database system that supports automated transfer and sharing across multiple cloud storage platforms.

## Core Features

- **Multi-platform support**: Supports 8+ cloud storage platforms including Baidu Pan, Aliyun Pan, Quark Pan, UC Pan, Xunlei Pan, 123 Pan, 115 Pan, and Tianyi Pan
- **Automated resource processing**: Automatic validity checking, transfer, and sharing of resources
- **Multi-account management**: Support for multiple accounts per platform
- **Public API**: RESTful API for resource ingestion and search
- **Search integration**: Meilisearch integration for fast full-text search
- **Bot integrations**: Telegram bot and WeChat official account auto-reply
- **SEO optimization**: Sitemap generation, Google indexing, Bing submission
- **Admin dashboard**: Full-featured web interface for resource and system management

## Target Users

- Content aggregators managing cloud storage resources
- Teams needing automated cloud storage workflows
- Services providing resource search and discovery

## Architecture

Full-stack application with:
- Go backend (Gin framework)
- Nuxt.js 3 frontend (Vue 3 + TypeScript)
- PostgreSQL database
- Optional Meilisearch for enhanced search
