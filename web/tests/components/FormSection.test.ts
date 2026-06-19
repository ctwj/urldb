/**
 * FormSection 渲染测试（T033 测试模板）
 *
 * 验证统一表单分组组件的核心渲染契约：
 * - title/description 默认渲染
 * - header slot 覆盖 title
 * - actions/footer slot 渲染
 * - 无 title 时不渲染 header 节点
 *
 * 配置页（site-config/feature-config/dev-config/seo）后续套用时复用此组件。
 */
import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import FormSection from '~/components/Admin/FormSection.vue'

const stub = { template: '<div />' }

describe('FormSection 渲染（T033 模板）', () => {
  it('渲染 title 与 description', () => {
    const wrapper = mount(FormSection, {
      props: { title: '基础配置', description: '站点核心信息' },
    })
    expect(wrapper.text()).toContain('基础配置')
    expect(wrapper.text()).toContain('站点核心信息')
  })

  it('未传 title 时不渲染 header 边界', () => {
    const wrapper = mount(FormSection, {})
    expect(wrapper.find('header').exists()).toBe(false)
  })

  it('header slot 覆盖默认 title', () => {
    const wrapper = mount(FormSection, {
      props: { title: '默认标题' },
      slots: { header: '<span class="custom">自定义标题</span>' },
    })
    expect(wrapper.find('.custom').exists()).toBe(true)
    expect(wrapper.text()).not.toContain('默认标题')
  })

  it('actions slot 渲染到 header 右侧', () => {
    const wrapper = mount(FormSection, {
      props: { title: '配置组' },
      slots: { actions: '<button class="reset">重置</button>' },
    })
    expect(wrapper.find('header button.reset').exists()).toBe(true)
  })

  it('footer slot 渲染到底部', () => {
    const wrapper = mount(FormSection, {
      props: { title: '配置组' },
      slots: { footer: '<button class="save">保存</button>' },
    })
    expect(wrapper.find('footer button.save').exists()).toBe(true)
  })

  it('无 footer slot 时不渲染 footer 节点', () => {
    const wrapper = mount(FormSection, {
      props: { title: '配置组' },
    })
    expect(wrapper.find('footer').exists()).toBe(false)
  })

  it('default slot 渲染到主体', () => {
    const wrapper = mount(FormSection, {
      props: { title: '配置组' },
      slots: { default: '<div class="field">字段</div>' },
      global: { stubs: { NButton: stub } },
    })
    expect(wrapper.find('.field').exists()).toBe(true)
  })
})
