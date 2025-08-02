<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50">
    <div class="text-center">
      <p class="text-xl font-semibold mb-4">Logging you in...</p>
      <svg class="animate-spin h-8 w-8 mx-auto text-purple-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z"></path>
      </svg>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '../auth/useAuth.js'

const router = useRouter()
const { handleCallback } = useAuth()

onMounted(async () => {
  try {
    await handleCallback()
    router.replace('/')  // redirect to dashboard/home after login processed
  } catch (error) {
    console.error('Login callback error:', error)
  }
})
</script>
