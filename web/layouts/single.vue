<template>
  <div class="single-layout">
    <n-notification-provider>
      <slot />
    </n-notification-provider>
  </div>
</template>

<script setup>

</script>

<style>
body, html {
  width: 100%;
  height: 100%;
  margin: 0;
  padding: 0;
  background-color: #fff;
}

.single-layout {
  width: 100%;
  height: 100vh;
  overflow: hidden;
}
</style> 