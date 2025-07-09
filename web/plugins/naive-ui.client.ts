import { setup } from '@css-render/vue3-ssr'
import { defineNuxtPlugin } from '#app'

export default defineNuxtPlugin((nuxtApp) => {
  if (process.server) {
    const { collect } = setup(nuxtApp.vueApp)
    const originalRenderMeta = nuxtApp.ssrContext?.renderMeta
    nuxtApp.ssrContext = nuxtApp.ssrContext || {}
    nuxtApp.ssrContext.renderMeta = () => {
      if (!originalRenderMeta) {
        return {
          headTags: collect()
        }
      }
      const originalMeta = originalRenderMeta()
      if ('headTags' in originalMeta) {
        originalMeta.headTags += collect()
      } else {
        originalMeta.headTags = collect()
      }
      return originalMeta
    }
  }
}) 