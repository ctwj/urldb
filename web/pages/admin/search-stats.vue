<template>
  <div class="min-h-screen bg-gray-50">
    <div class="container mx-auto px-4 py-8">
      <!-- 页面标题 -->
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-gray-900">搜索统计</h1>
        <p class="text-gray-600 mt-2">查看搜索量统计和热门关键词分析</p>
      </div>

      <!-- 统计卡片 -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <div class="bg-white rounded-lg shadow p-6">
          <div class="flex items-center">
            <div class="p-3 rounded-full bg-blue-100 text-blue-600">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
              </svg>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">今日搜索</p>
              <p class="text-2xl font-bold text-gray-900">{{ stats.todaySearches }}</p>
            </div>
          </div>
        </div>

        <div class="bg-white rounded-lg shadow p-6">
          <div class="flex items-center">
            <div class="p-3 rounded-full bg-green-100 text-green-600">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path>
              </svg>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">本周搜索</p>
              <p class="text-2xl font-bold text-gray-900">{{ stats.weekSearches }}</p>
            </div>
          </div>
        </div>

        <div class="bg-white rounded-lg shadow p-6">
          <div class="flex items-center">
            <div class="p-3 rounded-full bg-purple-100 text-purple-600">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"></path>
              </svg>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">本月搜索</p>
              <p class="text-2xl font-bold text-gray-900">{{ stats.monthSearches }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- 搜索趋势图表 -->
      <div class="bg-white rounded-lg shadow p-6 mb-8">
        <h2 class="text-xl font-semibold text-gray-900 mb-4">搜索趋势</h2>
        <div class="h-64">
          <canvas ref="trendChart"></canvas>
        </div>
      </div>

      <!-- 热门关键词 -->
      <div class="bg-white rounded-lg shadow p-6">
        <h2 class="text-xl font-semibold text-gray-900 mb-4">热门关键词</h2>
        <div class="space-y-4">
          <div v-for="keyword in stats.hotKeywords" :key="keyword.keyword" 
               class="flex items-center justify-between p-4 bg-gray-50 rounded-lg">
            <div class="flex items-center">
              <span class="inline-flex items-center justify-center w-8 h-8 bg-blue-100 text-blue-600 rounded-full text-sm font-medium mr-3">
                {{ keyword.rank }}
              </span>
              <span class="text-gray-900 font-medium">{{ keyword.keyword }}</span>
            </div>
            <div class="flex items-center">
              <span class="text-gray-600 mr-2">{{ keyword.count }}次</span>
              <div class="w-24 bg-gray-200 rounded-full h-2">
                <div class="bg-blue-600 h-2 rounded-full" 
                     :style="{ width: getPercentage(keyword.count) + '%' }"></div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
// 设置页面布局
definePageMeta({
  layout: 'admin'
})

import { ref, onMounted, computed } from 'vue'
import Chart from 'chart.js/auto'

const stats = ref({
  todaySearches: 0,
  weekSearches: 0,
  monthSearches: 0,
  hotKeywords: [],
  searchTrend: {
    days: [],
    values: []
  }
})

const trendChart = ref(null)
let chart = null

// 获取百分比
const getPercentage = (count) => {
  if (stats.value.hotKeywords.length === 0) return 0
  const maxCount = Math.max(...stats.value.hotKeywords.map(k => k.count))
  return Math.round((count / maxCount) * 100)
}

// 加载搜索统计
const loadSearchStats = async () => {
  try {
    const response = await fetch('/api/search-stats')
    if (response.ok) {
      const data = await response.json()
      stats.value = data
      
      // 更新图表
      updateChart()
    }
  } catch (error) {
    console.error('加载搜索统计失败:', error)
  }
}

// 更新图表
const updateChart = () => {
  if (chart) {
    chart.destroy()
  }

  const ctx = trendChart.value.getContext('2d')
  chart = new Chart(ctx, {
    type: 'line',
    data: {
      labels: stats.value.searchTrend.days,
      datasets: [{
        label: '搜索量',
        data: stats.value.searchTrend.values,
        borderColor: 'rgb(59, 130, 246)',
        backgroundColor: 'rgba(59, 130, 246, 0.1)',
        tension: 0.4,
        fill: true
      }]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: {
          display: false
        }
      },
      scales: {
        y: {
          beginAtZero: true,
          grid: {
            color: 'rgba(0, 0, 0, 0.1)'
          }
        },
        x: {
          grid: {
            color: 'rgba(0, 0, 0, 0.1)'
          }
        }
      }
    }
  })
}

onMounted(() => {
  loadSearchStats()
})
</script> 