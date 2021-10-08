import httpc, { HttpcResponse } from './httpc'

export interface Pod {
    name: string
    namespace: string
    images: string[]
    nodeName: string
    nodeIp: string
    createTime: string
    podIp: string
    status: {
        phase: string
        message: string
        reason: string
    },
    labels: {}
}


export interface PodEvent {
    name: string
    namespace: string
    reson: string
    message: string
}

interface PodEventResponse extends HttpcResponse {
    data: PodEvent
}

async function getPodEventByName(namespace: string, name: string): Promise<PodEventResponse> {
    const uri = `/pods/${name}/event?namespace=${namespace}`

    const resp = await httpc.get(uri)

    return resp.data
}

export default {
    getPodEventByName
}