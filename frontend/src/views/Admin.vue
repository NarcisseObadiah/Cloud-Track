<template>
  <div class="bg-white p-6 rounded-xl shadow-md max-w-7xl mx-auto mt-10">
    <h2 class="text-3xl font-bold mb-6 text-purple-700 flex items-center gap-2">
      <FontAwesomeIcon icon="server" />
      Admin Pod Overview
    </h2>

    <button
      @click="fetchAll"
      class="bg-purple-600 hover:bg-purple-700 text-white py-2 px-4 rounded mb-4 shadow-md flex items-center gap-2"
    >
  <FontAwesomeIcon icon="spinner" spin v-if="loading" />
      <FontAwesomeIcon icon="cubes" v-else />
      Refresh Pods
    </button>

    <div v-if="error" class="bg-red-100 text-red-700 p-3 rounded mb-4 border border-red-200">
      {{ error }}
    </div>

    <div v-if="pods.length" class="overflow-auto rounded shadow">
      <table class="min-w-full divide-y divide-gray-200 text-sm">
        <thead class="bg-gray-100">
          <tr>
            <th class="px-4 py-3 text-left font-semibold text-gray-600">Pod Name</th>
            <th class="px-4 py-3 text-left font-semibold text-gray-600">Tenant</th>
            <th class="px-4 py-3 text-left font-semibold text-gray-600">Namespace</th>
            <th class="px-4 py-3 text-left font-semibold text-gray-600">Status</th>
            <th class="px-4 py-3 text-left font-semibold text-gray-600">Node</th>
            <th class="px-4 py-3 text-left font-semibold text-gray-600">Uptime</th>
            <th class="px-4 py-3 text-left font-semibold text-gray-600">Restarts</th>
            <th class="px-4 py-3 text-left font-semibold text-gray-600">CPU</th>
            <th class="px-4 py-3 text-left font-semibold text-gray-600">Memory</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-100">
          <tr
            v-for="(pod, index) in pods"
            :key="index"
            class="hover:bg-gray-50 transition duration-150"
          >
            <td class="px-4 py-3 font-medium text-gray-800">{{ pod.name }}</td>
            <td class="px-4 py-3 capitalize text-purple-700 font-semibold">
              {{ pod.namespace.replace('tenant-', '') }}
            </td>
            <td class="px-4 py-3 text-gray-600">{{ pod.namespace }}</td>
            <td class="px-4 py-3">
              <span
                :class="[
                  'px-2 py-1 rounded-full text-xs font-semibold',
                  pod.status === 'Running'
                    ? 'bg-green-100 text-green-700'
                    : pod.status === 'Pending'
                    ? 'bg-yellow-100 text-yellow-700'
                    : 'bg-red-100 text-red-700'
                ]"
              >
                {{ pod.status }}
              </span>
            </td>
            <td class="px-4 py-3 text-gray-700">{{ pod.node }}</td>
            <td class="px-4 py-3 text-gray-500" :title="pod.age">{{ pod.age }}</td>
            <td class="px-4 py-3 text-center">{{ pod.restarts }}</td>
            <td class="px-4 py-3 text-blue-600 font-mono">{{ pod.cpu }}</td>
            <td class="px-4 py-3 text-indigo-600 font-mono">{{ pod.memory }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-else-if="!loading && !error" class="text-gray-500 mt-6">
      No pods found for tenants.
    </div>
  </div>
</template>



<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { useAuth } from '../auth/useAuth.js'
import { getAccessToken } from '../auth/useAuth.js'

const { user } = useAuth()

const pods = ref([])
const loading = ref(false)
const error = ref(null)

async function fetchAll() {
  try {
    loading.value = true
    error.value = null
    const token = await getAccessToken()
    const res = await axios.get(`http://91.99.215.178:2001/admin/tenants/pods`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    pods.value = res.data
  } catch (err) {
    console.error('Failed to fetch pods:', err)
    error.value = err.message
  } finally {
    loading.value = false
  }
}

onMounted(fetchAll)
</script>

