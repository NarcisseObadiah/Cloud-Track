<template>
  <div class="space-y-8">
    <!-- Admin Header -->
    <div class="bg-gradient-to-r from-purple-600 to-pink-600 rounded-2xl shadow-xl p-8 text-white">
      <h1 class="text-4xl font-bold mb-2 flex items-center">
        <FontAwesomeIcon icon="user-shield" class="mr-4" />
        System Administration
      </h1>
      <p class="text-purple-100 text-lg">Complete overview of tenant resources and system health</p>
    </div>

    <!-- System Overview Cards -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
      <div class="bg-white rounded-xl shadow-lg p-6">
        <div class="flex items-center">
          <div class="bg-blue-500 rounded-lg p-3 text-white mr-4">
            <FontAwesomeIcon icon="database" class="text-xl" />
          </div>
          <div>
            <div class="text-2xl font-bold text-gray-800">{{ adminStats.totalDatabases }}</div>
            <div class="text-gray-600 text-sm">Total Databases</div>
          </div>
        </div>
      </div>
      
      <div class="bg-white rounded-xl shadow-lg p-6">
        <div class="flex items-center">
          <div class="bg-green-500 rounded-lg p-3 text-white mr-4">
            <FontAwesomeIcon icon="cubes" class="text-xl" />
          </div>
          <div>
            <div class="text-2xl font-bold text-gray-800">{{ adminStats.totalPods }}</div>
            <div class="text-gray-600 text-sm">Running Pods</div>
          </div>
        </div>
      </div>
      
      <div class="bg-white rounded-xl shadow-lg p-6">
        <div class="flex items-center">
          <div class="bg-purple-500 rounded-lg p-3 text-white mr-4">
            <FontAwesomeIcon icon="users" class="text-xl" />
          </div>
          <div>
            <div class="text-2xl font-bold text-gray-800">{{ adminStats.totalTenants }}</div>
            <div class="text-gray-600 text-sm">Active Tenants</div>
          </div>
        </div>
      </div>
      
      <div class="bg-white rounded-xl shadow-lg p-6">
        <div class="flex items-center">
          <div class="bg-orange-500 rounded-lg p-3 text-white mr-4">
            <FontAwesomeIcon icon="server" class="text-xl" />
          </div>
          <div>
            <div class="text-2xl font-bold text-gray-800">{{ adminStats.totalNamespaces }}</div>
            <div class="text-gray-600 text-sm">Namespaces</div>
          </div>
        </div>
      </div>
    </div>

    <!-- Quick Actions -->
    <div class="bg-white rounded-xl shadow-lg p-6">
      <h2 class="text-2xl font-bold mb-6 flex items-center text-gray-800">
        <FontAwesomeIcon icon="sync-alt" class="mr-3 text-yellow-600" />
        Quick Actions
      </h2>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <button
          @click="refreshAllData"
          :disabled="loading"
          class="bg-blue-600 hover:bg-blue-700 disabled:opacity-50 text-white p-4 rounded-lg flex items-center justify-center transition-colors"
        >
          <FontAwesomeIcon :icon="loading ? 'spinner' : 'sync-alt'" :class="loading ? 'animate-spin' : ''" class="mr-2" />
          Refresh All Data
        </button>
        
        <button
          @click="exportData"
          class="bg-green-600 hover:bg-green-700 text-white p-4 rounded-lg flex items-center justify-center transition-colors"
        >
          <FontAwesomeIcon icon="download" class="mr-2" />
          Export Report
        </button>
        
        <button
          @click="showSystemLogs = !showSystemLogs"
          class="bg-purple-600 hover:bg-purple-700 text-white p-4 rounded-lg flex items-center justify-center transition-colors"
        >
          <FontAwesomeIcon icon="terminal" class="mr-2" />
          {{ showSystemLogs ? 'Hide' : 'View' }} System Logs
        </button>
      </div>
    </div>

    <!-- System Logs -->
    <div v-if="showSystemLogs" class="bg-white rounded-xl shadow-lg p-6">
      <h3 class="text-xl font-bold mb-4 flex items-center text-gray-800">
        <FontAwesomeIcon icon="terminal" class="mr-3 text-green-600" />
        System Activity Logs
      </h3>
      <div class="bg-gray-900 text-green-400 p-4 rounded-lg font-mono text-sm max-h-64 overflow-y-auto">
        <div v-for="(log, index) in systemLogs" :key="index" class="mb-1">
          <span class="text-gray-500">[{{ log.timestamp }}]</span>
          <span :class="log.level === 'ERROR' ? 'text-red-400' : log.level === 'WARN' ? 'text-yellow-400' : 'text-green-400'">
            {{ log.level }}
          </span>
          {{ log.message }}
        </div>
      </div>
    </div>

    <!-- Tenant Overview -->
    <div class="bg-white rounded-xl shadow-lg overflow-hidden">
      <div class="px-6 py-4 border-b border-gray-200 bg-gray-50">
        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <h2 class="text-xl font-semibold text-gray-800 flex items-center">
            <FontAwesomeIcon icon="cubes" class="mr-3 text-blue-600" />
            Tenant Resources Overview
          </h2>
          
          <div class="flex items-center gap-3">
            <input
              v-model="searchTenant"
              placeholder="Search tenants..."
              class="px-3 py-2 border border-gray-300 rounded-lg text-sm focus:ring-blue-500 focus:border-blue-500"
            />
            <button
              @click="loadTenantData"
              :disabled="loading"
              class="bg-blue-600 hover:bg-blue-700 disabled:opacity-50 text-white px-4 py-2 rounded-lg text-sm flex items-center transition-colors"
            >
              <FontAwesomeIcon :icon="loading ? 'spinner' : 'sync-alt'" :class="loading ? 'animate-spin' : ''" class="mr-2" />
              Refresh
            </button>
          </div>
        </div>
      </div>

      <div v-if="loading" class="flex justify-center py-12">
        <div class="text-center">
          <FontAwesomeIcon icon="spinner" class="text-4xl text-blue-600 animate-spin mb-4" />
          <p class="text-gray-600">Loading tenant data...</p>
        </div>
      </div>

      <div v-else-if="error" class="p-6 text-center">
        <FontAwesomeIcon icon="exclamation-triangle" class="text-4xl text-red-500 mb-4" />
        <p class="text-red-600 font-medium">{{ error }}</p>
        <button
          @click="loadTenantData"
          class="mt-4 bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded-lg"
        >
          Try Again
        </button>
      </div>

      <div v-else-if="filteredTenantPods.length === 0" class="text-center py-12">
        <FontAwesomeIcon icon="info-circle" class="text-4xl text-gray-400 mb-4" />
        <p class="text-gray-600 text-lg">No tenant resources found</p>
      </div>

      <div v-else class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Tenant & Pod
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Status & Health
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Resources
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Location & Age
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Actions
              </th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr
              v-for="(pod, index) in filteredTenantPods"
              :key="index"
              class="hover:bg-gray-50 transition-colors"
            >
              <td class="px-6 py-4">
                <div class="flex items-center">
                  <div class="flex-shrink-0 mr-4">
                    <div :class="[
                      'w-4 h-4 rounded-full',
                      pod.status === 'Running' ? 'bg-green-500' : 
                      pod.status === 'Pending' ? 'bg-yellow-500' : 'bg-red-500'
                    ]"></div>
                  </div>
                  <div>
                    <div class="text-sm font-medium text-gray-900">
                      {{ pod.name }}
                    </div>
                    <div class="text-xs text-gray-500">
                      <FontAwesomeIcon icon="tag" class="mr-1" />
                      {{ pod.namespace }}
                    </div>
                  </div>
                </div>
              </td>
              
              <td class="px-6 py-4">
                <div class="space-y-1">
                  <span :class="[
                    'inline-flex px-2 py-1 text-xs font-semibold rounded-full',
                    pod.status === 'Running' ? 'bg-green-100 text-green-800' :
                    pod.status === 'Pending' ? 'bg-yellow-100 text-yellow-800' :
                    'bg-red-100 text-red-800'
                  ]">
                    {{ pod.status }}
                  </span>
                  <div v-if="pod.restarts > 0" class="text-xs text-orange-600 flex items-center">
                    <FontAwesomeIcon icon="sync-alt" class="mr-1" />
                    {{ pod.restarts }} restarts
                  </div>
                </div>
              </td>
              
              <td class="px-6 py-4 text-sm text-gray-600">
                <div class="space-y-1">
                  <div class="flex items-center">
                    <FontAwesomeIcon icon="microchip" class="w-3 h-3 mr-2 text-blue-500" />
                    <span class="font-mono text-xs">{{ pod.cpu || 'N/A' }}</span>
                  </div>
                  <div class="flex items-center">
                    <FontAwesomeIcon icon="memory" class="w-3 h-3 mr-2 text-green-500" />
                    <span class="font-mono text-xs">{{ pod.memory || 'N/A' }}</span>
                  </div>
                </div>
              </td>
              
              <td class="px-6 py-4 text-sm text-gray-600">
                <div class="space-y-1">
                  <div class="flex items-center">
                    <FontAwesomeIcon icon="clock" class="w-3 h-3 mr-2 text-gray-400" />
                    <span>{{ pod.age }}</span>
                  </div>
                  <div class="flex items-center text-xs">
                    <FontAwesomeIcon icon="server" class="w-3 h-3 mr-2 text-gray-400" />
                    <span>{{ pod.node || 'N/A' }}</span>
                  </div>
                </div>
              </td>
              
              <td class="px-6 py-4 text-sm">
                <div class="flex space-x-2">
                  <button
                    @click="viewTenantDetails(pod)"
                    class="text-blue-600 hover:text-blue-800 transition-colors"
                    title="View Details"
                  >
                    <FontAwesomeIcon icon="eye" />
                  </button>
                  <button
                    @click="manageTenant(pod.namespace)"
                    class="text-green-600 hover:text-green-800 transition-colors"
                    title="Manage Tenant"
                  >
                    <FontAwesomeIcon icon="cog" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Detailed Modal -->
    <div v-if="selectedTenant" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div class="bg-white rounded-xl shadow-2xl max-w-4xl w-full max-h-96 overflow-y-auto">
        <div class="flex items-center justify-between p-6 border-b border-gray-200">
          <h3 class="text-xl font-semibold text-gray-800">Tenant Details: {{ selectedTenant.namespace }}</h3>
          <button
            @click="selectedTenant = null"
            class="text-gray-400 hover:text-gray-600 transition-colors"
          >
            <FontAwesomeIcon icon="times" class="text-xl" />
          </button>
        </div>
        <div class="p-6">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div class="space-y-4">
              <div>
                <label class="text-sm font-medium text-gray-700">Pod Name</label>
                <div class="mt-1 font-mono text-sm p-3 bg-gray-100 rounded">{{ selectedTenant.name }}</div>
              </div>
              <div>
                <label class="text-sm font-medium text-gray-700">Namespace</label>
                <div class="mt-1 text-sm p-3 bg-gray-100 rounded">{{ selectedTenant.namespace }}</div>
              </div>
              <div>
                <label class="text-sm font-medium text-gray-700">Status</label>
                <div class="mt-1">
                  <span :class="[
                    'inline-flex px-3 py-1 text-sm font-semibold rounded-full',
                    selectedTenant.status === 'Running' ? 'bg-green-100 text-green-800' :
                    selectedTenant.status === 'Pending' ? 'bg-yellow-100 text-yellow-800' :
                    'bg-red-100 text-red-800'
                  ]">
                    {{ selectedTenant.status }}
                  </span>
                </div>
              </div>
            </div>
            <div class="space-y-4">
              <div>
                <label class="text-sm font-medium text-gray-700">Age</label>
                <div class="mt-1 text-sm p-3 bg-gray-100 rounded">{{ selectedTenant.age }}</div>
              </div>
              <div>
                <label class="text-sm font-medium text-gray-700">Node</label>
                <div class="mt-1 text-sm p-3 bg-gray-100 rounded">{{ selectedTenant.node || 'N/A' }}</div>
              </div>
              <div>
                <label class="text-sm font-medium text-gray-700">Restarts</label>
                <div class="mt-1 text-sm p-3 bg-gray-100 rounded">{{ selectedTenant.restarts || 0 }}</div>
              </div>
            </div>
          </div>
          
          <div class="mt-6 grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label class="text-sm font-medium text-gray-700">CPU Request</label>
              <div class="mt-1 font-mono text-sm p-3 bg-gray-100 rounded">{{ selectedTenant.cpu || 'N/A' }}</div>
            </div>
            <div>
              <label class="text-sm font-medium text-gray-700">Memory Request</label>
              <div class="mt-1 font-mono text-sm p-3 bg-gray-100 rounded">{{ selectedTenant.memory || 'N/A' }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import axios from 'axios'
import { useAuth } from '../auth/useAuth.js'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import { API_BASE_URL, API_ENDPOINTS } from '../config/api.js'

const { user } = useAuth()

// Reactive data
const tenantPods = ref([])
const loading = ref(false)
const error = ref(null)
const searchTenant = ref('')
const selectedTenant = ref(null)
const showSystemLogs = ref(false)

const adminStats = ref({
  totalDatabases: 0,
  totalPods: 0,
  totalTenants: 0,
  totalNamespaces: 0
})

const systemLogs = ref([
  { timestamp: new Date().toISOString().slice(11, 19), level: 'INFO', message: 'System startup completed' },
  { timestamp: new Date().toISOString().slice(11, 19), level: 'INFO', message: 'Database operator healthy' },
  { timestamp: new Date().toISOString().slice(11, 19), level: 'INFO', message: 'All tenant namespaces synchronized' }
])

// Computed properties
const filteredTenantPods = computed(() => {
  if (!searchTenant.value) return tenantPods.value
  const search = searchTenant.value.toLowerCase()
  return tenantPods.value.filter(pod => 
    pod.name.toLowerCase().includes(search) || 
    pod.namespace.toLowerCase().includes(search)
  )
})

// Methods
async function loadTenantData() {
  try {
    loading.value = true
    error.value = null

    const response = await axios.get(`${API_BASE_URL}${API_ENDPOINTS.ADMIN_TENANT_PODS}`, {
      headers: {
        Authorization: `Bearer ${user.value.access_token}`
      }
    })

    tenantPods.value = response.data || []
    
    // Calculate stats
    const namespaces = new Set(tenantPods.value.map(pod => pod.namespace))
    adminStats.value = {
      totalPods: tenantPods.value.length,
      totalDatabases: tenantPods.value.filter(pod => pod.name.includes('db') || pod.name.includes('postgres')).length,
      totalTenants: namespaces.size,
      totalNamespaces: namespaces.size
    }

  } catch (err) {
    error.value = err.response?.data?.error || 'Failed to load tenant data'
    tenantPods.value = []
  } finally {
    loading.value = false
  }
}

function refreshAllData() {
  loadTenantData()
  // Add log entry
  systemLogs.value.unshift({
    timestamp: new Date().toISOString().slice(11, 19),
    level: 'INFO',
    message: 'Admin data refresh initiated'
  })
}

function exportData() {
  const data = {
    timestamp: new Date().toISOString(),
    stats: adminStats.value,
    tenants: tenantPods.value
  }
  
  const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = `tenant-report-${new Date().toISOString().slice(0, 10)}.json`
  link.click()
  URL.revokeObjectURL(url)
}

function viewTenantDetails(pod) {
  selectedTenant.value = pod
}

function manageTenant(namespace) {
  // This could navigate to a tenant management page
  console.log('Managing tenant:', namespace)
}

onMounted(() => {
  loadTenantData()
})
</script>

<style scoped>
.animate-spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
</style>

