<template>
  <div class="min-h-screen bg-gray-100 dark: bg-gray‑900">
    <NavBar />
    <main class="container mx-auto py-6">
      <RouterView />
    </main>
  </div>
</template>

<script setup>
import NavBar from './components/NavBar.vue'
</script>
