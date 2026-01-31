/**
 * MCP配置JSON Schema定义
 * 用于配置验证和智能提示
 */

export const mcpJsonSchema = {
  $schema: 'http://json-schema.org/draft-07/schema#',
  type: 'object',
  properties: {
    mcpServers: {
      type: 'object',
      patternProperties: {
        '^[a-zA-Z][a-zA-Z0-9_-]*$': {
          type: 'object',
          properties: {
            command: {
              type: 'string',
              description: '命令行工具路径（stdio传输时必需）',
              examples: ['npx', 'python', 'node', 'docker']
            },
            args: {
              type: 'array',
              items: {
                type: 'string'
              },
              description: '命令行参数',
              examples: [['duckduckgo-websearch'], ['python', 'script.py'], ['node', 'server.js']]
            },
            env: {
              type: 'object',
              patternProperties: {
                '.*': {
                  type: 'string'
                }
              },
              description: '环境变量配置'
            },
            transport: {
              type: 'string',
              enum: ['stdio', 'http', 'https', 'sse'],
              description: '传输类型',
              default: 'stdio'
            },
            endpoint: {
              type: 'string',
              format: 'uri',
              description: '远程端点URL（http/https/sse传输时必需）',
              examples: ['https://api.example.com/mcp', 'http://localhost:3000/mcp']
            },
            headers: {
              type: 'object',
              patternProperties: {
                '.*': {
                  type: 'string'
                }
              },
              description: 'HTTP头部信息（http/https/sse传输时使用）',
              properties: {
                Authorization: {
                  type: 'string',
                  description: '认证头',
                  examples: ['Bearer token123', 'ApiKey abc123']
                },
                'Content-Type': {
                  type: 'string',
                  description: '内容类型',
                  default: 'application/json'
                }
              }
            },
            enabled: {
              type: 'boolean',
              description: '是否启用此服务',
              default: true
            },
            auto_start: {
              type: 'boolean',
              description: '是否自动启动此服务',
              default: true
            }
          },
          required: ['transport', 'enabled', 'auto_start'],
          allOf: [
            {
              if: {
                properties: {
                  transport: {
                    const: 'stdio'
                  }
                }
              },
              then: {
                required: ['command']
              },
              else: {}
            },
            {
              if: {
                properties: {
                  transport: {
                    enum: ['http', 'https', 'sse']
                  }
                }
              },
              then: {
                required: ['endpoint']
              },
              else: {}
            }
          ]
        }
      },
      description: 'MCP服务器配置'
    }
  },
  required: ['mcpServers'],
  additionalProperties: false
}

/**
 * MCP字段智能提示配置
 */
export const mcpCompletionConfig = {
  // 传输类型选项
  transportTypes: [
    { label: 'stdio - 本地进程通信', value: 'stdio' },
    { label: 'http - HTTP REST API', value: 'http' },
    { label: 'https - 安全HTTP连接', value: 'https' },
    { label: 'sse - Server-Sent Events', value: 'sse' }
  ],

  // 常见命令提示
  commonCommands: [
    { label: 'npx - Node包执行器', value: 'npx' },
    { label: 'python - Python解释器', value: 'python' },
    { label: 'python3 - Python3解释器', value: 'python3' },
    { label: 'node - Node.js运行时', value: 'node' },
    { label: 'docker - Docker容器', value: 'docker' },
    { label: 'python -m pip - Python包管理', value: 'python', args: ['-m', 'pip'] }
  ],

  // 常见MCP工具
  commonMCPTools: [
    { label: 'duckduckgo-websearch', value: 'duckduckgo-websearch' },
    { label: '@modelcontextprotocol/server-filesystem', value: '@modelcontextprotocol/server-filesystem' },
    { label: '@modelcontextprotocol/server-git', value: '@modelcontextprotocol/server-git' },
    { label: '@modelcontextprotocol/server-github', value: '@modelcontextprotocol/server-github' },
    { label: '@modelcontextprotocol/server-sqlite', value: '@modelcontextprotocol/server-sqlite' },
    { label: '@modelcontextprotocol/server-postgres', value: '@modelcontextprotocol/server-postgres' },
    { label: 'time-mcp', value: 'time-mcp' }
  ],

  // 环境变量常见配置
  commonEnvVars: [
    { label: 'API_KEY - API密钥', key: 'API_KEY', description: 'API访问密钥' },
    { label: 'BING_API_KEY - Bing搜索API', key: 'BING_API_KEY', description: 'Bing搜索API密钥' },
    { label: 'OPENAI_API_KEY - OpenAI API', key: 'OPENAI_API_KEY', description: 'OpenAI API密钥' },
    { label: 'NODE_ENV - Node环境', key: 'NODE_ENV', value: 'production' },
    { label: 'DEBUG - 调试模式', key: 'DEBUG', value: 'true' }
  ],

  // HTTP头部常见配置
  commonHeaders: [
    { label: 'Authorization - Bearer Token', key: 'Authorization', value: 'Bearer your-token' },
    { label: 'Authorization - API Key', key: 'Authorization', value: 'ApiKey your-api-key' },
    { label: 'Content-Type - JSON', key: 'Content-Type', value: 'application/json' },
    { label: 'User-Agent - 自定义', key: 'User-Agent', value: 'MCP-Client/1.0' },
    { label: 'Accept - JSON响应', key: 'Accept', value: 'application/json' }
  ]
}

/**
 * 验证MCP配置的逻辑函数
 */
export const validateMCPConfig = (config: any): { valid: boolean; errors: string[]; warnings: string[] } => {
  const errors: string[] = []
  const warnings: string[] = []

  try {
    if (!config || typeof config !== 'object') {
      errors.push('配置必须是一个有效的JSON对象')
      return { valid: false, errors, warnings }
    }

    if (!config.mcpServers) {
      errors.push('缺少必需的字段: mcpServers')
      return { valid: false, errors, warnings }
    }

    if (typeof config.mcpServers !== 'object' || Array.isArray(config.mcpServers)) {
      errors.push('mcpServers 必须是一个对象')
      return { valid: false, errors, warnings }
    }

    // 验证每个服务器配置
    Object.entries(config.mcpServers).forEach(([serverName, serverConfig]: [string, any]) => {
      if (!serverConfig || typeof serverConfig !== 'object') {
        errors.push(`服务器 "${serverName}" 的配置必须是一个对象`)
        return
      }

      // 验证传输类型
      if (!serverConfig.transport) {
        errors.push(`服务器 "${serverName}" 缺少必需的字段: transport`)
      } else if (!['stdio', 'http', 'https', 'sse'].includes(serverConfig.transport)) {
        errors.push(`服务器 "${serverName}" 的传输类型 "${serverConfig.transport}" 无效`)
      }

      // 验证stdio配置
      if (serverConfig.transport === 'stdio') {
        if (!serverConfig.command) {
          errors.push(`服务器 "${serverName}" 缺少必需的字段: command`)
        }
      }

      // 验证HTTP/SSE配置
      if (['http', 'https', 'sse'].includes(serverConfig.transport)) {
        if (!serverConfig.endpoint) {
          errors.push(`服务器 "${serverName}" 缺少必需的字段: endpoint`)
        } else {
          try {
            new URL(serverConfig.endpoint)
          } catch {
            errors.push(`服务器 "${serverName}" 的端点URL格式无效: ${serverConfig.endpoint}`)
          }

          if (serverConfig.transport === 'https' && !serverConfig.endpoint.startsWith('https://')) {
            warnings.push(`服务器 "${serverName}" 使用https传输但端点不是https协议`)
          }
        }
      }

      // 检查必需的字段
      if (serverConfig.enabled === undefined) {
        errors.push(`服务器 "${serverName}" 缺少必需的字段: enabled`)
      } else if (typeof serverConfig.enabled !== 'boolean') {
        warnings.push(`服务器 "${serverName}" 的 enabled 字段应该是布尔值`)
      }

      if (serverConfig.auto_start === undefined) {
        errors.push(`服务器 "${serverName}" 缺少必需的字段: auto_start`)
      } else if (typeof serverConfig.auto_start !== 'boolean') {
        warnings.push(`服务器 "${serverName}" 的 auto_start 字段应该是布尔值`)
      }
    })

  } catch (error) {
    errors.push(`配置验证时发生错误: ${error}`)
  }

  return {
    valid: errors.length === 0,
    errors,
    warnings
  }
}