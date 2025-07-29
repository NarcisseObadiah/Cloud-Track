<template>
  <div class="p-6 text-center text-lg font-bold text-blue-600">
    Redirecting...
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { handleCallback } from '../auth'
import { useRouter } from 'vue-router'

const router = useRouter()

onMounted(async () => {
  await handleCallback()
  router.push('/')
})
</script>
