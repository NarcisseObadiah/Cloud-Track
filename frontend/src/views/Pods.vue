<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="bg-white rounded-xl shadow-lg p-6">
      <div class="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4">
        <div>
          <h1 class="text-3xl font-bold text-gray-800 flex items-center">
            <FontAwesomeIcon icon="cubes" class="mr-3 text-purple-600" />
            Pod Management
          </h1>
          <p class="text-gray-600 mt-2">Monitor and manage your containerized applications</p>
        </div>
        
        <div class="flex flex-col sm:flex-row gap-3">
          <div class="flex-1">
            <input
              v-model="searchNamespace"
              placeholder="Enter namespace (e.g., tenant-mike)"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg shadow-sm focus:ring-purple-500 focus:border-purple-500"
            />
          </div>
          <button
            @click="fetchPods"
            :disabled="loading"
            class="bg-purple-600 hover:bg-purple-700 disabled:opacity-50 text-white px-6 py-2 rounded-lg shadow flex items-center justify-center transition-colors"
          >
            <FontAwesomeIcon :icon="loading ? 'spinner' : 'search'" :class="loading ? 'animate-spin' : ''" class="mr-2" />
            {{ loading ? 'Loading...' : 'Search Pods' }}
          </button>
          
          <button
            @click="loadMyPods"
            :disabled="loading"
            class="bg-blue-600 hover:bg-blue-700 disabled:opacity-50 text-white px-6 py-2 rounded-lg shadow flex items-center justify-center transition-colors"
          >
            <FontAwesomeIcon icon="home" class="mr-2" />
            My Pods
          </button>
        </div>
      </div>
    </div>

    <!-- Statistics Cards -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
      <div class="bg-white rounded-xl shadow-lg p-6">
        <div class="flex items-center">
          <div class="bg-green-500 rounded-lg p-3 text-white mr-4">
            <FontAwesomeIcon icon="check-circle" class="text-xl" />
          </div>
          <div>
            <div class="text-2xl font-bold text-gray-800">{{ podStats.running }}</div>
            <div class="text-gray-600 text-sm">Running</div>
          </div>
        </div>
      </div>
      
      <div class="bg-white rounded-xl shadow-lg p-6">
        <div class="flex items-center">
          <div class="bg-yellow-500 rounded-lg p-3 text-white mr-4">
            <FontAwesomeIcon icon="clock" class="text-xl" />
          </div>
          <div>
            <div class="text-2xl font-bold text-gray-800">{{ podStats.pending }}</div>
            <div class="text-gray-600 text-sm">Pending</div>
          </div>
        </div>
      </div>
      
      <div class="bg-white rounded-xl shadow-lg p-6">
        <div class="flex items-center">
          <div class="bg-red-500 rounded-lg p-3 text-white mr-4">
            <FontAwesomeIcon icon="times" class="text-xl" />
          </div>
          <div>
            <div class="text-2xl font-bold text-gray-800">{{ podStats.failed }}</div>
            <div class="text-gray-600 text-sm">Failed</div>
          </div>
        </div>
      </div>
      
      <div class="bg-white rounded-xl shadow-lg p-6">
        <div class="flex items-center">
          <div class="bg-blue-500 rounded-lg p-3 text-white mr-4">
            <FontAwesomeIcon icon="cubes" class="text-xl" />
          </div>
          <div>
            <div class="text-2xl font-bold text-gray-800">{{ pods.length }}</div>
            <div class="text-gray-600 text-sm">Total Pods</div>
          </div>
        </div>
      </div>
    </div>

    <!-- Error Message -->
    <div v-if="error" class="bg-red-50 border border-red-200 rounded-xl p-4">
      <div class="flex items-center">
        <FontAwesomeIcon icon="times" class="text-red-500 mr-3" />
        <span class="text-red-800 font-medium">{{ error }}</span>
      </div>
    </div>

    <!-- Pods Table -->
    <div class="bg-white rounded-xl shadow-lg overflow-hidden">
      <div class="px-6 py-4 border-b border-gray-200">
        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <h2 class="text-xl font-semibold text-gray-800">
            Pods in {{ currentNamespace || 'No namespace selected' }}
          </h2>
          
          <div class="flex items-center gap-3">
            <div class="flex items-center">
              <label class="text-sm text-gray-600 mr-2">Filter:</label>
              <select 
                v-model="statusFilter" 
                class="px-3 py-1 border border-gray-300 rounded-lg text-sm focus:ring-purple-500 focus:border-purple-500"
              >
                <option value="">All Status</option>
                <option value="Running">Running</option>
                <option value="Pending">Pending</option>
                <option value="Failed">Failed</option>
              </select>
            </div>
            
            <button
              @click="fetchPods"
              :disabled="!currentNamespace || loading"
              class="bg-gray-600 hover:bg-gray-700 disabled:opacity-50 text-white px-3 py-1 rounded-lg text-sm flex items-center transition-colors"
            >
              <FontAwesomeIcon icon="sync-alt" class="mr-1" />
              Refresh
            </button>
          </div>
        </div>
      </div>

      <div v-if="loading" class="flex justify-center py-12">
        <div class="text-center">
          <FontAwesomeIcon icon="spinner" class="text-4xl text-purple-600 animate-spin mb-4" />
          <p class="text-gray-600">Loading pods...</p>
        </div>
      </div>

      <div v-else-if="filteredPods.length === 0" class="text-center py-12">
        <FontAwesomeIcon icon="info-circle" class="text-4xl text-gray-400 mb-4" />
        <p class="text-gray-600 text-lg mb-2">No pods found</p>
        <p class="text-gray-500 text-sm">
          {{ currentNamespace ? 'No pods match your current filter.' : 'Enter a namespace to view pods.' }}
        </p>
      </div>

      <div v-else class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Pod Details
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Status
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Resources
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Age & Location
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Actions
              </th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr
              v-for="(pod, index) in filteredPods"
              :key="index"
              class="hover:bg-gray-50 transition-colors"
            >
              <td class="px-6 py-4">
                <div class="flex items-center">
                  <div class="flex-shrink-0 mr-4">
                    <div :class="[
                      'w-3 h-3 rounded-full',
                      pod.status === 'Running' ? 'bg-green-500' : 
                      pod.status === 'Pending' ? 'bg-yellow-500' : 'bg-red-500'
                    ]"></div>
                  </div>
                  <div>
                    <div class="text-sm font-medium text-gray-900 font-mono">
                      {{ pod.name }}
                    </div>
                    <div class="text-xs text-gray-500">
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
                  <div v-if="pod.restarts > 0" class="text-xs text-orange-600">
                    {{ pod.restarts }} restarts
                  </div>
                </div>
              </td>
              
              <td class="px-6 py-4 text-sm text-gray-600">
                <div class="space-y-1">
                  <div class="flex items-center">
                    <span class="w-12 text-xs text-gray-500">CPU:</span>
                    <span class="font-mono">{{ pod.cpu || 'N/A' }}</span>
                  </div>
                  <div class="flex items-center">
                    <span class="w-12 text-xs text-gray-500">RAM:</span>
                    <span class="font-mono">{{ pod.memory || 'N/A' }}</span>
                  </div>
                </div>
              </td>
              
              <td class="px-6 py-4 text-sm text-gray-600">
                <div class="space-y-1">
                  <div>{{ pod.age }}</div>
                  <div class="text-xs text-gray-500">{{ pod.node || 'N/A' }}</div>
                </div>
              </td>
              
              <td class="px-6 py-4 text-sm">
                <div class="flex space-x-2">
                  <button
                    @click="viewPodDetails(pod)"
                    class="text-blue-600 hover:text-blue-800 transition-colors"
                    title="View Details"
                  >
                    <FontAwesomeIcon icon="eye" />
                  </button>
                  <button
                    @click="copyPodName(pod.name)"
                    class="text-green-600 hover:text-green-800 transition-colors"
                    title="Copy Name"
                  >
                    <FontAwesomeIcon icon="key" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Pod Details Modal -->
    <div v-if="selectedPod" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div class="bg-white rounded-xl shadow-2xl max-w-2xl w-full max-h-96 overflow-y-auto">
        <div class="flex items-center justify-between p-6 border-b border-gray-200">
          <h3 class="text-xl font-semibold text-gray-800">Pod Details</h3>
          <button
            @click="selectedPod = null"
            class="text-gray-400 hover:text-gray-600 transition-colors"
          >
            <FontAwesomeIcon icon="times" class="text-xl" />
          </button>
        </div>
        <div class="p-6 space-y-4">
          <div>
            <label class="text-sm font-medium text-gray-700">Pod Name</label>
            <div class="mt-1 font-mono text-sm p-2 bg-gray-100 rounded">{{ selectedPod.name }}</div>
          </div>
          <div>
            <label class="text-sm font-medium text-gray-700">Namespace</label>
            <div class="mt-1 text-sm p-2 bg-gray-100 rounded">{{ selectedPod.namespace }}</div>
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="text-sm font-medium text-gray-700">Status</label>
              <div class="mt-1 text-sm p-2 bg-gray-100 rounded">{{ selectedPod.status }}</div>
            </div>
            <div>
              <label class="text-sm font-medium text-gray-700">Age</label>
              <div class="mt-1 text-sm p-2 bg-gray-100 rounded">{{ selectedPod.age }}</div>
            </div>
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="text-sm font-medium text-gray-700">CPU Request</label>
              <div class="mt-1 font-mono text-sm p-2 bg-gray-100 rounded">{{ selectedPod.cpu || 'N/A' }}</div>
            </div>
            <div>
              <label class="text-sm font-medium text-gray-700">Memory Request</label>
              <div class="mt-1 font-mono text-sm p-2 bg-gray-100 rounded">{{ selectedPod.memory || 'N/A' }}</div>
            </div>
          </div>
          <div>
            <label class="text-sm font-medium text-gray-700">Node</label>
            <div class="mt-1 text-sm p-2 bg-gray-100 rounded">{{ selectedPod.node || 'N/A' }}</div>
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
const searchNamespace = ref('')
const currentNamespace = ref('')
const pods = ref([])
const loading = ref(false)
const error = ref(null)
const statusFilter = ref('')
const selectedPod = ref(null)

// Computed properties
const username = computed(() => {
  const raw = user.value?.profile?.preferred_username || user.value?.profile?.email || 'user'
  const base = raw.split('@')[0]
  return base.toLowerCase().replace(/[^a-z0-9-]/g, '')
})

const podStats = computed(() => {
  return {
    running: pods.value.filter(p => p.status === 'Running').length,
    pending: pods.value.filter(p => p.status === 'Pending').length,
    failed: pods.value.filter(p => p.status !== 'Running' && p.status !== 'Pending').length
  }
})

const filteredPods = computed(() => {
  if (!statusFilter.value) return pods.value
  return pods.value.filter(pod => pod.status === statusFilter.value)
})

// Methods
function loadMyPods() {
  searchNamespace.value = `tenant-${username.value}`
  fetchPods()
}

async function fetchPods() {
  const namespace = searchNamespace.value.trim()
  if (!namespace) {
    error.value = 'Please enter a namespace.'
    return
  }

  try {
    loading.value = true
    error.value = null
    currentNamespace.value = namespace

    const response = await axios.get(`${API_BASE_URL}${API_ENDPOINTS.PODS.replace(':namespace', namespace)}`, {
      headers: {
        Authorization: `Bearer ${user.value.access_token}`
      }
    })

    pods.value = response.data || []
  } catch (err) {
    error.value = err.response?.data?.error || 'Failed to fetch pods'
    pods.value = []
    currentNamespace.value = ''
  } finally {
    loading.value = false
  }
}

function viewPodDetails(pod) {
  selectedPod.value = pod
}

function copyPodName(name) {
  navigator.clipboard.writeText(name).then(() => {
    // You could add a toast notification here
    console.log('Pod name copied to clipboard')
  })
}

onMounted(() => {
  // Auto-load user's pods if available
  if (username.value) {
    loadMyPods()
  }
})
</script>

<style scoped>
.animate-spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style>
      headers: {
        Authorization: `Bearer ${user.value.access_token}`
      }
    })

    pods.value = res.data
  } catch (err) {
    error.value = err.response?.data?.message || err.message || 'Failed to fetch pods.'
  } finally {
    loading.value = false
  }
}
