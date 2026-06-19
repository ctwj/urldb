import { describe, it, expect } from 'vitest'
import { filterItems } from '~/utils/commandPalette'
import type { CommandPaletteItem } from '~/composables/useAdminNav'

// 构造测试数据：覆盖多分组顺序、keywords、相同字符串不同位置
const items: CommandPaletteItem[] = [
  { id: '/admin', title: '仪表盘', group: '仪表盘', icon: 'fas fa-tachometer-alt', to: '/admin', keywords: ['仪表盘', 'dashboard'] },
  { id: '/admin/resources', title: '资源管理', group: '数据管理', icon: 'fas fa-database', to: '/admin/resources', keywords: ['资源管理', '数据管理', 'resources'] },
  { id: '/admin/copyright-claims', title: '版权申述', group: '数据管理', icon: 'fas fa-balance-scale', to: '/admin/copyright-claims', keywords: ['版权申述', '数据管理', 'copyright'] },
  { id: '/admin/site-config', title: '站点配置', group: '系统配置', icon: 'fas fa-globe', to: '/admin/site-config', keywords: ['站点配置', 'site-config'] },
  { id: '/admin/seo', title: 'SEO', group: '运营管理', icon: 'fas fa-search', to: '/admin/seo', keywords: ['SEO', '搜索引擎优化'] },
]

describe('filterItems', () => {
  describe('空 query', () => {
    it('空字符串返回全部按分组顺序', () => {
      const result = filterItems(items, '')
      expect(result).toHaveLength(items.length)
      expect(result.map((i) => i.id)).toEqual(items.map((i) => i.id))
    })

    it('纯空白字符视为空 query', () => {
      const result = filterItems(items, '   ')
      expect(result).toHaveLength(items.length)
    })
  })

  describe('标题匹配', () => {
    it('标题精确匹配排在最前', () => {
      const result = filterItems(items, 'SEO')
      expect(result[0].id).toBe('/admin/seo')
    })

    it('标题包含匹配（大小写不敏感）', () => {
      const result = filterItems(items, '资源')
      expect(result.map((i) => i.id)).toContain('/admin/resources')
      expect(result[0].id).toBe('/admin/resources')
    })

    it('英文大小写不敏感', () => {
      const result = filterItems(items, 'seo')
      expect(result[0].id).toBe('/admin/seo')
    })
  })

  describe('keywords 匹配', () => {
    it('keywords 包含匹配也能命中', () => {
      // "数据管理" 仅出现在 resources 的 keywords 中（不作为任何 title）
      const result = filterItems(items, '数据管理')
      // 多个项的 keywords 含 "数据管理"：resources 与 copyright-claims 同组
      expect(result.map((i) => i.id)).toContain('/admin/resources')
    })

    it('标题精确匹配优先于 keywords 匹配', () => {
      // 构造一项：A 标题精确匹配、B 仅 keywords 匹配
      const local: CommandPaletteItem[] = [
        { id: 'b', title: '其他', group: 'X', icon: 'i', to: '/b', keywords: ['target'] },
        { id: 'a', title: 'target', group: 'X', icon: 'i', to: '/a', keywords: [] },
      ]
      const result = filterItems(local, 'target')
      expect(result[0].id).toBe('a')
    })
  })

  describe('无匹配', () => {
    it('返回空数组', () => {
      const result = filterItems(items, '不存在的关键字xyz')
      expect(result).toEqual([])
    })
  })

  describe('排序稳定性', () => {
    it('同优先级的命中按原数组顺序返回', () => {
      // 两个均仅 keywords 命中 "数据管理"
      const result = filterItems(items, '数据管理')
      const ids = result.map((i) => i.id)
      const resourcesIdx = ids.indexOf('/admin/resources')
      const copyrightIdx = ids.indexOf('/admin/copyright-claims')
      // 两项都应命中
      expect(resourcesIdx).toBeGreaterThanOrEqual(0)
      expect(copyrightIdx).toBeGreaterThanOrEqual(0)
      // resources 在 copyright-claims 之前（原顺序）
      expect(resourcesIdx).toBeLessThan(copyrightIdx)
    })
  })
})
