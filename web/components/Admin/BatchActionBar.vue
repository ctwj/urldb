<template>
  <div
    v-if="selectedIds.length > 0"
    class="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-3 flex items-center justify-between"
  >
    <div class="text-sm text-blue-700 dark:text-blue-300">
      共 <span class="font-bold">{{ total }}</span> 项，已选
      <span class="font-bold">{{ selectedIds.length }}</span> 项
    </div>
    <div class="flex items-center gap-2">
      <template v-for="action in actions" :key="action.key">
        <n-popconfirm v-if="action.confirm" @positive-click="runAction(action)">
          <template #trigger>
            <n-button :type="action.type" size="small" :disabled="selectedIds.length === 0">
              <template #icon>
                <i :class="action.icon"></i>
              </template>
              {{ action.label }}
            </n-button>
          </template>
          <div class="max-w-xs">
            <p class="font-medium text-sm">{{ action.confirm.title }}</p>
            <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">{{ action.confirm.content }}</p>
          </div>
        </n-popconfirm>

        <n-button
          v-else
          :type="action.type"
          size="small"
          :disabled="selectedIds.length === 0"
          :loading="loadingKey === action.key"
          @click="runAction(action)"
        >
          <template #icon>
            <i :class="action.icon"></i>
          </template>
          {{ action.label }}
        </n-button>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useMessage } from 'naive-ui'

/**
 * 统一批量操作栏
 *
 * 显示"共 N 项已选 M 项"+ 按钮组。
 * BatchAction.confirm 存在时用 n-popconfirm 二次确认；
 * 否则直接执行。完成后用 useMessage 反馈，并 emit completed 通知父组件刷新。
 */

export interface BatchActionResult {
  success: boolean
  affected?: number
  message?: string
}

export interface BatchAction {
  key: string
  label: string
  type: 'primary' | 'info' | 'warning' | 'error'
  icon: string
  confirm?: { title: string; content: string }
  handler: (ids: (string | number)[]) => Promise<BatchActionResult>
}

const props = defineProps<{
  actions: BatchAction[]
  selectedIds: (string | number)[]
  total: number
}>()

const emit = defineEmits<{
  completed: [key: string, result: BatchActionResult]
}>()

const message = useMessage()
const loadingKey = ref<string | null>(null)

const runAction = async (action: BatchAction) => {
  if (props.selectedIds.length === 0) return
  loadingKey.value = action.key
  try {
    const result = await action.handler(props.selectedIds)
    if (result.success) {
      message.success(
        result.message ||
          `操作成功${result.affected ? `（影响 ${result.affected} 项）` : ''}`,
      )
    } else {
      message.error(result.message || '操作失败')
    }
    emit('completed', action.key, result)
  } catch (err) {
    const msg = (err as Error)?.message || '操作失败'
    message.error(msg)
    emit('completed', action.key, { success: false, message: msg })
  } finally {
    loadingKey.value = null
  }
}
</script>
