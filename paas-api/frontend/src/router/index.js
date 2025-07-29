import { createRouter, createWebHistory } from 'vue-router'

import Dashboard from '../views/Dashboard.vue'
import Databases from '../views/Databases.vue'
import Pods from '../views/Pods.vue'
import Admin from '../views/Admin.vue'
import Callback from '../views/Callback.vue'

const routes = [
  { path: '/', component: Dashboard },
  { path: '/databases', component: Databases },
  { path: '/pods', component: Pods },
  { path: '/admin', component: Admin },
  { path: '/callback', component: Callback },

]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
