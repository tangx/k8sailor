<template>
<div class="app-container">
  <div class="header">
    <div class="link-container">
      <router-link to="/deployments">Deployments</router-link>
      <router-link to="/deployments/my-nginx-1">Deployments Detail</router-link>
      <hr>
      <router-link to="/detail">Details</router-link>
    </div>
  </div>
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
