<template>
  <div class="p-6">
    <h2 class="text-xl font-bold mb-4">Login</h2>
    <button @click="login" class="bg-blue-500 text-white px-4 py-2 rounded">Login with ZITADEL</button>
  </div>
</template>

<script setup>
import { login } from '../auth'
</script>
