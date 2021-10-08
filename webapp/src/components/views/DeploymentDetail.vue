<template>
  <h3>Deployment Detail</h3>
  <div class="deployment-container">
    namespace: {{ data.Item.namespace }} &nbsp;
    deployment: {{ data.Item.name }} &nbsp;
    <br />replicas:
    <input
      type="text"
      v-model="data.Item.replicas"
      @keyup.enter="setDeploymentReplicas(data.Item)"
    />
    <button @click="setDeploymentReplicas(data.Item)">设置</button>
  </div>

  <div class="pod-container">
    <table class="table">
      <thead>
        <tr>
          <th scope="col">Pod Name</th>
          <th scope="col">Status</th>
          <th scope="col">Node Info</th>
          <th scope="col">CreateTime</th>
          <th scope="col">Json Info</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(pod,idx) of data.Pods" :key="pod.name">
          <td>
            <div>{{ pod.name }}</div>
            <div class="color-gray">{{ pod.podIp }}</div>
          </td>
          <td>
            {{ pod.status.phase }}
            <template v-if="!validPodPhase(pod.status.phase)">
              <Suspense>
                <template #default>
                  <div>
                    <PodEventDetail :pod="pod" :namespace="pod.namespace" :name="pod.name" />
                  </div>
                </template>
                <template #fallback>
                  <div>等待中</div>
                </template>
              </Suspense>
            </template>
          </td>
          <td>
            <div>{{ pod.nodeName }}</div>
            <div class="color-gray">{{ pod.nodeIp }}</div>
          </td>
          <td>{{ pod.createTime }}</td>
          <td>
            <a :href="podOutput(pod)" target="_blank">详细信息</a>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang='ts'>import { reactive } from "@vue/reactivity";
import { onMounted, onUnmounted } from "@vue/runtime-core";
import { useRouter } from "vue-router";
import client, { Deployment } from "../../apis/deployment";
import { Pod } from "../../apis/pod"
// import PodEventDetail from "./PodEventDetail.vue";
import { defineAsyncComponent } from "vue";
const PodEventDetail = defineAsyncComponent(() => import("./PodEventDetail.vue"));



// data 页面展示参数
let data = reactive({
  Item: {} as Deployment,
  Error: "",
  Pods: [] as Pod[],
})

// req 路径参数与请求参数
let req = reactive({
  Params: {
    name: "",
    namespace: "",
  }
})

// Swither 开关
const Swither = reactive({
  LoopData: false
})

const router = useRouter()

const fetchDataLoop = async function () {
  while (Swither.LoopData) {
    getDeploymentPods()

    await new Promise(f => setTimeout(f, 2000))
  }
}

const fetchData = function () {
  getDeployment()
  getDeploymentPods()
}

const getDeployment = async function () {
  const p = req.Params
  let resp = await client.getDeploymentByName(p.namespace, p.name)
  data.Item = resp.data
  data.Error = resp.error
}

const getDeploymentPods = async function () {
  const p = req.Params

  let resp = await client.getDeploymentPodsByName(p.namespace, p.name)
  data.Pods = resp.data.sort((p1, p2): number => {
    if (p1.name >= p2.name) {
      return 1
    }
    return -1
  })
}

const setDeploymentReplicas = function (item: Deployment) {
  // console.log("item:::",item);
  client.setDeploymentReplicas(item.namespace, item.name, item.replicas)
}


// 获取 url 中的变量信息
const fetchUrlParams = function () {
  // 获取全路径
  // console.log("fullpath::::",router.currentRoute.value.fullPath);

  // 获取 query 参数
  // console.log("query::::",router.currentRoute.value.query);

  // 获取 路径参数
  // console.log("params::::",router.currentRoute.value.params);

  req.Params.name = router.currentRoute.value.params.name as string
  req.Params.namespace = router.currentRoute.value.query.namespace as string

}

let validPodPhase = function (phase: string): boolean {
  if (phase === "Running") {
    return true
  }
  return false
}

const podOutput=function(pod: Pod): string{
  return `http://127.0.0.1:8088/k8sailor/v0/pods/${pod.name}/output?namespace=${pod.namespace}&outputFormat=json`
}

onMounted(() => {
  Swither.LoopData = true
  fetchUrlParams()

  fetchData()
  fetchDataLoop()
})

onUnmounted(() => {

  Swither.LoopData = false
})
</script>

<style lang='less' scoped>
.deployment-container {
  font-size: 20px;
  font-weight: 5000;
  // font: bold;
  color: #0f65fd;
}

.color-gray {
  color: gray;
}
</style>