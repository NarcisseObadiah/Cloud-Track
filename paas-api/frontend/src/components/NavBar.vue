<script setup>
import { ref, onMounted, computed } from 'vue'
import { login, logout, getUser } from '../auth'
import { RouterLink } from 'vue-router'  // import RouterLink

const user = ref(null)

onMounted(async () => {
  user.value = await getUser()
})

const isAdmin = computed(() => {
  return user.value?.profile?.roles?.includes('admin')
})
</script>
ß
<template>
  <nav class="p-4 bg-gray-800 text-white flex justify-between items-center">
    <div class="space-x-4">
      <RouterLink to="/" class="hover:underline">Home</RouterLink>
      <RouterLink to="/databases" class="hover:underline">Databases</RouterLink>
      <RouterLink to="/pods" class="hover:underline">Pods</RouterLink>

      <RouterLink v-if="isAdmin" to="/admin" class="hover:underline">Admin</RouterLink>
    </div>
    <div>
      <span v-if="user" class="mr-4">{{ user.profile.name }}</span>
      <button v-if="user" @click="logout" class="bg-red-500 px-3 py-1 rounded">Logout</button>
      <button v-else @click="login" class="bg-green-500 px-3 py-1 rounded">Login</button>
    </div>
  </nav>
</template>
