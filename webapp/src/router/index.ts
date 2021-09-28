// 导入创建路由所需的组件
import { createRouter, createWebHistory } from "vue-router";

// 路由表
const routes = [
    {
        path: "/deployments",
        name: "Deployments",
        component: () => import('@/components/views/Deployment.vue'),
    },
    {
        path: "/deployments/:name",
        name: "DeploymentDetail",
        component: () => import('@/components/views/DeploymentDetail.vue')
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes: routes,
})

export default router
