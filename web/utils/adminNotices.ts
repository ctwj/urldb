/**
 * 后台管理页统一提示文案常量
 *
 * 集中维护避免散落硬编码，便于未来国际化或文案调整。
 */

/**
 * 多标签页编辑提示
 *
 * 决策依据：spec Clarifications —— 多标签页采用 last-write-wins，不引入冲突检测。
 * 配置页（site-config / feature-config / dev-config / seo）顶部统一展示此提示。
 */
export const LAST_WRITE_WINS_NOTICE =
  '多标签页同时编辑时，后保存的版本将覆盖先保存的内容，请注意刷新查看最新值'
