<template>
  <aside class="w-64 bg-white dark:bg-gray-800 shadow-sm border-r border-gray-200 dark:border-gray-700 min-h-screen">
    <nav class="mt-8">
      <div class="px-4 space-y-2">
        <!-- 导航菜单 -->
        <NuxtLink
          v-for="item in navigationItems"
          :key="item.to"
          :to="item.to"
          class="flex items-center px-4 py-3 text-gray-700 dark:text-gray-300 hover:bg-blue-50 dark:hover:bg-blue-900/20 hover:text-blue-600 dark:hover:text-blue-400 rounded-lg transition-colors"
          :class="{ 'bg-blue-50 dark:bg-blue-900/20 text-blue-600 dark:text-blue-400': item.active($route) }"
        >
          <i :class="item.icon + ' w-5 h-5 mr-3'"></i>
          <span>{{ item.label }}</span>
        </NuxtLink>


      </div>
    </nav>
  </aside>
</template>

<script setup lang="ts">
// 简化的侧边栏组件
const navigationItems = ref([
  {
    to: '/admin',
    label: '仪表盘',
    icon: 'fas fa-tachometer-alt',
    active: (route: any) => route.path === '/admin'
  },
  {
    to: '/admin/resources',
    label: '资源管理',
    icon: 'fas fa-database',
    active: (route: any) => route.path.startsWith('/admin/resources')
  },
  {
    to: '/admin/ready-resources',
    label: '待处理资源',
    icon: 'fas fa-clock',
    active: (route: any) => route.path.startsWith('/admin/ready-resources')
  },
  {
    to: '/admin/categories',
    label: '分类管理',
    icon: 'fas fa-folder',
    active: (route: any) => route.path.startsWith('/admin/categories')
  },
  {
    to: '/admin/tags',
    label: '标签管理',
    icon: 'fas fa-tags',
    active: (route: any) => route.path.startsWith('/admin/tags')
  },
  {
    to: '/admin/platforms',
    label: '平台管理',
    icon: 'fas fa-cloud',
    active: (route: any) => route.path.startsWith('/admin/platforms')
  },
  {
    to: '/admin/accounts',
    label: '账号管理',
    icon: 'fas fa-user-shield',
    active: (route: any) => route.path.startsWith('/admin/accounts')
  },
  {
    to: '/admin/hot-dramas',
    label: '热播剧管理',
    icon: 'fas fa-film',
    active: (route: any) => route.path.startsWith('/admin/hot-dramas')
  },
  {
    to: '/admin/users',
    label: '用户管理',
    icon: 'fas fa-users',
    active: (route: any) => route.path.startsWith('/admin/users')
  },
  {
    to: '/admin/search-stats',
    label: '搜索统计',
    icon: 'fas fa-chart-line',
    active: (route: any) => route.path.startsWith('/admin/search-stats')
  },
  {
    to: '/admin/system-config',
    label: '系统配置',
    icon: 'fas fa-cog',
    active: (route: any) => route.path.startsWith('/admin/system-config')
  }
])
</script>

<style scoped>
/* 确保Font Awesome图标正确显示 */
.fas {
  font-family: 'Font Awesome 6 Free';
  font-weight: 900;
}
</style> 