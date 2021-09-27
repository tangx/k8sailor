import axios from 'axios'

// 使用 config 模式
// https://github.com/axios/axios#config-defaults
let httpc = axios.create({
    baseURL:"http://127.0.0.1:8088/k8sailor/v0"
})

export default httpc