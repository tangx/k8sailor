<template>
  <h3>Deployment Detail</h3>
  <div class="deployment-container">
    namespace: {{ data.Item.namespace }} &nbsp;
    deployment: {{ data.Item.name }} &nbsp;
    <br>
    replicas: <input type="text" v-model="data.Item.replicas">
    <button @click="setDeploymentReplicas(data.Item)">设置</button>
  </div>

  <div class="pod-container">
    <table class="table">
  <thead>
    <tr>
      <th scope="col">Pod Name</th>
      <th scope="col">Status</th>
      <th scope="col">Node</th>
      <th scope="col">Ipaddr</th>
      <th scope="col">CreateTime</th>
    </tr>
  </thead>
  <tbody>
    <tr v-for="(pod,idx) of data.Pods" :key="pod.name">
      <td>{{ pod.name }}</td>
      <td>{{ pod.status.phase }}</td>
      <td>{{ pod.nodeName }}</td>
      <td>{{ pod.podIp }}</td>
      <td>{{ pod.createTime }}</td>
    </tr>
  </tbody>
</table>
  </div>
</template>

<script setup lang='ts'>import { reactive } from "@vue/reactivity";
import { onMounted } from "@vue/runtime-core";
import { useRouter } from "vue-router";
import client, { Deployment } from "../../apis/deployment";

let data = reactive({
  Item: {} as Deployment,
  Error: "",
  Pods: [] as Pod[],
})

let req = reactive({
  Params: {
    name: "",
    namespace: "",
  }
})

const router = useRouter()

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

const getDeploymentPods=async function(){
  const p = req.Params

  let resp=await client.getDeploymentPodsByName(p.namespace,p.name)
  data.Pods=resp.data
}

const setDeploymentReplicas=function(item:Deployment){
  // console.log("item:::",item);
  client.setDeploymentReplicas(item.namespace,item.name,item.replicas)

  fetchData()
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


onMounted(() => {
  fetchUrlParams()

  fetchData()
})
</script>

<style lang='less' scoped>
.deployment-container{
  font-size: 20px;
  font-weight: 5000;
  // font: bold;
  color: #0F65FD;
}
</style>