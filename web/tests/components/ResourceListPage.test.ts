/**
 * 列表页交互测试模板（T022b）
 *
 * 覆盖 resources.vue 及其余 7 个列表页（ready-resources/reports/copyright-claims/
 * accounts/files/failed-resources/untransferred-resources）共享的选择/批量/分页行为。
 *
 * 采用纯函数测试而非完整页面 mount：完整 mount resources.vue 需 mock 大量
 * Nuxt auto-import / Pinia / Naive UI provider，脆弱且偏离逻辑核心。
 * 选择逻辑提取为 utils/listSelection.ts 纯函数，跨页面复用且测试健壮。
 *
 * 其余 7 个列表页套用此模板时：只需 import 相同纯函数，无需重复编写。
 */
import { describe, it, expect } from 'vitest'
import {
  toggleItemSelection,
  toggleSelectAll,
  isAllSelected,
  isIndeterminate,
  clearOnPageChange,
} from '~/utils/listSelection'

describe('列表页选择逻辑（T022b 测试模板）', () => {
  const pageIds = [1, 2, 3] as const

  describe('(a) 全选/取消全选', () => {
    it('toggleSelectAll(true) 返回当前页所有 id', () => {
      const selected = toggleSelectAll([...pageIds], true)
      expect(selected).toEqual([1, 2, 3])
    })

    it('toggleSelectAll(false) 返回空数组', () => {
      const selected = toggleSelectAll([...pageIds], false)
      expect(selected).toEqual([])
    })

    it('全选后 isAllSelected 为 true', () => {
      const selected = toggleSelectAll([...pageIds], true)
      expect(isAllSelected(pageIds.length, selected.length)).toBe(true)
    })

    it('取消全选后 isAllSelected 为 false', () => {
      const selected = toggleSelectAll([...pageIds], false)
      expect(isAllSelected(pageIds.length, selected.length)).toBe(false)
    })

    it('部分选中时 isIndeterminate 为 true', () => {
      // 选中 1 项（共 3 项）
      const selected = toggleItemSelection([], 1, true)
      expect(isIndeterminate(pageIds.length, selected.length)).toBe(true)
      expect(isAllSelected(pageIds.length, selected.length)).toBe(false)
    })
  })

  describe('单项选择', () => {
    it('toggleItemSelection 追加未选中的 id', () => {
      const selected = toggleItemSelection([1, 2], 3, true)
      expect(selected).toEqual([1, 2, 3])
    })

    it('toggleItemSelection 不重复追加已选中的 id', () => {
      const selected = toggleItemSelection([1, 2], 2, true)
      expect(selected).toEqual([1, 2])
    })

    it('toggleItemSelection 移除已选中的 id', () => {
      const selected = toggleItemSelection([1, 2, 3], 2, false)
      expect(selected).toEqual([1, 3])
    })

    it('toggleItemSelection 对未选中 id 执行 false 无副作用', () => {
      const selected = toggleItemSelection([1, 2], 99, false)
      expect(selected).toEqual([1, 2])
    })

    it('不修改原数组（返回新数组）', () => {
      const original = [1, 2]
      const selected = toggleItemSelection(original, 3, true)
      expect(original).toEqual([1, 2])
      expect(selected).not.toBe(original)
    })
  })

  describe('(b)(c) 批量操作触发条件', () => {
    it('未选中时 BatchActionBar 应隐藏（selectedIds.length === 0）', () => {
      const selected = toggleSelectAll([...pageIds], false)
      expect(selected.length).toBe(0)
    })

    it('选中后 BatchActionBar 应显示（selectedIds.length > 0）', () => {
      const selected = toggleItemSelection([], 1, true)
      expect(selected.length).toBeGreaterThan(0)
    })

    it('批量操作 handler 接收当前 selectedIds', () => {
      // 模拟 BatchAction.handler 接收 ids
      const selected = toggleSelectAll([...pageIds], true)
      const handler = (ids: number[]) => ids.length
      expect(handler(selected as number[])).toBe(3)
    })
  })

  describe('(d) 切换分页后 selectedIds 清空', () => {
    it('clearOnPageChange 返回空数组', () => {
      const selected = toggleSelectAll([...pageIds], true)
      expect(selected.length).toBe(3)
      const after = clearOnPageChange()
      expect(after).toEqual([])
      expect(after.length).toBe(0)
    })
  })

  describe('边界值', () => {
    it('空页面 isAllSelected 为 false', () => {
      expect(isAllSelected(0, 0)).toBe(false)
    })

    it('空页面 isIndeterminate 为 false', () => {
      expect(isIndeterminate(0, 0)).toBe(false)
    })

    it('单条数据选中后 isAllSelected 为 true', () => {
      expect(isAllSelected(1, 1)).toBe(true)
      expect(isIndeterminate(1, 1)).toBe(false)
    })

    it('选中数大于页面总数时（异常防御）isAllSelected 为 false', () => {
      expect(isAllSelected(3, 5)).toBe(false)
    })
  })
})
