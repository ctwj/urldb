import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import FilterBar from '~/components/Admin/FilterBar.vue'
import type { FilterConfig } from '~/components/Admin/FilterBar.vue'

const config: FilterConfig = {
  search: { placeholder: '搜索资源...', key: 'search' },
  selects: [
    { key: 'category', placeholder: '选择分类', options: [{ label: '分类A', value: 1 }] },
    { key: 'platform', placeholder: '选择平台', options: [{ label: '平台X', value: 'x' }] },
  ],
}

describe('FilterBar', () => {
  it('根据 config 正确渲染搜索框与下拉筛选', () => {
    const wrapper = mount(FilterBar, {
      props: { config, modelValue: { search: '', category: null, platform: null } },
      global: { stubs: { NInput: true, NSelect: true, NButton: true } },
    })

    // 搜索框存在
    const inputs = wrapper.findAllComponents({ name: 'NInput' })
    expect(inputs.length).toBeGreaterThanOrEqual(1)

    // 两个下拉
    const selects = wrapper.findAllComponents({ name: 'NSelect' })
    expect(selects).toHaveLength(2)
  })

  it('config 无 search 时仅渲染下拉', () => {
    const noSearchConfig: FilterConfig = {
      selects: [{ key: 'status', placeholder: '状态', options: [] }],
    }
    const wrapper = mount(FilterBar, {
      props: { config: noSearchConfig, modelValue: { status: null } },
      global: { stubs: { NInput: true, NSelect: true, NButton: true } },
    })
    expect(wrapper.findAllComponents({ name: 'NSelect' })).toHaveLength(1)
  })

  it('点击搜索按钮触发 search 事件', async () => {
    const wrapper = mount(FilterBar, {
      props: { config, modelValue: { search: '关键词', category: null, platform: null } },
      global: { stubs: { NInput: true, NSelect: true, NButton: { template: '<button @click="$emit(\'click\')"><slot /></button>' } } },
    })

    // 找到搜索按钮（第一个 NButton 是搜索，第二个是重置）
    const buttons = wrapper.findAll('button')
    const searchBtn = buttons[0]
    await searchBtn.trigger('click')
    expect(wrapper.emitted('search')).toBeTruthy()
  })

  it('点击重置按钮触发 reset 事件并清空值', async () => {
    const wrapper = mount(FilterBar, {
      props: { config, modelValue: { search: '关键词', category: 1, platform: 'x' } },
      global: { stubs: { NInput: true, NSelect: true, NButton: { template: '<button @click="$emit(\'click\')"><slot /></button>' } } },
    })

    const buttons = wrapper.findAll('button')
    const resetBtn = buttons[1]
    await resetBtn.trigger('click')

    expect(wrapper.emitted('reset')).toBeTruthy()
    // reset 时应同时 emit update:modelValue 为空状态
    const updateEvents = wrapper.emitted('update:modelValue')
    expect(updateEvents).toBeTruthy()
    const lastUpdate = updateEvents![updateEvents!.length - 1][0] as Record<string, any>
    expect(lastUpdate.search).toBe('')
    expect(lastUpdate.category).toBeNull()
    expect(lastUpdate.platform).toBeNull()
  })
})
