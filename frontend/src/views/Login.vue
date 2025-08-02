<template>
  <div class="h-screen flex items-center justify-center">
    <button @click="login" class="bg-blue-600 text-white px-6 py-3 rounded">Login with Zitadel</button>
  </div>
</template>

<script setup>
import { useAuth } from '../auth/useAuth.js'
const { login } = useAuth()
</script>
