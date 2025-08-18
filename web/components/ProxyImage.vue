<template>
  <n-image
    :src="proxyUrl"
    v-bind="$attrs"
  />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useImageUrl } from '~/composables/useImageUrl'

interface Props {
  src: string
}

const props = defineProps<Props>()
const { getImageUrl } = useImageUrl()

const proxyUrl = computed(() => {
  return getImageUrl(props.src)
})
</script> 