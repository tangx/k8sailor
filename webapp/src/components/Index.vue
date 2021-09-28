<template>
<div class="app-container">
  <div class="header">
    <div class="link-container">
      <router-link to="/deployments">Deployments</router-link>
    </div>
  </div>
      <hr>
  <div class="body">
    <div class="view-container">
      <router-view />
    </div>
  </div>
</div>

</template>

<script setup lang='ts'>
</script>

<style lang='less' scoped>
</style>
