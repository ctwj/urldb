<template>
  <div class="forbidden-layout">
    <slot />
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

.forbidden-layout {
  width: 100%;
  height: 100vh;
  overflow: hidden;
}
</style> 