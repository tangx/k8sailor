<template>
  <CreateDeploymentDashboard />

  <h3>deployments</h3>

  <div>
    <label>namespace</label>
    <input type="text" placeholder="default" v-model="data.namespace" />
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
        <th scope="col">Action</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="(item,id) in data.items" :key="item.name + '-' + item.namespace">
        <th scope="row">{{ id }}</th>
        <td>{{ isActived(item.replicas, item.status.availableReplicas) }}</td>
        <td>
          <a :href="depDetailLink(item)">{{ item.name }}</a>
        </td>
        <td>{{ item.namespace }}</td>
        <td>{{ item.replicas }}</td>
        <td>{{ imagesJoin(item.images) }}</td>
        <td>
          <span>{{ item.status.replicas }}</span>
          <span>/</span>
          <span>{{ item.status.availableReplicas }}</span>
          <span>/</span>
          <span>{{ item.status.unavailableReplicas }}</span>
        </td>
        <td>
          <button class="btn btn-danger" @click="deleteDeploymentByname(item)">删除</button>
        </td>
      </tr>
    </tbody>
  </table>
</template>

<script setup lang='ts'>
import { computed, reactive } from '@vue/reactivity'
import { onMounted, onUnmounted } from '@vue/runtime-core'
import client, { Deployment } from '@/apis/deployment'
import CreateDeploymentDashboard from './deployment/CreateDeploymentDashboard.vue'

// data 页面展示数据
let data = reactive({
  namespace: "default",
  error: "",
  items: [] as Deployment[]
})


// Switcher 开关
let Switcher = reactive({
  LoopData: false
})


// getAll 返回所有 deployment 清单
const getAllByNamespace = async function (namespace = "default") {
  const resp = await client.getAllDeployments(namespace)

  // 对数组进行排序， 避免返回结果数据相同但顺序不同时， vue 不断重新渲染。
  let _items = resp.data.sort(
    (n1: Deployment, n2: Deployment) => {
      if (n1.name >= n2.name) {
        return 1
      }
      return -1
    }
  )

  // console.log("_items:", _items);

  data.items = _items

  data.error = resp.error
  // console.log("guodegang");

}

const getAllByNamespaceLoop = async function () {
  while (Switcher.LoopData) {
    let f = getAllByNamespace("default")
    await new Promise(f => setTimeout(f, 2000));
  }
}

// imagesJoin 将镜像列表数组连接并返回成字符串
const imagesJoin = function (images: string[]): string {
  return images.join(",")
}

// isActived 判断 deployment 是否处于可用状态
const isActived = function (replicas: number, availableReplicas: number): boolean {
  return replicas === availableReplicas
}

// deployment 详情页计算属性
const depDetailLink = function (item: Deployment): string {
  return `/deployments/${item.name}?namespace=${item.namespace}`
}

// deleteDeploymentByname 根据名字删除 deployment 
const deleteDeploymentByname = async function (dep: Deployment) {
  const resp = await client.deleteDeploymentByName(dep.namespace, dep.name)

  if (resp.code !== 0) {
    console.log("delete pod failed: ", resp.error);
    return
  }
}

onMounted(() => {
  Switcher.LoopData = true
  console.log("onMounted: onOFF.loop", Switcher.LoopData);

  getAllByNamespaceLoop()
})

onUnmounted(() => {
  Switcher.LoopData = false

  console.log("onUnmounted: onOFF.loop", Switcher.LoopData);
})


</script>

<style lang='less' scoped>
</style>