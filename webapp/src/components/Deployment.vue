<template>
<h3>hello deployments</h3>

{{ hello }}

{{ data.items }}

</template>

<script setup lang='ts'>
import {reactive } from '@vue/reactivity'
import { onMounted } from '@vue/runtime-core'
import client,{ Deployment,DeploymentItem} from '../apis/deployment'

let data = reactive({
  items: [] as DeploymentItem[]
})



const getAll= async function() {
  const resp=await client.getAllDeployments("default")

  data.items=resp.data
  
}


onMounted(()=>{
  getAll()
})

const hello="123"
</script>

<style lang='less' scoped>
</style>