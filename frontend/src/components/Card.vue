<template>
  <div @click="go" class="p-4 bg-white dark:bg-gray‑800 rounded shadow hover:shadow-md cursor-pointer transition">
    <h3 class="text-lg font-semibold mb-1">{{ title }}</h3>
    <p class="text-gray-600 dark:text-gray‑400">Click to open</p>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
const props = defineProps({ title: String, link: String })
const router = useRouter()
function go() { router.push(link) }
</script>
