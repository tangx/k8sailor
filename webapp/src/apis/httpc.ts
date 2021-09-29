import axios from 'axios'

// HttpcResponse 是 Server 端的基础响应结构体
// 具体的接口响应结果， 需要接口自行实现 HttpcReponse 的继承与 data 字段的覆盖
export interface HttpcResponse {
    code: number
    error: string
    data: Object
}


// 使用 config 模式
// https://github.com/axios/axios#config-defaults
let httpc = axios.create({
    baseURL: "http://127.0.0.1:8088/k8sailor/v0"
})

export default httpc