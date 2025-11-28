interface VersionInfo {
  version: string
  build_time: string
  git_commit: string
  git_branch: string
  go_version: string
  node_version: string
  platform: string
  arch: string
}

interface VersionResponse {
  success: boolean
  data: VersionInfo
  message: string
  time: string
}

export const useVersion = () => {
  const versionInfo = ref<VersionInfo>({
    version: '1.3.5',
    build_time: '',
    git_commit: 'unknown',
    git_branch: 'unknown',
    go_version: '',
    node_version: '',
    platform: '',
    arch: ''
  })

  const loading = ref(false)
  const error = ref('')

  // 获取版本信息
  const fetchVersionInfo = async () => {
    loading.value = true
    error.value = ''
    
    try {
      const response = await $fetch('/api/version') as VersionResponse
      if (response.success) {
        versionInfo.value = response.data
      } else {
        error.value = response.message || '获取版本信息失败'
      }
    } catch (err: any) {
      error.value = err.message || '网络错误'
      console.error('获取版本信息失败:', err)
    } finally {
      loading.value = false
    }
  }

  // 获取版本字符串
  const getVersionString = async () => {
    try {
      const response = await $fetch('/api/version/string') as any
      if (response.success) {
        return response.data.version
      }
    } catch (err) {
      console.error('获取版本字符串失败:', err)
    }
    return versionInfo.value.version
  }

  // 检查更新
  const checkUpdate = async () => {
    try {
      const response = await $fetch('/api/version/check-update') as any
      if (response.success) {
        return response.data
      }
    } catch (err) {
      console.error('检查更新失败:', err)
    }
    return null
  }

  // 格式化版本信息
  const formatVersionInfo = computed(() => {
    const info = versionInfo.value
    return {
      version: info.version,
      gitCommit: info.git_commit !== 'unknown' ? info.git_commit : null,
      gitBranch: info.git_branch !== 'unknown' ? info.git_branch : null,
      buildTime: info.build_time ? new Date(info.build_time).toLocaleString('zh-CN') : null,
      platform: `${info.platform}/${info.arch}`,
      goVersion: info.go_version,
      nodeVersion: info.node_version
    }
  })

  // 获取完整版本信息
  const getFullVersionInfo = async () => {
    try {
      const response = await $fetch('/api/version/full') as any
      if (response.success) {
        return response.data.version_info
      }
    } catch (err) {
      console.error('获取完整版本信息失败:', err)
    }
    return null
  }

  return {
    versionInfo: readonly(versionInfo),
    loading: readonly(loading),
    error: readonly(error),
    formatVersionInfo,
    fetchVersionInfo,
    getVersionString,
    checkUpdate,
    getFullVersionInfo
  }
} 