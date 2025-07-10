<template>
  <div v-if="show" class="fixed top-4 right-4 z-50">
    <div class="bg-green-50 border border-green-200 rounded-lg p-4 shadow-lg max-w-sm">
      <div class="flex items-start">
        <div class="flex-shrink-0">
          <i class="fas fa-check-circle text-green-400"></i>
        </div>
        <div class="ml-3 flex-1">
          <h3 class="text-sm font-medium text-green-800">成功</h3>
          <div class="mt-1 text-sm text-green-700">
            {{ message }}
          </div>
        </div>
        <div class="ml-4 flex-shrink-0">
          <button
            @click="close"
            class="inline-flex text-green-400 hover:text-green-600 focus:outline-none"
          >
            <i class="fas fa-times"></i>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

interface Props {
  message: string
  duration?: number
}

const props = withDefaults(defineProps<Props>(), {
  duration: 3000
})

const emit = defineEmits<{
  close: []
}>()

const show = ref(false)

const close = () => {
  show.value = false
  emit('close')
}

onMounted(() => {
  show.value = true
  if (props.duration > 0) {
    setTimeout(() => {
      close()
    }, props.duration)
  }
})
</script> 