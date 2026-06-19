<template>
  <ClientOnly>
    <Teleport to="body">
      <Transition name="cmd-fade">
        <div
          v-if="open"
          class="fixed inset-0 z-50 flex justify-center"
          role="dialog"
          aria-modal="true"
          aria-label="快速跳转"
        >
          <!-- 背景遮罩 -->
          <div
            class="absolute inset-0 bg-black/50 backdrop-blur-sm"
            @click="close"
          ></div>

          <!-- 面板容器 -->
          <div class="relative w-full max-w-xl mx-4 mt-[12vh] self-start">
            <div
              class="bg-white dark:bg-gray-800 rounded-xl shadow-2xl border border-gray-200 dark:border-gray-700 overflow-hidden"
            >
              <!-- 输入区 -->
              <div
                class="flex items-center gap-3 px-4 py-3 border-b border-gray-200 dark:border-gray-700"
              >
                <i class="fas fa-search text-gray-400 dark:text-gray-500 shrink-0"></i>
                <input
                  ref="inputRef"
                  v-model="query"
                  type="text"
                  placeholder="搜索功能页…"
                  aria-label="搜索功能页"
                  aria-autocomplete="list"
                  aria-controls="cmd-listbox"
                  :aria-activedescendant="activeDomId"
                  class="flex-1 bg-transparent outline-none text-gray-900 dark:text-white placeholder-gray-400 dark:placeholder-gray-500"
                  @keydown.down.prevent="moveDown"
                  @keydown.up.prevent="moveUp"
                  @keydown.enter.prevent="selectActive"
                  @keydown.esc.prevent="close"
                />
                <kbd
                  class="px-2 py-0.5 text-xs text-gray-500 dark:text-gray-400 bg-gray-100 dark:bg-gray-700 rounded border border-gray-200 dark:border-gray-600 font-sans"
                >ESC</kbd>
              </div>

              <!-- 结果区 -->
              <ul
                v-if="filtered.length > 0"
                id="cmd-listbox"
                ref="listRef"
                role="listbox"
                aria-label="功能页列表"
                class="max-h-[400px] overflow-y-auto py-2"
              >
                <template v-for="group in groupedFiltered" :key="group.key">
                  <li
                    v-if="groupedFiltered.length > 1"
                    role="presentation"
                    class="px-4 pt-3 pb-1 text-[11px] font-semibold uppercase tracking-wider text-gray-400 dark:text-gray-500"
                  >
                    {{ group.title }}
                  </li>
                  <li
                    v-for="item in group.items"
                    :id="domIdOf(item.id)"
                    :key="item.id"
                    role="option"
                    :aria-selected="item.id === activeId"
                  >
                    <NuxtLink
                      :to="item.to"
                      class="flex items-center gap-3 px-4 py-2.5 transition-colors cursor-pointer outline-none"
                      :class="{
                        'bg-blue-50 dark:bg-blue-900/30 text-blue-700 dark:text-blue-300': item.id === activeId,
                        'text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700/40': item.id !== activeId,
                      }"
                      @click="close"
                      @mouseenter="activeId = item.id"
                    >
                      <i
                        :class="item.icon"
                        class="w-5 text-center text-gray-400 dark:text-gray-500 shrink-0"
                      ></i>
                      <div class="flex-1 min-w-0">
                        <div
                          class="text-sm font-medium truncate"
                          v-html="highlight(item.title)"
                        ></div>
                        <div
                          class="text-xs text-gray-400 dark:text-gray-500 truncate font-mono"
                        >{{ item.to }}</div>
                      </div>
                      <i
                        v-if="item.id === activeId"
                        class="fas fa-arrow-right text-xs text-blue-500 dark:text-blue-400 shrink-0"
                      ></i>
                    </NuxtLink>
                  </li>
                </template>
              </ul>

              <!-- 空状态 -->
              <div v-else class="px-4 py-12 text-center">
                <i class="fas fa-search text-3xl text-gray-300 dark:text-gray-600 mb-3"></i>
                <p class="text-sm text-gray-500 dark:text-gray-400">未找到匹配的功能</p>
                <p class="text-xs text-gray-400 dark:text-gray-500 mt-1">尝试其他关键词</p>
              </div>

              <!-- 屏幕阅读器公告 -->
              <div aria-live="polite" class="sr-only">
                <span v-if="filtered.length > 0">{{ filtered.length }} 个结果</span>
                <span v-else>无匹配结果</span>
              </div>

              <!-- 底部状态栏 -->
              <div
                class="flex items-center justify-between px-4 py-2 border-t border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-900/30 text-xs text-gray-500 dark:text-gray-400"
              >
                <div class="flex items-center gap-3">
                  <span class="flex items-center gap-1">
                    <kbd
                      class="px-1.5 py-0.5 bg-white dark:bg-gray-700 rounded border border-gray-200 dark:border-gray-600 font-sans"
                    >↑</kbd>
                    <kbd
                      class="px-1.5 py-0.5 bg-white dark:bg-gray-700 rounded border border-gray-200 dark:border-gray-600 font-sans"
                    >↓</kbd>
                    <span>选择</span>
                  </span>
                  <span class="flex items-center gap-1">
                    <kbd
                      class="px-1.5 py-0.5 bg-white dark:bg-gray-700 rounded border border-gray-200 dark:border-gray-600 font-sans"
                    >↵</kbd>
                    <span>跳转</span>
                  </span>
                  <span class="flex items-center gap-1">
                    <kbd
                      class="px-1.5 py-0.5 bg-white dark:bg-gray-700 rounded border border-gray-200 dark:border-gray-600 font-sans"
                    >esc</kbd>
                    <span>关闭</span>
                  </span>
                </div>
                <span>{{ filtered.length }} 项</span>
              </div>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </ClientOnly>
</template>

<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue'
import { useMagicKeys, whenever } from '@vueuse/core'
import { useAdminNav, buildCommandPaletteItems } from '~/composables/useAdminNav'
import { filterItems } from '~/utils/commandPalette'
import type { CommandPaletteItem } from '~/composables/useAdminNav'

// 支持 v-model:open（与父组件如顶部"快速跳转"按钮双向联动）
const props = defineProps<{ open?: boolean }>()
const emit = defineEmits<{ 'update:open': [value: boolean] }>()

const { groups } = useAdminNav()
const allItems = computed<CommandPaletteItem[]>(() => buildCommandPaletteItems(groups))

// 内部状态与 prop 双向同步：父组件未传 open 时退化为内部自管
const open = computed({
  get: () => (props.open !== undefined ? props.open : internalOpen.value),
  set: (v) => {
    if (props.open !== undefined) emit('update:open', v)
    else internalOpen.value = v
  },
})
const internalOpen = ref(false)

const query = ref('')
// 用 id 而非 index 作为激活标识：filtered 顺序随 query 变化时不会错位
const activeId = ref<string | null>(null)
const inputRef = ref<HTMLInputElement | null>(null)
const listRef = ref<HTMLUListElement | null>(null)

const filtered = computed(() => filterItems(allItems.value, query.value))

// 分组策略：
// - query 为空：按 nav 原始分组聚合，配合分组标题便于浏览
// - query 非空：filterItems 已按相关性（精确 > 标题包含 > keywords）排序，单组展示不打乱顺序
const groupedFiltered = computed(() => {
  if (query.value.trim()) {
    return [{ key: 'results', title: '搜索结果', items: filtered.value }]
  }
  const map = new Map<string, CommandPaletteItem[]>()
  for (const item of filtered.value) {
    if (!map.has(item.group)) map.set(item.group, [])
    map.get(item.group)!.push(item)
  }
  const result: { key: string; title: string; items: CommandPaletteItem[] }[] = []
  for (const g of groups) {
    const items = map.get(g.title)
    if (items && items.length > 0) result.push({ key: g.key, title: g.title, items })
  }
  return result
})

// id → DOM id（aria-activedescendant 需要稳定 ID）
const domIdOf = (id: string) => `cmd-opt-${id.replace(/[^a-zA-Z0-9_-]/g, '-')}`
const activeDomId = computed(() => (activeId.value ? domIdOf(activeId.value) : undefined))

// 命中文字高亮：先 HTML 转义防 XSS，再用 RegExp 包裹 <mark>
function escapeHtml(s: string): string {
  return s.replace(/[&<>"']/g, (c) => {
    const m: Record<string, string> = { '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' }
    return m[c]
  })
}
function highlight(text: string): string {
  const q = query.value.trim()
  const safeText = escapeHtml(text)
  if (!q) return safeText
  const safeQ = escapeHtml(q)
  const re = new RegExp(safeQ.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'), 'gi')
  return safeText.replace(re, (m) => `<mark class="bg-yellow-200 dark:bg-yellow-500/40 text-inherit rounded-sm px-0.5">${m}</mark>`)
}

// Cmd/Ctrl+K 唤起/关闭
const keys = useMagicKeys()
whenever(keys['Cmd+K'], () => toggleOpen())
whenever(keys['Ctrl+K'], () => toggleOpen())

function toggleOpen() {
  open.value = !open.value
}
function close() {
  open.value = false
}

// 唤起时聚焦输入框并重置状态
watch(open, (v) => {
  if (v) {
    query.value = ''
    activeId.value = filtered.value[0]?.id ?? null
    nextTick(() => inputRef.value?.focus())
  }
})

// 输入变化时重置激活到首项
watch(query, () => {
  activeId.value = filtered.value[0]?.id ?? null
})

// 滚动激活项到可视区（仅键盘操作触发，避免 mouseenter 引起不必要滚动）
async function scrollActiveIntoView() {
  await nextTick()
  if (!activeId.value) return
  document.getElementById(domIdOf(activeId.value))?.scrollIntoView({ block: 'nearest' })
}

const moveDown = async () => {
  const list = filtered.value
  if (list.length === 0) return
  const idx = list.findIndex((i) => i.id === activeId.value)
  const next = idx < 0 ? 0 : (idx + 1) % list.length
  activeId.value = list[next].id
  await scrollActiveIntoView()
}

const moveUp = async () => {
  const list = filtered.value
  if (list.length === 0) return
  const idx = list.findIndex((i) => i.id === activeId.value)
  const prev = idx <= 0 ? list.length - 1 : idx - 1
  activeId.value = list[prev].id
  await scrollActiveIntoView()
}

const selectActive = () => {
  const item = filtered.value.find((i) => i.id === activeId.value)
  if (item) {
    close()
    navigateTo(item.to)
  }
}
</script>

<style scoped>
.cmd-fade-enter-active,
.cmd-fade-leave-active {
  transition: opacity 120ms ease;
}
.cmd-fade-enter-from,
.cmd-fade-leave-to {
  opacity: 0;
}
.fas {
  font-family: 'Font Awesome 6 Free';
  font-weight: 900;
}
/* 屏幕阅读器专用，视觉隐藏 */
.sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
  border: 0;
}
</style>
