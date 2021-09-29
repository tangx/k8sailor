<template>
  <div class="app-container">
    <div class="header">
      <h3>k8sailor - 下一代 k8s 管理平台</h3>
      <h5 class="grey">虽然代码臭，但是口号响啊</h5>
    </div>
    <hr />
    <div class="body">
      <div class="body-left">
        <div class="link-container">
          <router-link to="/deployments">Deployments</router-link>
        </div>
      </div>
      <div class="view-container body-right">
        <router-view />
      </div>
    </div>
  </div>
</template>

<script setup lang='ts'>
</script>

<style lang='less' scoped>
.grey{
  color: grey;
}
.body{
  // 悬浮排列
  display: flex;
  
  .body-left{
    width: 20%;
    margin: 5px;
    // background-color: aquamarine;
    .link-container{
      // 垂直排列
      display: flex;
      flex-direction: column;

      margin: 0px;
      padding: 3px;
      width: 100%;
    }
  }
  .body-right{
    width: 80%;
    margin: 5px;
    padding: 5px;
    // background-color: rosybrown;
  }
}
</style>
