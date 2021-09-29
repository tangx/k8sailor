
interface Pod {
    name: string
    namespace: string
    images: string[]
    nodeName: string
    createTime: string
    podIp: string
    status: {
        phase: string
        message: string
        reason: string
    },
    labels: {}
}
