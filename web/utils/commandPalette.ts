import type { CommandPaletteItem } from '~/composables/useAdminNav'

/**
 * 命令面板过滤纯函数
 *
 * 优先级（高 → 低）：
 *   1. 标题完全匹配（大小写不敏感）
 *   2. 标题包含匹配
 *   3. keywords 包含匹配
 *
 * 同优先级内保持原数组顺序（即分组顺序）。
 * 空 query 返回全部（拷贝以避免外部突变）。
 */
export function filterItems(
  items: CommandPaletteItem[],
  query: string,
): CommandPaletteItem[] {
  const q = (query ?? '').trim().toLowerCase()
  if (!q) return items.slice()

  const exact: CommandPaletteItem[] = []
  const titleContains: CommandPaletteItem[] = []
  const keywordContains: CommandPaletteItem[] = []

  for (const item of items) {
    const titleLower = item.title.toLowerCase()
    if (titleLower === q) {
      exact.push(item)
      continue
    }
    if (titleLower.includes(q)) {
      titleContains.push(item)
      continue
    }
    if (item.keywords?.some((k) => k.toLowerCase().includes(q))) {
      keywordContains.push(item)
    }
  }

  return [...exact, ...titleContains, ...keywordContains]
}
