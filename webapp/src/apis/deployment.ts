import httpc from './httpc'

export interface Deployment {
    code: number
    error: string
    data: DeploymentItem[] 
}

export interface DeploymentItem {
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

// 获取所有 deployment 信息
// namespace 默认值为 defualt
// 使用 async await 解析内容
async function getAllDeployments(namespace = "default"): Promise<Deployment>{
    // const resp2 = await httpc.get(`/deployments?namespace=${namespace}`)
    const resp = await httpc.get(`/deployments?namespace=${namespace}`)
    // console.log(resp.data)
    return resp.data

}


export default  {
    getAllDeployments
}