<template>
  <button @click="displaySwitcher">打开/关闭创建面板</button>
  <div class="edit-box" v-if="isDisplay">
    <div class="deployment">
      <label>
        Name
        <input type="text" v-model="deployment.name" />
      </label>
      <label>
        Namespace
        <input type="text" v-model="deployment.namespace" />
      </label>
      <label>
        Replicas
        <input type="number" v-model="deployment.replicas" />
      </label>

      <div class="containers">
        <div class="container" v-for="(c,id) in deployment.containers" :key="id">
          <div class="title">
            <div>Container 编号: {{ id }}</div>
            <div>
              <button @click="deleteContainer(deployment.containers, id)">删除当前container</button>
            </div>
          </div>

          <div class="image">
            <label>Image</label>
            <input type="text" v-model="c.image" />
          </div>
          <div class="ports">
            <div class="port" v-for="(port,id) in c.ports" :key="id">
              <label>Port</label>
              <!-- input 不能直接使用 v-model 将数组绑定到 value 中 -->
              <!-- 但是 input 可以使用 v-bind 绑定 -->
              <!-- 应该是双向绑定数据修改的问题导致的 -->
              <input type="text" v-model.number="port.value" />
              <button @click="deleteArrayItem(c.ports, id)">删除</button>
            </div>

            <button @click="appendPort(c.ports)">增加 Port</button>
          </div>
        </div>

        <button @click="appendContainer(deployment.containers)">增加 Container</button>
      </div>
    </div>

    <div class="actions">
      <button @click="applyCreatation">确认创建</button>
      <button @click="isDisplay = false">取消</button>
    </div>
  </div>
</template>

<script setup lang='ts'>
import { reactive, ref } from '@vue/reactivity';
import httpc, { createDeploymentByNameInput } from '@/apis/deployment'

interface container {
  image: string
  ports: port[],
}

interface port {
  value: number
}

interface _container {
  image: string
  ports: number[]
}

let deployment = reactive({
  name: "my-nginx-23421" as string,
  namespace: "default" as string,
  replicas: 10 as number,
  containers: [
    {
      image: "hub.example.com/nginx:alpine",
      ports: [
        {
          value: 80,
        }
      ],
    }
  ] as container[]
})

const isDisplay = ref(false)

const displaySwitcher = function () {
  isDisplay.value = !isDisplay.value
}

const appendPort = function (ports: port[]) {
  ports.push({
    value: 0,
  })
}

const deleteArrayItem = function (arr: port[] | undefined, index: number) {
  if (arr) {
    arr.splice(index, 1)
  }
}

const appendContainer = function (cs: container[]) {
  cs.push({
    ports: [] as port[]
  } as container)
}

const deleteContainer = function (arr: container[] | undefined, index: number) {
  if (arr) {
    arr.splice(index, 1)
  }
}

const applyCreatation = async function () {
  // console.log(deployment);

  const input = {
    replicas: deployment.replicas,
    containers: transContainers(deployment.containers),
  } as createDeploymentByNameInput

  console.log(input);


  const resp =await httpc.createDeploymentByName(deployment.namespace, deployment.name, input)

  console.log(resp.data);
  
}

const transContainers = function (cs: container[]): _container[] {
  const _cs = [] as _container[]
  for (const c of cs) {
    const ps = [] as number[]
    for (const p of c.ports) {
      ps.push(p.value)
    }

    _cs.push({
      image: c.image,
      ports: ps,
    })
  }

  return _cs
}

</script>

<style lang='less' scoped>
.edit-box {
  .actions {
    background-color: cornflowerblue;
    margin-top: 10px;
    display: flex;
    justify-content: space-evenly;
  }
  .containers {
    border-style: solid;
    border-width: 1px;

    .container {
      margin-top: 20px;
      margin-bottom: 20px;
      width: 800px;
      background-color: cadetblue;

      display: flex;
      flex-direction: column;
      justify-content: start;

      .title {
        display: flex;
        justify-content: space-between;
      }

      .image {
        width: 100%;
        background-color: chocolate;

        label {
          width: 20%;
        }
        input {
          width: 80%;
        }
      }
    }

    .ports {
      border-style: solid;
      border-width: 1px;

      .port {
        label {
          width: 20%;
        }
        input {
          width: 70%;
        }
        button {
          width: 10%;
        }
      }
    }
  }
}
</style>
