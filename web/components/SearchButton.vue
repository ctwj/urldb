<template>
  <!-- 搜索按钮组件 -->
  <div class="search-button-container">
    <!-- 搜索按钮 -->
    <n-button
      size="tiny"
      type="tertiary"
      round
      ghost
      class="!px-2 !py-1 !text-xs !text-white dark:!text-white !border-white/30 hover:!border-white"
      @click="openSearch"
    >
      <i class="fas fa-search text-xs"></i>
      <span class="ml-1 hidden sm:inline">搜索</span>
    </n-button>

    <!-- 完整的搜索弹窗组件 -->
    <SearchModal ref="searchModalRef" />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import SearchModal from './SearchModal.vue'

// 搜索弹窗的引用
const searchModalRef = ref()

// 打开搜索弹窗
const openSearch = () => {
  searchModalRef.value?.show()
}

// 暴露给父组件的方法
defineExpose({
  openSearch,
  closeSearch: () => searchModalRef.value?.hide(),
  toggleSearch: () => searchModalRef.value?.toggle()
})
</script>

<style scoped>
.search-button-container {
  display: inline-block;
}
</style>