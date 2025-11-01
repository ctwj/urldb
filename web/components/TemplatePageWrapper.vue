<template>
  <div>
    <SmartPageRouter :page-path="pagePath" />
  </div>
</template>

<script setup lang="ts">
const props = defineProps({
  pagePath: {
    type: String,
    required: true
  }
})
</script>