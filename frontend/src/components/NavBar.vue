<template>
<nav class="bg-white dark:bg-gray-800 shadow p-4 flex justify-between items-center">
    <div class="space-x-6 hidden md:flex text-lg font-semibold">
      <RouterLink to="/" class="flex items-center hover:underline">
        <FontAwesomeIcon icon="home" class="mr-1" /> Dashboard
      </RouterLink>
      <RouterLink v-if="hasTenant" to="/databases" class="flex items-center hover:underline">
        <FontAwesomeIcon icon="database" class="mr-1" /> Databases
      </RouterLink>
      <RouterLink v-if="hasTenant" to="/pods" class="flex items-center hover:underline">
        <FontAwesomeIcon icon="cubes" class="mr-1" /> Pods
      </RouterLink>
      <RouterLink v-if="isAdmin" to="/admin" class="flex items-center hover:underline">
        <FontAwesomeIcon icon="user-shield" class="mr-1" /> Admin
      </RouterLink>
    </div>

    <div class="flex items-center space-x-4">
      <span v-if="user" class="text-gray-700 dark:text-gray-200">
        Hello, {{ user.profile.name }}
      </span>
      <button
        v-if="user"
        @click="logout"
        class="bg-red-600 hover:bg-red-700 text-white px-4 py-1 rounded flex items-center"
      >
        <FontAwesomeIcon icon="sign-out-alt" class="mr-1" /> Logout
      </button>
      <button
        v-else
        @click="login"
        class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-1 rounded flex items-center"
      >
        <FontAwesomeIcon icon="sign-in-alt" class="mr-1" /> Login
      </button>
    </div>
  </nav>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { useAuth } from '../auth/useAuth.js'

const { user, fetchUser, login, logout } = useAuth()
onMounted(fetchUser)

const roles = computed(() => {
  const rolesObj = user.value?.profile?.['urn:zitadel:iam:org:project:328396684147088791:roles'] || {}
  return Object.keys(rolesObj)
})

const hasTenant = computed(() => roles.value.includes('tenant'))
const isAdmin = computed(() => roles.value.includes('admin'))
</script>
