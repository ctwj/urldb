module.exports = {
  apps: [
    {
      name: 'urldb-nuxt',
      port: '3030',
      exec_mode: 'cluster',
      instances: 'max', // 使用所有可用的CPU核心
      script: './.output/server/index.mjs',
      error_file: './logs/err.log',
      out_file: './logs/out.log',
      log_file: './logs/combined.log',
      time: true,
      env: {
        NODE_ENV: 'production',
        HOST: '0.0.0.0',
        PORT: 3030,
        NUXT_PUBLIC_API_SERVER: 'http://localhost:8080/api'
      }
    }
  ]
};