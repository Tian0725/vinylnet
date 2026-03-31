import { createRouter, createWebHistory } from 'vue-router'
import Login from '../components/Login.vue'
import Dashboard from '../components/Dashboard.vue' // 1. Importas el componente


  const routes = [
  { path: '/', component: Login }, 
  { path: '/login', component: Login }, // Agrega esta línea explícitamente
  { path: '/dashboard', component: Dashboard }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router