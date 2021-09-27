<template>
<h3>hello deployments</h3>

{{ hello }}

{{ data.items }}

<table class="table">
  <thead>
    <tr>
      <th scope="col">#</th>
      <th scope="col">Active</th>
      <th scope="col">Name</th>
      <th scope="col">Replicas</th>
      <th scope="col">Images</th>
      <th scope="col">Status<hr>replicas/availableReplicas/unavailableReplicas</th>
    </tr>
  </thead>
  <tbody>
    <tr v-for="(item,id) in data.items" key=:id>
      <th scope="row">{{ id }}</th>
      <td>active</td>
      <td>{{ item.name }}</td>
      <td>{{ item.replicas }}</td>
      <td>{{ item.images }}</td>
      <td>
        <span>{{ item.status.replicas }}</span> <span>/</span> 
        <span>{{ item.status.availableReplicas }}</span> <span>/</span> 
        <span>{{ item.status.unavailableReplicas }}</span>

      </td>
      
    </tr>
  </tbody>
</table>

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