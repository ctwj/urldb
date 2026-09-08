<template>
  <n-card :bordered="false" class="shadow-sm">
    <template #header>
      <div class="flex items-center gap-3">
        <div class="w-8 h-8 rounded-lg flex items-center justify-center flex-shrink-0"
          :class="verified ? 'bg-green-50 dark:bg-green-900/30' : 'bg-amber-50 dark:bg-amber-900/30'">
          <i class="fas fa-envelope-open-text text-sm"
            :class="verified ? 'text-green-500 dark:text-green-400' : 'text-amber-500 dark:text-amber-400'"></i>
        </div>
        <div>
          <p class="font-semibold text-gray-900 dark:text-white leading-tight">邮箱认证</p>
          <p class="text-xs text-gray-400 dark:text-gray-500">
            {{ verified ? '邮箱已完成认证' : '完成邮箱认证后才能申请 API 权限' }}
          </p>
        </div>
      </div>
    </template>
    <template #header-extra>
      <n-tag :type="verified ? 'success' : 'warning'" size="small" round>
        <template #icon>
          <i :class="verified ? 'fas fa-check' : 'fas fa-exclamation'"></i>
        </template>
        {{ verified ? '已认证' : '未认证' }}
      </n-tag>
    </template>

    <div v-if="verified" class="text-sm text-gray-500 dark:text-gray-400">
      邮箱 <span class="text-gray-800 dark:text-gray-200 font-medium">{{ email || '未绑定' }}</span> 已完成认证。
    </div>
    <div v-else>
      <p class="text-sm text-gray-500 dark:text-gray-400 mb-4">
        向 <span class="text-gray-800 dark:text-gray-200 font-medium">{{ email || '（未绑定邮箱）' }}</span> 发送 6 位验证码，输入后即完成认证。
      </p>
      <div class="flex flex-wrap items-center gap-3">
        <n-button type="primary" secondary :loading="sending" :disabled="!email || countdown > 0" @click="sendCode">
          <template #icon><i class="fas fa-paper-plane"></i></template>
          {{ countdown > 0 ? `${countdown}s 后可重发` : '发送验证码' }}
        </n-button>
        <n-input v-model:value="code" type="text" placeholder="6 位验证码" maxlength="6"
          style="width: 150px" @keyup.enter="verify" />
        <n-button type="primary" :loading="verifying" :disabled="code.length === 0" @click="verify">
          验证
        </n-button>
      </div>
    </div>
  </n-card>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useNotification } from 'naive-ui'
import { useApiFetch } from '~/composables/useApiFetch'
import { useUserStore } from '~/stores/user'

// 邮箱认证卡片（016-api-access-application）
// API 访问页与设置页共用：展示认证状态，未认证时提供发码/验码流程
const props = defineProps<{
  verified: boolean
}>()

const emit = defineEmits<{
  (e: 'verified'): void
}>()

const notification = useNotification()
const userStore = useUserStore()

const sending = ref(false)
const verifying = ref(false)
const code = ref('')
const countdown = ref(0)
let timer: ReturnType<typeof setInterval> | null = null

const email = computed(() => userStore.user?.email || '')

const notify = (type: 'success' | 'error' | 'warning' | 'info', title: string, content: string) => {
  notification[type]({ title, content, duration: 4000 })
}

const startCountdown = () => {
  countdown.value = 60
  if (timer) clearInterval(timer)
  timer = setInterval(() => {
    countdown.value--
    if (countdown.value <= 0 && timer) {
      clearInterval(timer)
      timer = null
    }
  }, 1000)
}

const sendCode = async () => {
  sending.value = true
  try {
    await useApiFetch('/user/email/verification-code', { method: 'POST' })
    notify('success', '验证码已发送', '请查收邮箱（注意垃圾邮件夹）')
    startCountdown()
  } catch (e: any) {
    notify('error', '发送失败', e?.message || '请稍后重试')
  } finally {
    sending.value = false
  }
}

const verify = async () => {
  verifying.value = true
  try {
    await useApiFetch('/user/email/verify', { method: 'POST', body: { code: code.value } })
    notify('success', '邮箱认证成功', '现在可以使用需要邮箱认证的功能了')
    code.value = ''
    emit('verified')
  } catch (e: any) {
    notify('error', '验证失败', e?.message || '验证码错误或已过期')
  } finally {
    verifying.value = false
  }
}

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>
