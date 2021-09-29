<template>
  <h3>Deployment Detail</h3>
  <div class="deployment-container">
    namespace: {{ data.Item.namespace }} &nbsp;
    deployment: {{ data.Item.name }} &nbsp;
  </div>

  <div class="pod-container">
    <h5>pod info</h5>
  </div>
</template>

<script setup lang='ts'>import { reactive } from "@vue/reactivity";
import { onMounted } from "@vue/runtime-core";
import { useRouter } from "vue-router";
import client, { Deployment } from "../../apis/deployment";

let data = reactive({
  Item: {} as Deployment,
  Error: "",

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
  // getPod()
}


const getDeployment = async function () {
  const p = req.Params
  let resp = await client.getDeploymentByName(p.namespace, p.name)
  data.Item = resp.data
  data.Error = resp.error

  console.log(resp);

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
</style>