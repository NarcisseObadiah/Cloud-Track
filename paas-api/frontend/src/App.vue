<script setup>
import NavBar from './components/NavBar.vue'
</script>

<template>
  <div class="min-h-screen bg-gray-100">
    <header>
      <NavBar />
    </header>

    <main class="container mx-auto mt-6">
      <RouterView />
    </main>
  </div>
</template>
