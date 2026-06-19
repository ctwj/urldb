/**
 * 列表页选择逻辑纯函数
 *
 * 跨 8 个列表页（resources/reports/...）共享的选择/批量/分页行为。
 * 提取为纯函数便于单元测试（T022b 测试模板）与跨页面复用。
 */

export type SelectableId = string | number

/**
 * 切换单项选择
 * - checked=true 且 id 不在已选中：追加
 * - checked=false 且 id 在已选中：移除
 * - 返回新数组（不修改原数组）
 */
export function toggleItemSelection(
  selected: SelectableId[],
  id: SelectableId,
  checked: boolean,
): SelectableId[] {
  if (checked) {
    return selected.includes(id) ? selected : [...selected, id]
  }
  return selected.filter((x) => x !== id)
}

/**
 * 全选/取消全选当前页
 * - checked=true：返回当前页所有 id
 * - checked=false：返回空数组
 */
export function toggleSelectAll(
  allIds: SelectableId[],
  checked: boolean,
): SelectableId[] {
  return checked ? [...allIds] : []
}

/**
 * 是否全部选中（当前页非空且选中数等于当前页总数）
 */
export function isAllSelected(
  pageItemCount: number,
  selectedCount: number,
): boolean {
  return pageItemCount > 0 && selectedCount === pageItemCount
}

/**
 * 是否部分选中（选中数大于 0 且小于当前页总数）
 */
export function isIndeterminate(
  pageItemCount: number,
  selectedCount: number,
): boolean {
  return selectedCount > 0 && selectedCount < pageItemCount
}

/**
 * 分页切换时清空已选（跨页选择在当前架构下不保留）
 */
export function clearOnPageChange(): SelectableId[] {
  return []
}
