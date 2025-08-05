<template>
  <div class="space-y-8">
    <!-- Hero Section -->
    <div class="bg-gradient-to-r from-blue-600 to-purple-600 rounded-2xl shadow-xl p-8 text-white">
      <div class="flex items-center justify-between mb-6">
        <div>
          <h1 class="text-4xl font-bold mb-2">Welcome Back, {{ displayName }}!</h1>
          <p class="text-blue-100 text-lg">Your cloud infrastructure at a glance</p>
        </div>
        <div class="text-right">
          <div class="text-3xl font-bold">{{ currentTime }}</div>
          <div class="text-blue-100">{{ currentDate }}</div>
        </div>
      </div>
      
      <!-- Quick Stats -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mt-8">
        <div class="bg-white/10 backdrop-blur-sm rounded-xl p-6">
          <div class="flex items-center">
            <FontAwesomeIcon icon="database" class="text-3xl mr-4" />
            <div>
              <div class="text-2xl font-bold">{{ stats.databases }}</div>
              <div class="text-blue-100">Active Databases</div>
            </div>
          </div>
        </div>
        <div class="bg-white/10 backdrop-blur-sm rounded-xl p-6">
          <div class="flex items-center">
            <FontAwesomeIcon icon="cubes" class="text-3xl mr-4" />
            <div>
              <div class="text-2xl font-bold">{{ stats.pods }}</div>
              <div class="text-blue-100">Running Pods</div>
            </div>
          </div>
        </div>
        <div class="bg-white/10 backdrop-blur-sm rounded-xl p-6">
          <div class="flex items-center">
            <FontAwesomeIcon icon="server" class="text-3xl mr-4" />
            <div>
              <div class="text-2xl font-bold">{{ stats.uptime }}</div>
              <div class="text-blue-100">System Uptime</div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Quick Actions -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div class="bg-white rounded-xl shadow-lg p-6 hover:shadow-xl transition-shadow duration-300">
        <div class="flex items-center mb-4">
          <div class="bg-blue-500 rounded-lg p-3 text-white mr-4">
            <FontAwesomeIcon icon="plus-circle" class="text-2xl" />
          </div>
          <h3 class="text-xl font-semibold">Create Database</h3>
        </div>
        <p class="text-gray-600 mb-4">Provision a new PostgreSQL database with just one click</p>
        <RouterLink 
          to="/databases" 
          class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg inline-block transition-colors"
        >
          Get Started
        </RouterLink>
      </div>

      <div class="bg-white rounded-xl shadow-lg p-6 hover:shadow-xl transition-shadow duration-300">
        <div class="flex items-center mb-4">
          <div class="bg-green-500 rounded-lg p-3 text-white mr-4">
            <FontAwesomeIcon icon="cubes" class="text-2xl" />
          </div>
          <h3 class="text-xl font-semibold">Manage Pods</h3>
        </div>
        <p class="text-gray-600 mb-4">Monitor and manage your running containerized applications</p>
        <RouterLink 
          to="/pods" 
          class="bg-green-600 hover:bg-green-700 text-white px-4 py-2 rounded-lg inline-block transition-colors"
        >
          View Pods
        </RouterLink>
      </div>

      <div v-if="isAdmin" class="bg-white rounded-xl shadow-lg p-6 hover:shadow-xl transition-shadow duration-300">
        <div class="flex items-center mb-4">
          <div class="bg-purple-500 rounded-lg p-3 text-white mr-4">
            <FontAwesomeIcon icon="user-shield" class="text-2xl" />
          </div>
          <h3 class="text-xl font-semibold">Admin Panel</h3>
        </div>
        <p class="text-gray-600 mb-4">Manage all tenant resources and system administration</p>
        <RouterLink 
          to="/admin" 
          class="bg-purple-600 hover:bg-purple-700 text-white px-4 py-2 rounded-lg inline-block transition-colors"
        >
          Admin View
        </RouterLink>
      </div>
    </div>

    <!-- Recent Activity -->
    <div class="bg-white rounded-xl shadow-lg p-6">
      <h2 class="text-2xl font-bold mb-6 flex items-center">
        <FontAwesomeIcon icon="clock" class="mr-3 text-blue-600" />
        Recent Activity
      </h2>
      
      <div v-if="loading" class="flex justify-center py-8">
        <FontAwesomeIcon icon="spinner" class="text-3xl animate-spin text-blue-600" />
      </div>
      
      <div v-else-if="recentActivity.length === 0" class="text-center py-8 text-gray-500">
        <FontAwesomeIcon icon="info-circle" class="text-4xl mb-4" />
        <p>No recent activity. Start by creating your first database!</p>
      </div>
      
      <div v-else class="space-y-4">
        <div 
          v-for="(activity, index) in recentActivity" 
          :key="index"
          class="flex items-center p-4 bg-gray-50 rounded-lg"
        >
          <div class="flex-shrink-0 mr-4">
            <div :class="[
              'rounded-full p-2 text-white',
              activity.type === 'database' ? 'bg-blue-500' : 'bg-green-500'
            ]">
              <FontAwesomeIcon :icon="activity.type === 'database' ? 'database' : 'cubes'" />
            </div>
          </div>
          <div class="flex-1">
            <p class="font-medium">{{ activity.message }}</p>
            <p class="text-sm text-gray-500">{{ activity.time }}</p>
          </div>
          <div :class="[
            'px-3 py-1 rounded-full text-xs font-medium',
            activity.status === 'Running' ? 'bg-green-100 text-green-800' : 
            activity.status === 'Pending' ? 'bg-yellow-100 text-yellow-800' : 
            'bg-gray-100 text-gray-800'
          ]">
            {{ activity.status }}
          </div>
        </div>
      </div>
    </div>

    <!-- System Status -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <div class="bg-white rounded-xl shadow-lg p-6">
        <h3 class="text-xl font-bold mb-4 flex items-center">
          <FontAwesomeIcon icon="heartbeat" class="mr-3 text-red-500" />
          System Health
        </h3>
        <div class="space-y-3">
          <div class="flex justify-between items-center">
            <span>API Status</span>
            <span class="flex items-center text-green-600">
              <div class="w-3 h-3 bg-green-500 rounded-full mr-2 animate-pulse"></div>
              Healthy
            </span>
          </div>
          <div class="flex justify-between items-center">
            <span>Database Operator</span>
            <span class="flex items-center text-green-600">
              <div class="w-3 h-3 bg-green-500 rounded-full mr-2 animate-pulse"></div>
              Online
            </span>
          </div>
          <div class="flex justify-between items-center">
            <span>Kubernetes Cluster</span>
            <span class="flex items-center text-green-600">
              <div class="w-3 h-3 bg-green-500 rounded-full mr-2 animate-pulse"></div>
              Ready
            </span>
          </div>
        </div>
      </div>

      <div class="bg-white rounded-xl shadow-lg p-6">
        <h3 class="text-xl font-bold mb-4 flex items-center">
          <FontAwesomeIcon icon="chart-line" class="mr-3 text-blue-500" />
          Quick Stats
        </h3>
        <div class="space-y-3">
          <div class="flex justify-between items-center">
            <span>Your Namespace</span>
            <code class="bg-gray-100 px-2 py-1 rounded text-sm">tenant-{{ username }}</code>
          </div>
          <div class="flex justify-between items-center">
            <span>Available Resources</span>
            <span class="text-green-600 font-medium">Unlimited</span>
          </div>
          <div class="flex justify-between items-center">
            <span>Storage Quota</span>
            <span class="text-blue-600 font-medium">100GB</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { RouterLink } from 'vue-router'
import { useAuth } from '../auth/useAuth.js'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import { API_BASE_URL, API_ENDPOINTS } from '../config/api.js'
import axios from 'axios'

const { user, ensureValidToken, clearAuthAndRestart, login } = useAuth()

const stats = ref({
  databases: 0,
  pods: 0,
  uptime: '99.9%'
})

const loading = ref(false)
const recentActivity = ref([])
const currentTime = ref('')
const currentDate = ref('')

let timeInterval = null

// Computed properties
const displayName = computed(() => {
  return user.value?.profile?.name || user.value?.profile?.preferred_username || 'User'
})

const username = computed(() => {
  const raw = user.value?.profile?.preferred_username || user.value?.profile?.email || 'user'
  const base = raw.split('@')[0]
  return base.toLowerCase().replace(/[^a-z0-9-]/g, '')
})

const roles = computed(() => {
  const rolesObj = user.value?.profile?.['urn:zitadel:iam:org:project:328396684147088791:roles'] || {}
  return Object.keys(rolesObj)
})

const isAdmin = computed(() => roles.value.includes('admin'))

// Methods
function updateTime() {
  const now = new Date()
  currentTime.value = now.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
  currentDate.value = now.toLocaleDateString([], { 
    weekday: 'long', 
    year: 'numeric', 
    month: 'long', 
    day: 'numeric' 
  })
}

async function fetchStats() {
  if (!user.value) {
    console.log('No user found, skipping stats fetch')
    return
  }
  
  loading.value = true
  try {
    // Check token status first
    const token = await ensureValidToken()
    if (!token) {
      console.log('No valid token available, skipping API call')
      recentActivity.value = [
        {
          type: 'info',
          message: 'Please login to view your data',
          time: 'Now',
          status: 'Pending'
        }
      ]
      return
    }
    
    console.log('Token available, making API call...')
    
    // Fetch user's pods to get statistics
    const namespace = `tenant-${username.value}`
    const response = await axios.get(`${API_BASE_URL}${API_ENDPOINTS.PODS.replace(':namespace', namespace)}`, {
      headers: {
        Authorization: `Bearer ${token}`
      }
    })
    
    // Check if response.data exists and is an array
    const pods = response.data || []
    
    stats.value.pods = pods.length
    
    // Count databases (PostgreSQL pods)
    const dbPods = pods.filter(pod => 
      pod.name && (pod.name.includes('db') || pod.name.includes('postgres'))
    )
    stats.value.databases = dbPods.length

    // Create recent activity from pods
    recentActivity.value = pods.slice(0, 5).map(pod => ({
      type: (pod.name && pod.name.includes('db')) ? 'database' : 'pod',
      message: `${pod.name || 'Unknown'} is ${(pod.status || 'unknown').toLowerCase()}`,
      time: `Started ${pod.age || 'unknown'} ago`,
      status: pod.status || 'unknown'
    }))
    
  } catch (error) {
    console.error('Failed to fetch stats:', error)
    
    // Don't auto-restart for now to avoid infinite reload
    // if (error.response?.status === 401) {
    //   console.log('Authentication failed, clearing auth and restarting...')
    //   await clearAuthAndRestart()
    //   return
    // }
    
    // Set some default activity for demo
    recentActivity.value = [
      {
        type: 'database',
        message: 'Welcome to your PaaS Dashboard',
        time: 'Just now',
        status: 'Ready'
      }
    ]
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  updateTime()
  timeInterval = setInterval(updateTime, 1000)
  fetchStats()
})

onUnmounted(() => {
  if (timeInterval) {
    clearInterval(timeInterval)
  }
})
</script>

<style scoped>
.animate-pulse {
  animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: .5;
  }
}
</style>
