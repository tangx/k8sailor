import httpc, { HttpcResponse } from './httpc'
import { Pod } from './pod'
import qs from 'qs'

// Deployment 是 Deployment 的数据结构
export interface Deployment {
    images: string[]
    name: string
    namespace: string
    replicas: number
    status: {
        availableReplicas: number
        replicas: number
        unavailableReplicas: number
    }
}

// DeploymentListResponse 继承并覆盖 data， 返回 deployment 的列表
export interface DeploymentListResponse extends HttpcResponse {
    data: Deployment[]
}
// 获取所有 deployment 信息
// namespace 默认值为 defualt
// 使用 async await 解析内容
async function getAllDeployments(namespace = "default"): Promise<DeploymentListResponse> {
    const resp = await httpc.get(`/deployments?namespace=${namespace}`)
    // console.log(resp.data)
    return resp.data
}


// DeploymentResponse 继承并覆盖 data， 返回 deployment 的单个字段
export interface DeploymentResponse extends HttpcResponse {
    data: Deployment
}

async function getDeploymentByName(namespace: string, name: string): Promise<DeploymentResponse> {
    const resp = await httpc.get(`/deployments/${name}?namespace=${namespace}`)
    // console.log("deployment by name ::::", resp.data);
    return resp.data
}

// 获取与 deployment 相关的所有 pod
interface DeploymentPodsResponse extends HttpcResponse {
    data: Pod[]
}
async function getDeploymentPodsByName(namespace: string, name: string): Promise<DeploymentPodsResponse> {
    const resp = await httpc.get(`/deployments/${name}/pods?namespace=${namespace}`)

    // console.log("deployment pod::::", resp.data);
    return resp.data
}



interface SetDeploymentReplicasResponse extends HttpcResponse {
    data: Pod[]
}
// setDeploymentReplicas 设置 deployment 副本数量
async function setDeploymentReplicas(namespace: string, name: string, replicas: number): Promise<SetDeploymentReplicasResponse> {
    const params = qs.stringify({
        namespace: namespace,
        replicas: replicas,
    })

    const u = `/deployments/${name}/replicas?${params}`
    // /deployments/failed-nginx/replicas?namespace=default&replicas=3

    const resp = await httpc.put(u)

    return resp.data
}

interface DeleteDeploymentByName extends HttpcResponse {
    data: boolean
}
// deleteDeploymentByName 根据名字删除 deployment
async function deleteDeploymentByName(namespace: string, name: string): Promise<DeleteDeploymentByName> {

    const u = `/deployments/${name}?namespace=${namespace}`
    const resp = await httpc.delete(u)
    return resp.data

}

export interface createDeploymentByNameInput {
    replicas: number
    containers: {
        image: string
        ports: number[]
    }[]
}

interface CreateDeploymentByNameResponse extends HttpcResponse {
    data: boolean
}
async function createDeploymentByName(namespace: string, name: string, input: createDeploymentByNameInput): Promise<CreateDeploymentByNameResponse> {
    const u = `/deployments/${name}?namespace=${namespace}`
    // const resp = await httpc.post(u)
    const resp = await httpc.post(u, input)

    console.log("apppppp", resp.data);

    return resp.data
}



// 导出所有方法
export default {
    getAllDeployments,
    getDeploymentByName,
    getDeploymentPodsByName,
    setDeploymentReplicas,
    deleteDeploymentByName,
    createDeploymentByName,
}