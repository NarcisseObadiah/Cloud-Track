<template>
  <div class="bg-white p-6 rounded shadow max-w-4xl mx-auto">
    <h2 class="text-2xl font-bold mb-6 flex items-center text-gray-800">
      <FontAwesomeIcon icon="cubes" class="mr-2 text-purple-600" />
      Tenant Pod Viewer
    </h2>

    <div class="flex flex-col sm:flex-row sm:items-center gap-4 mb-4">
      <input
        v-model="namespace"
        placeholder="Enter namespace (e.g., tenant-marhez)"
        class="flex-1 p-2 border border-gray-300 rounded text-sm shadow-sm focus:outline-none focus:ring-2 focus:ring-purple-500"
      />

      <button
        @click="fetchPods"
        class="bg-purple-600 hover:bg-purple-700 text-white py-2 px-4 rounded shadow flex items-center justify-center text-sm"
      >
        <FontAwesomeIcon icon="search" class="mr-2" />
        Get Pods
      </button>
    </div>

    <div v-if="loading" class="text-sm text-gray-500 mb-4 flex items-center">
      <FontAwesomeIcon icon="spinner" spin class="mr-2" /> Loading pods...
    </div>

    <div v-if="error" class="text-red-600 text-sm mb-4">
      {{ error }}
    </div>

    <div v-if="pods.length" class="overflow-x-auto">
      <table class="min-w-full table-auto border border-gray-200 rounded shadow text-sm">
        <thead class="bg-gray-100 text-gray-700">
          <tr>
            <th class="text-left px-4 py-2">Pod Name</th>
            <th class="text-left px-4 py-2">Namespace</th>
            <th class="text-left px-4 py-2">Status</th>
            <th class="text-left px-4 py-2">Age</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="(pod, index) in pods"
            :key="index"
            class="border-t hover:bg-gray-50"
          >
            <td class="px-4 py-2 font-mono">{{ pod.name }}</td>
            <td class="px-4 py-2">{{ pod.namespace }}</td>
            <td class="px-4 py-2">
              <span
                :class="[
                  'text-xs font-semibold px-2 py-1 rounded',
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
            <td class="px-4 py-2">{{ pod.age }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-else-if="!loading && !error" class="text-gray-500 text-sm mt-4">
      No pods found in this namespace.
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import axios from 'axios'
import { useAuth } from '../auth/useAuth.js'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

const { user } = useAuth()

const namespace = ref('')
const pods = ref([])
const loading = ref(false)
const error = ref(null)

async function fetchPods() {
  if (!namespace.value.trim()) {
    error.value = 'Namespace is required.'
    return
  }

  try {
    loading.value = true
    error.value = null
    pods.value = []

    const res = await axios.get(`http://91.99.215.178:2001/pods/${namespace.value}`, {
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
</script>
