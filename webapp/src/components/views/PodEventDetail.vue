<template>
  <div>
    <!-- <h3>123123</h3> -->

    <!-- <h3>{{ event }}</h3> -->
    <label class="error">{{ event.message }}</label>
  </div>
</template>

<script setup lang='ts'>
import { PropType } from '@vue/runtime-core'
import client, { Pod, PodEvent } from '../../apis/pod'

let props = defineProps({
  namespace: String,
  name: String,
  pod: {
    type: Object as PropType<Pod>,
  },
})

let event = {} as PodEvent

// let resp = await client.getPodEventByName(props.namespace as string, props.name as string)
let resp = await client.getPodEventByName(props.pod.namespace, props.pod.name)

event = resp.data
// console.log(`${props.name} ::`,event);

</script>


<style lang='less' scoped>
.error {
  color: red;
}
</style>