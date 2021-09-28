<template>
<h3>deployments</h3>

<div>
  <label>namespace</label>
  <input type="text" placeholder="default" v-model="data.namespace">
  <button @click="getAllByNamespace(data.namespace)">更新数据</button>
</div>

<!-- 当数据异常的时候显示 -->
<div class="error-container" v-if="data.error">
  <h3>error</h3>
  <h3>{{ data.error }}</h3>
</div>

<!-- 当数据正常的时候显示 -->
<table class="table" v-if="!data.error">
  <thead>
    <tr>
      <th scope="col">#</th>
      <th scope="col">状态</th>
      <th scope="col">名字</th>
      <th scope="col">Namespace</th>
      <th scope="col">期望副本数量</th>
      <th scope="col">镜像列表</th>
      <th scope="col">Pod 总数/可用/不可用</th>
    </tr>
  </thead>
  <tbody>
    <tr v-for="(item,id) in data.items" key=:id>
      <th scope="row">{{ id }}</th>
      <td>{{ isActived(item.replicas,item.status.availableReplicas) }}</td>
      <td>{{ item.name }}</td>
      <td>{{ item.namespace }}</td>
      <td>{{ item.replicas }}</td>
      <td >{{ imagesJoin(item.images) }}</td>
      <td>
        <span>{{ item.status.replicas }}</span> <span>/</span> 
        <span>{{ item.status.availableReplicas }}</span> <span>/</span> 
        <span>{{ item.status.unavailableReplicas }}</span>
      </td>
    </tr>
  </tbody>

  <hr>

      <div class="view-container">
      <router-view />
    </div>
</table>

</template>

<script setup lang='ts'>
import {reactive } from '@vue/reactivity'
import { onMounted } from '@vue/runtime-core'
import client,{ DeploymentItem } from '../../apis/deployment'

let data = reactive({
  namespace:"default",
  error: "",
  items: [] as DeploymentItem[]
})

// getAll 返回所有 deployment 清单
const getAllByNamespace= async function(namespace="default") {
  const resp=await client.getAllDeployments(namespace)
  data.items=resp.data
  data.error=resp.error
}

// imagesJoin 将镜像列表数组连接并返回成字符串
const imagesJoin=function(images:string[]):string{
  return images.join(",")
}

// isActived 判断 deployment 是否处于可用状态
const isActived=function(replicas:number,availableReplicas:number):boolean{
  return replicas===availableReplicas
}

onMounted(()=>{
  getAllByNamespace()
})


</script>

<style lang='less' scoped>
</style>