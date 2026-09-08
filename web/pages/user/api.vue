<template>
  <div class="max-w-4xl mx-auto space-y-6">
    <!-- 页头 -->
    <div>
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white flex items-center">
        <i class="fas fa-key text-blue-500 mr-2"></i>
        API 访问
      </h1>
      <p class="text-gray-600 dark:text-gray-400 mt-1">申请开通 API 权限，通过程序化方式查询平台公开资源</p>
    </div>

    <n-spin :show="loading">
      <!-- 修复：n-spin 会包裹一层容器，内部需要自己的纵向间距 -->
      <div class="space-y-6">
        <!-- 邮箱认证（未认证时显示） -->
        <UserEmailVerificationCard v-if="status && !status.email_verified" :verified="false" @verified="fetchStatus" />

        <!-- 凭证卡（已开通） -->
        <div v-if="credential && credential.exists" class="bg-white dark:bg-gray-800 rounded-xl shadow-sm border border-gray-200 dark:border-gray-700 p-6">
          <div class="flex items-center justify-between mb-4">
            <div class="flex items-center">
              <div class="w-1 h-6 bg-green-500 rounded-full mr-3"></div>
              <h2 class="text-lg font-semibold text-gray-800 dark:text-gray-200">API 凭证</h2>
            </div>
            <n-tag :type="credentialTagType" size="small">{{ credentialTagLabel }}</n-tag>
          </div>

          <!-- 重置后：明文凭证仍在本次会话中，可复制 / 重新查看 -->
          <n-alert v-if="newKey" type="warning" :show-icon="true" class="mb-4">
            <div class="flex flex-wrap items-center gap-x-3 gap-y-2">
              <span class="text-sm">新凭证已生成，明文仅本次访问可复制，刷新页面后需重置获取。</span>
              <n-button size="tiny" type="warning" @click="copyNewKey">
                <template #icon><i :class="copiedRecently ? 'fas fa-check' : 'fas fa-copy'"></i></template>
                {{ copiedRecently ? '已复制' : '复制凭证' }}
              </n-button>
              <n-button size="tiny" quaternary @click="showKeyModal = true">重新查看</n-button>
            </div>
          </n-alert>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 mb-4">
            <div>
              <p class="text-xs text-gray-400 dark:text-gray-500 mb-1">凭证标识（前缀）</p>
              <code class="text-sm bg-gray-100 dark:bg-gray-700 rounded px-2 py-1 text-gray-800 dark:text-gray-200">{{ credential.key_prefix }}...</code>
            </div>
            <div>
              <p class="text-xs text-gray-400 dark:text-gray-500 mb-1">到期时间</p>
              <p class="text-sm text-gray-800 dark:text-gray-200">
                {{ credential.expires_at ? formatTime(credential.expires_at) : '永久有效' }}
                <n-tag v-if="credential.expired" type="error" size="tiny" class="ml-1">已过期</n-tag>
              </p>
            </div>
            <div>
              <p class="text-xs text-gray-400 dark:text-gray-500 mb-1">最近调用</p>
              <p class="text-sm text-gray-800 dark:text-gray-200">{{ credential.last_used_at ? formatTime(credential.last_used_at) : '从未调用' }}</p>
            </div>
            <div>
              <p class="text-xs text-gray-400 dark:text-gray-500 mb-1">创建时间</p>
              <p class="text-sm text-gray-800 dark:text-gray-200">{{ formatTime(credential.created_at) }}</p>
            </div>
          </div>

          <div class="flex items-center gap-3">
            <n-button type="warning" secondary :disabled="credential.status !== 'active' || credential.expired" :loading="resetting" @click="showResetConfirm = true">
              <template #icon><i class="fas fa-sync-alt"></i></template>
              重置凭证
            </n-button>
            <span class="text-xs text-gray-400 dark:text-gray-500">重置后旧凭证立即失效，到期时间不变</span>
          </div>
        </div>

        <!-- API 使用说明 -->
        <div class="bg-white dark:bg-gray-800 rounded-xl shadow-sm border border-gray-200 dark:border-gray-700 p-6">
          <div class="flex items-center justify-between mb-4">
            <div class="flex items-center">
              <div class="w-1 h-6 bg-emerald-500 rounded-full mr-3"></div>
              <h2 class="text-lg font-semibold text-gray-800 dark:text-gray-200">API 使用说明</h2>
            </div>
            <NuxtLink to="/api-docs" class="text-sm text-blue-500 hover:text-blue-600 dark:text-blue-400">
              查看完整 API 文档 <i class="fas fa-arrow-right text-xs ml-0.5"></i>
            </NuxtLink>
          </div>

          <div class="space-y-5 text-sm text-gray-700 dark:text-gray-300">
            <!-- 端点 -->
            <div class="flex flex-wrap items-center gap-2">
              <span class="px-2 py-0.5 rounded bg-emerald-100 dark:bg-emerald-900/40 text-emerald-700 dark:text-emerald-300 text-xs font-bold">GET</span>
              <code class="text-sm bg-gray-100 dark:bg-gray-700 rounded px-2 py-1 text-gray-800 dark:text-gray-200">/api/open/resources/search</code>
              <span class="text-xs text-gray-400 dark:text-gray-500">检索公开且有效的资源</span>
            </div>

            <!-- 认证 -->
            <p>
              认证：在请求头携带个人凭证
              <code class="bg-gray-100 dark:bg-gray-700 rounded px-1.5 py-0.5 text-xs">X-API-Key: urldb_xxx</code>
              （也支持 <code class="bg-gray-100 dark:bg-gray-700 rounded px-1.5 py-0.5 text-xs">?api_key=</code> 查询参数）。
              凭证仅在重置时展示一次，请妥善保存。
            </p>

            <!-- 参数 -->
            <div>
              <p class="font-medium text-gray-800 dark:text-gray-200 mb-2">请求参数</p>
              <div class="overflow-x-auto">
                <table class="w-full text-sm">
                  <thead>
                    <tr class="text-left text-xs text-gray-400 dark:text-gray-500 border-b border-gray-200 dark:border-gray-700">
                      <th class="py-2 pr-4 font-medium">参数</th>
                      <th class="py-2 pr-4 font-medium">必填</th>
                      <th class="py-2 font-medium">说明</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr class="border-b border-gray-100 dark:border-gray-700/60">
                      <td class="py-2 pr-4"><code class="text-xs bg-gray-100 dark:bg-gray-700 rounded px-1.5 py-0.5">keyword</code></td>
                      <td class="py-2 pr-4">是</td>
                      <td class="py-2">搜索关键词，不超过 100 字</td>
                    </tr>
                    <tr>
                      <td class="py-2 pr-4"><code class="text-xs bg-gray-100 dark:bg-gray-700 rounded px-1.5 py-0.5">size</code></td>
                      <td class="py-2 pr-4">否</td>
                      <td class="py-2">返回数量，默认 20，最大 100（超过按 100 处理）；接口无翻页，固定返回命中的前 size 条</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>

            <!-- 请求示例 -->
            <div>
              <p class="font-medium text-gray-800 dark:text-gray-200 mb-2">请求示例</p>
              <div class="relative group rounded-lg bg-gray-900 dark:bg-black/60 overflow-hidden">
                <button
                  class="absolute top-2 right-2 px-2 py-1 rounded text-xs bg-gray-700/80 text-gray-200 hover:bg-gray-600 transition-colors"
                  @click="copyCurl"
                >
                  <i :class="copiedCurl ? 'fas fa-check' : 'fas fa-copy'" class="mr-1"></i>
                  {{ copiedCurl ? '已复制' : '复制' }}
                </button>
                <pre class="p-4 text-xs leading-relaxed text-gray-100 overflow-x-auto"><code>{{ curlExample }}</code></pre>
              </div>
            </div>

            <!-- 响应结构 -->
            <div>
              <p class="font-medium text-gray-800 dark:text-gray-200 mb-2">响应结构</p>
              <div class="rounded-lg bg-gray-50 dark:bg-gray-900/60 border border-gray-100 dark:border-gray-700/60 p-4 overflow-x-auto">
                <pre class="text-xs leading-relaxed text-gray-700 dark:text-gray-300"><code>{
  "success": true,
  "data": {
    "list": [{ "title", "description", "url",
               "key", "cover", "tags", "created_at" }],
    "size": 20
  }
}</code></pre>
              </div>
              <p class="text-xs text-gray-500 dark:text-gray-400 -mt-2">
                字段说明：<code class="text-xs">url</code> 优先返回转存链接（无转存时为原始链接）；<code class="text-xs">cover</code> 为封面图片；不返回 total / id / save_url / category / pan_name / file_size，以 <code class="text-xs">key</code> 标识资源。
              </p>
            </div>

            <!-- 限频与错误码 -->
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div class="rounded-lg border border-gray-100 dark:border-gray-700/60 p-4">
                <p class="font-medium text-gray-800 dark:text-gray-200 mb-2">
                  <i class="fas fa-tachometer-alt text-amber-500 mr-1.5"></i>频率限制
                </p>
                <p class="text-xs leading-relaxed">每分钟、每小时、每天三档配额同时生效（以平台配置为准）。超限返回 <code class="text-xs">429</code>，响应消息中会标明具体上限，请稍后重试。</p>
              </div>
              <div class="rounded-lg border border-gray-100 dark:border-gray-700/60 p-4">
                <p class="font-medium text-gray-800 dark:text-gray-200 mb-2">
                  <i class="fas fa-circle-exclamation text-red-500 mr-1.5"></i>错误码
                </p>
                <ul class="text-xs space-y-1.5">
                  <li><code class="text-xs">400</code> 参数错误（keyword 为空或超长等）</li>
                  <li><code class="text-xs">401</code> 凭证无效或缺失</li>
                  <li><code class="text-xs">403</code> 权限已停用 / 凭证已过期</li>
                  <li><code class="text-xs">429</code> 请求过于频繁</li>
                </ul>
              </div>
            </div>
          </div>
        </div>

        <!-- 申请状态卡 -->
        <div v-if="application" class="bg-white dark:bg-gray-800 rounded-xl shadow-sm border border-gray-200 dark:border-gray-700 p-6">
          <div class="flex items-center mb-4">
            <div class="w-1 h-6 bg-blue-500 rounded-full mr-3"></div>
            <h2 class="text-lg font-semibold text-gray-800 dark:text-gray-200">我的申请</h2>
          </div>
          <div class="space-y-3">
            <div class="flex items-center gap-3">
              <n-tag :type="applicationTagType">{{ applicationTagLabel }}</n-tag>
              <span class="text-sm text-gray-500 dark:text-gray-400">提交于 {{ formatTime(application.created_at) }}</span>
            </div>
            <div>
              <p class="text-xs text-gray-400 dark:text-gray-500 mb-1">用途说明</p>
              <p class="text-sm text-gray-800 dark:text-gray-200 whitespace-pre-wrap">{{ application.purpose }}</p>
            </div>
            <n-alert v-if="application.status === 'rejected'" type="error" :show-icon="true" class="mt-2">
              申请被拒绝{{ application.reject_reason ? `：${application.reject_reason}` : '' }}。您可以修改用途后重新提交。
            </n-alert>
            <n-alert v-if="application.status === 'pending'" type="info" :show-icon="true">
              申请正在等待管理员审核，审核结果将通过本页面展示。
            </n-alert>
          </div>
        </div>

        <!-- 申请表单（无申请 / 已拒绝 / 已过期时可提交） -->
        <div v-if="canApply" class="bg-white dark:bg-gray-800 rounded-xl shadow-sm border border-gray-200 dark:border-gray-700 p-6">
          <div class="flex items-center mb-4">
            <div class="w-1 h-6 bg-purple-500 rounded-full mr-3"></div>
            <h2 class="text-lg font-semibold text-gray-800 dark:text-gray-200">{{ application && application.status === 'rejected' ? '重新提交申请' : '申请开通 API' }}</h2>
          </div>
          <p class="text-sm text-gray-500 dark:text-gray-400 mb-4">
            请填写 API 使用用途说明（必填，不超过 500 字），管理员审核通过后将自动为您开通。
          </p>
          <n-input
            v-model:value="purpose"
            type="textarea"
            placeholder="例如：接入个人聚合搜索站，展示本站公开资源"
            :maxlength="500"
            show-count
            :rows="4"
          />
          <div class="mt-4">
            <n-button type="primary" :loading="applying" :disabled="!purpose.trim() || (status && !status.email_verified)" @click="submitApply">
              提交申请
            </n-button>
            <span v-if="status && !status.email_verified" class="text-xs text-amber-500 ml-3">完成邮箱认证后才能提交</span>
          </div>
        </div>

        <!-- 空状态 -->
        <div v-if="!loading && !status" class="text-center py-16 text-gray-400">
          <i class="fas fa-plug text-3xl mb-3"></i>
          <p>加载失败，请刷新重试</p>
        </div>
      </div>
    </n-spin>

    <!-- 重置确认 -->
    <n-modal v-model:show="showResetConfirm" preset="dialog" type="warning" title="重置 API 凭证"
      content="重置后旧凭证立即失效，所有使用旧凭证的程序将无法继续调用。确定重置吗？"
      positive-text="确定重置" negative-text="取消"
      :loading="resetting"
      @positive-click="confirmReset" />

    <!-- 一次性明文展示 -->
    <n-modal v-model:show="showKeyModal" preset="card" title="您的新 API 凭证" style="max-width: 560px">
      <n-alert type="warning" :show-icon="true" class="mb-4">
        完整凭证仅此一次展示，请立即复制并妥善保存。刷新页面后将无法再次查看。
      </n-alert>
      <div class="relative group">
        <button
          class="absolute top-2 right-2 px-2 py-1 rounded text-xs bg-gray-200/80 dark:bg-gray-600/80 text-gray-700 dark:text-gray-100 hover:bg-gray-300 dark:hover:bg-gray-500 transition-colors"
          @click="copyNewKey"
        >
          <i :class="copiedRecently ? 'fas fa-check' : 'fas fa-copy'" class="mr-1"></i>
          {{ copiedRecently ? '已复制' : '复制' }}
        </button>
        <code class="block text-sm bg-gray-100 dark:bg-gray-700 rounded p-3 pr-16 break-all text-gray-800 dark:text-gray-200">{{ newKey }}</code>
      </div>
      <template #footer>
        <div class="flex justify-end gap-2">
          <n-button type="primary" @click="showKeyModal = false">我已保存</n-button>
        </div>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useDialog, useNotification } from 'naive-ui'
import { useApiFetch } from '~/composables/useApiFetch'
import { useTimeFormat } from '~/composables/useTimeFormat'

definePageMeta({
  layout: 'user',
  title: 'API 访问'
})

const notification = useNotification()
const dialog = useDialog()
const { formatDateTime: formatTime } = useTimeFormat()

const loading = ref(true)
const status = ref<any>(null)

// 申请
const purpose = ref('')
const applying = ref(false)

// 凭证重置
const resetting = ref(false)
const showResetConfirm = ref(false)
const showKeyModal = ref(false)
const newKey = ref('')

// 复制成功态（短暂显示"已复制"）
const copiedRecently = ref(false)
const copiedCurl = ref(false)
let copyTimer: ReturnType<typeof setTimeout> | null = null
let curlTimer: ReturnType<typeof setTimeout> | null = null

const application = computed(() => status.value?.application || null)
const credential = computed(() => (status.value?.credential?.exists ? status.value.credential : null))

// 请求示例（挂载后取当前站点域名避免 SSR hydration 不一致；凭证以占位符展示——明文只在重置时可见）
const siteOrigin = ref('')
const curlExample = computed(() => {
  const origin = siteOrigin.value || 'https://your-domain.com'
  return `curl "${origin}/api/open/resources/search?keyword=电影&size=20" \\\n  -H "X-API-Key: urldb_你的凭证"`
})

// 可提交申请：未认证不行；有待审核不行；已开通（未过期）不行；无申请/已拒绝/已过期可以
const canApply = computed(() => {
  if (!status.value || !status.value.email_verified) return false
  if (application.value && application.value.status === 'pending') return false
  if (credential.value && credential.value.status === 'active' && !credential.value.expired) return false
  return true
})

const applicationTagType = computed(() => {
  switch (application.value?.status) {
    case 'pending': return 'warning'
    case 'approved': return 'success'
    case 'rejected': return 'error'
    default: return 'default'
  }
})
const applicationTagLabel = computed(() => {
  switch (application.value?.status) {
    case 'pending': return '待审核'
    case 'approved': return '已同意'
    case 'rejected': return '已拒绝'
    default: return '未知'
  }
})
const credentialTagType = computed(() => {
  if (!credential.value) return 'default'
  if (credential.value.expired) return 'error'
  return credential.value.status === 'active' ? 'success' : 'error'
})
const credentialTagLabel = computed(() => {
  if (!credential.value) return ''
  if (credential.value.expired) return '已过期'
  return credential.value.status === 'active' ? '有效' : '已停用'
})

const notify = (type: 'success' | 'error' | 'warning' | 'info', title: string, content: string) => {
  notification[type]({ title, content, duration: 4000 })
}

const fetchStatus = async () => {
  loading.value = true
  try {
    status.value = await useApiFetch('/user/api/status').then((r: any) => r?.data ?? r)
  } catch (e: any) {
    console.error('获取 API 状态失败:', e)
  } finally {
    loading.value = false
  }
}

const submitApply = async () => {
  applying.value = true
  try {
    await useApiFetch('/user/api/apply', { method: 'POST', body: { purpose: purpose.value.trim() } })
    notify('success', '申请已提交', '请等待管理员审核')
    purpose.value = ''
    await fetchStatus()
  } catch (e: any) {
    notify('error', '提交失败', e?.message || '请稍后重试')
  } finally {
    applying.value = false
  }
}

const confirmReset = async () => {
  resetting.value = true
  try {
    const res: any = await useApiFetch('/user/api/reset-key', { method: 'POST' })
    newKey.value = res?.data?.key || ''
    showKeyModal.value = true
    showResetConfirm.value = false
    await fetchStatus()
  } catch (e: any) {
    notify('error', '重置失败', e?.message || '请稍后重试')
  } finally {
    resetting.value = false
  }
}

// 写入剪贴板：优先 Clipboard API（需 secure context），非 https 部署时回退 execCommand
const copyText = async (text: string): Promise<boolean> => {
  if (!text) return false
  try {
    if (navigator.clipboard && window.isSecureContext) {
      await navigator.clipboard.writeText(text)
      return true
    }
  } catch { /* 回退到 execCommand */ }
  try {
    const ta = document.createElement('textarea')
    ta.value = text
    ta.style.position = 'fixed'
    ta.style.opacity = '0'
    document.body.appendChild(ta)
    ta.select()
    const ok = document.execCommand('copy')
    document.body.removeChild(ta)
    return ok
  } catch {
    return false
  }
}

const flashCopied = (which: 'key' | 'curl') => {
  if (which === 'key') {
    copiedRecently.value = true
    if (copyTimer) clearTimeout(copyTimer)
    copyTimer = setTimeout(() => (copiedRecently.value = false), 1500)
  } else {
    copiedCurl.value = true
    if (curlTimer) clearTimeout(curlTimer)
    curlTimer = setTimeout(() => (copiedCurl.value = false), 1500)
  }
}

const copyNewKey = async () => {
  if (await copyText(newKey.value)) {
    flashCopied('key')
    notify('success', '已复制', '凭证已复制到剪贴板')
  } else {
    dialog.warning({ title: '复制失败', content: '请手动选中凭证文本复制', positiveText: '知道了' })
  }
}

const copyCurl = async () => {
  if (await copyText(curlExample.value)) {
    flashCopied('curl')
  }
}

onMounted(() => {
  siteOrigin.value = window.location.origin
  fetchStatus()
})

onUnmounted(() => {
  if (copyTimer) clearTimeout(copyTimer)
  if (curlTimer) clearTimeout(curlTimer)
})
</script>
