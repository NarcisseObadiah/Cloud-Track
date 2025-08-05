<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100 dark:from-gray-900 dark:to-gray-800">
    <NavBar />
    <main class="max-w-7xl mx-auto py-8 px-4 sm:px-6 lg:px-8">
      <RouterView />
    </main>
    
    <!-- Footer -->
    <footer class="bg-white dark:bg-gray-800 border-t border-gray-200 dark:border-gray-700 mt-12">
      <div class="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
        <div class="flex flex-col md:flex-row justify-between items-center">
          <div class="flex items-center space-x-2 mb-4 md:mb-0">
            <div class="bg-gradient-to-r from-blue-600 to-purple-600 text-white p-1 rounded">
              <FontAwesomeIcon icon="database" class="text-sm" />
            </div>
            <span class="text-sm text-gray-600 dark:text-gray-400">
              © 2025 PaaS Platform. All rights reserved.
            </span>
          </div>
          <div class="flex items-center space-x-4 text-sm text-gray-500 dark:text-gray-400">
            <a href="#" class="hover:text-gray-700 dark:hover:text-gray-300 transition-colors">
              <FontAwesomeIcon icon="question-circle" class="mr-1" />
              Support
            </a>
            <a href="#" class="hover:text-gray-700 dark:hover:text-gray-300 transition-colors">
              <FontAwesomeIcon icon="book" class="mr-1" />
              Documentation
            </a>
            <a href="#" class="hover:text-gray-700 dark:hover:text-gray-300 transition-colors">
              <FontAwesomeIcon icon="shield-alt" class="mr-1" />
              Privacy
            </a>
          </div>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup>
import NavBar from './components/NavBar.vue'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
</script>

<style>
/* Custom scrollbar */
::-webkit-scrollbar {
  width: 8px;
}

::-webkit-scrollbar-track {
  background: #f1f1f1;
}

::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: #a1a1a1;
}

/* Smooth transitions for all elements */
* {
  transition: color 0.2s ease, background-color 0.2s ease, border-color 0.2s ease, opacity 0.2s ease, transform 0.2s ease;
}

/* Focus styles */
button:focus,
input:focus,
select:focus,
textarea:focus {
  outline: 2px solid #3b82f6;
  outline-offset: 2px;
}

/* Custom button hover effects */
.btn-hover {
  transform: translateY(0);
  transition: all 0.2s ease;
}

.btn-hover:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}
</style>
