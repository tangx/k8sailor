import { createApp } from 'vue'
import App from './App.vue'

// 引入 vue-router
import router from './router/index'


const app = createApp(App)
app.use(router)
app.mount('#app')
