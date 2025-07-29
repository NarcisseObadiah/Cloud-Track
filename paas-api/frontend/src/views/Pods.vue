<template>
  <div class="p-6">
    <h2 class="text-xl font-bold mb-4">My Pods</h2>
    <ul v-if="pods.length">
      <li v-for="pod in pods" :key="pod.name">{{ pod.name }}</li>
    </ul>
    <p v-else>Loading...</p>
  </div>
</template>

<script setup>
import axios from 'axios'
import { ref, onMounted } from 'vue'
import { getUser } from '../auth'

const pods = ref([])

onMounted(async () => {
  const user = await getUser()
  const token = user?.access_token

  try {
    const res = await axios.get('http://localhost:8080/pods/my-namespace', {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
    pods.value = res.data
  } catch (err) {
    console.error('Failed to fetch pods:', err)
  }
})
</script>
