<template>
  <nav class="bg-white dark:bg-gray-800 shadow-lg border-b border-gray-200">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between h-16">
        <!-- Left side - Logo and Navigation -->
        <div class="flex items-center">
          <!-- Logo -->
          <div class="flex-shrink-0 flex items-center">
            <div class="bg-gradient-to-r from-blue-600 to-purple-600 text-white p-2 rounded-lg mr-3">
              <FontAwesomeIcon icon="database" class="text-xl" />
            </div>
            <span class="text-xl font-bold text-gray-800 dark:text-white">PaaS Platform</span>
          </div>

          <!-- Navigation Links -->
          <div class="hidden md:ml-8 md:flex md:space-x-1">
            <RouterLink 
              to="/" 
              class="nav-link flex items-center px-3 py-2 rounded-lg text-sm font-medium transition-colors duration-200"
              :class="$route.path === '/' ? 'nav-link-active' : 'nav-link-inactive'"
            >
              <FontAwesomeIcon icon="home" class="mr-2" />
              Dashboard
            </RouterLink>
            
            <RouterLink 
              v-if="hasTenant" 
              to="/databases" 
              class="nav-link flex items-center px-3 py-2 rounded-lg text-sm font-medium transition-colors duration-200"
              :class="$route.path === '/databases' ? 'nav-link-active' : 'nav-link-inactive'"
            >
              <FontAwesomeIcon icon="database" class="mr-2" />
              Databases
            </RouterLink>
            
            <RouterLink 
              v-if="hasTenant" 
              to="/pods" 
              class="nav-link flex items-center px-3 py-2 rounded-lg text-sm font-medium transition-colors duration-200"
              :class="$route.path === '/pods' ? 'nav-link-active' : 'nav-link-inactive'"
            >
              <FontAwesomeIcon icon="cubes" class="mr-2" />
              Pods
            </RouterLink>
            
            <RouterLink 
              v-if="isAdmin" 
              to="/admin" 
              class="nav-link flex items-center px-3 py-2 rounded-lg text-sm font-medium transition-colors duration-200"
              :class="$route.path === '/admin' ? 'nav-link-active' : 'nav-link-inactive'"
            >
              <FontAwesomeIcon icon="user-shield" class="mr-2" />
              Admin
            </RouterLink>
          </div>
        </div>

        <!-- Right side - User info and actions -->
        <div class="flex items-center space-x-4">
          <!-- User Info -->
          <div v-if="user" class="hidden sm:flex sm:items-center sm:space-x-3">
            <div class="flex items-center">
              <div class="bg-gradient-to-r from-blue-500 to-purple-500 text-white rounded-full w-8 h-8 flex items-center justify-center text-sm font-semibold">
                {{ userInitials }}
              </div>
              <div class="ml-3">
                <div class="text-sm font-medium text-gray-800 dark:text-gray-200">
                  {{ user.profile.name || user.profile.preferred_username }}
                </div>
                <div class="text-xs text-gray-500 dark:text-gray-400 flex items-center">
                  <FontAwesomeIcon icon="user" class="mr-1" />
                  {{ roleDisplay }}
                </div>
              </div>
            </div>
          </div>

          <!-- Action Buttons -->
          <div class="flex items-center space-x-2">
            <!-- Notifications (placeholder) -->
            <button 
              v-if="user"
              class="p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-lg transition-colors duration-200"
              title="Notifications"
            >
              <FontAwesomeIcon icon="bell" />
            </button>

            <!-- Theme Toggle (placeholder) -->
            <button 
              class="p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-lg transition-colors duration-200"
              title="Toggle Theme"
            >
              <FontAwesomeIcon icon="moon" />
            </button>

            <!-- Logout/Login Button -->
            <button
              v-if="user"
              @click="logout"
              class="bg-red-500 hover:bg-red-600 text-white px-4 py-2 rounded-lg text-sm font-medium flex items-center transition-colors duration-200 shadow-sm"
            >
              <FontAwesomeIcon icon="sign-out-alt" class="mr-2" />
              <span class="hidden sm:inline">Logout</span>
            </button>
            
            <button
              v-else
              @click="login"
              class="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-lg text-sm font-medium flex items-center transition-colors duration-200 shadow-sm"
            >
              <FontAwesomeIcon icon="sign-in-alt" class="mr-2" />
              <span class="hidden sm:inline">Login</span>
            </button>
          </div>

          <!-- Mobile menu button -->
          <div class="md:hidden">
            <button
              @click="mobileMenuOpen = !mobileMenuOpen"
              class="p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-lg transition-colors duration-200"
            >
              <FontAwesomeIcon :icon="mobileMenuOpen ? 'times' : 'bars'" />
            </button>
          </div>
        </div>
      </div>

      <!-- Mobile menu -->
      <div v-if="mobileMenuOpen" class="md:hidden border-t border-gray-200">
        <div class="px-2 pt-2 pb-3 space-y-1">
          <RouterLink 
            to="/" 
            @click="mobileMenuOpen = false"
            class="mobile-nav-link flex items-center px-3 py-2 rounded-lg text-base font-medium"
            :class="$route.path === '/' ? 'mobile-nav-link-active' : 'mobile-nav-link-inactive'"
          >
            <FontAwesomeIcon icon="home" class="mr-3" />
            Dashboard
          </RouterLink>
          
          <RouterLink 
            v-if="hasTenant" 
            to="/databases" 
            @click="mobileMenuOpen = false"
            class="mobile-nav-link flex items-center px-3 py-2 rounded-lg text-base font-medium"
            :class="$route.path === '/databases' ? 'mobile-nav-link-active' : 'mobile-nav-link-inactive'"
          >
            <FontAwesomeIcon icon="database" class="mr-3" />
            Databases
          </RouterLink>
          
          <RouterLink 
            v-if="hasTenant" 
            to="/pods" 
            @click="mobileMenuOpen = false"
            class="mobile-nav-link flex items-center px-3 py-2 rounded-lg text-base font-medium"
            :class="$route.path === '/pods' ? 'mobile-nav-link-active' : 'mobile-nav-link-inactive'"
          >
            <FontAwesomeIcon icon="cubes" class="mr-3" />
            Pods
          </RouterLink>
          
          <RouterLink 
            v-if="isAdmin" 
            to="/admin" 
            @click="mobileMenuOpen = false"
            class="mobile-nav-link flex items-center px-3 py-2 rounded-lg text-base font-medium"
            :class="$route.path === '/admin' ? 'mobile-nav-link-active' : 'mobile-nav-link-inactive'"
          >
            <FontAwesomeIcon icon="user-shield" class="mr-3" />
            Admin
          </RouterLink>
        </div>

        <!-- Mobile user info -->
        <div v-if="user" class="pt-4 pb-3 border-t border-gray-200">
          <div class="px-4 flex items-center">
            <div class="bg-gradient-to-r from-blue-500 to-purple-500 text-white rounded-full w-10 h-10 flex items-center justify-center text-sm font-semibold">
              {{ userInitials }}
            </div>
            <div class="ml-3">
              <div class="text-base font-medium text-gray-800 dark:text-gray-200">
                {{ user.profile.name || user.profile.preferred_username }}
              </div>
              <div class="text-sm text-gray-500 dark:text-gray-400">
                {{ roleDisplay }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </nav>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { useAuth } from '../auth/useAuth.js'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

const { user, fetchUser, login, logout } = useAuth()
const mobileMenuOpen = ref(false)

onMounted(fetchUser)

const roles = computed(() => {
  const rolesObj = user.value?.profile?.['urn:zitadel:iam:org:project:328396684147088791:roles'] || {}
  return Object.keys(rolesObj)
})

const hasTenant = computed(() => roles.value.includes('tenant'))
const isAdmin = computed(() => roles.value.includes('admin'))

const userInitials = computed(() => {
  if (!user.value?.profile) return 'U'
  const name = user.value.profile.name || user.value.profile.preferred_username || 'User'
  return name.split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2)
})

const roleDisplay = computed(() => {
  const rolesList = roles.value
  if (rolesList.includes('admin')) return 'Administrator'
  if (rolesList.includes('tenant')) return 'Tenant User'
  return 'User'
})

// Close mobile menu when clicking outside
document.addEventListener('click', (e) => {
  if (!e.target.closest('nav')) {
    mobileMenuOpen.value = false
  }
})
</script>

<style scoped>
.nav-link-active {
  @apply bg-blue-100 text-blue-700 dark:bg-blue-900 dark:text-blue-200;
}

.nav-link-inactive {
  @apply text-gray-600 hover:text-gray-900 hover:bg-gray-100 dark:text-gray-300 dark:hover:text-white dark:hover:bg-gray-700;
}

.mobile-nav-link-active {
  @apply bg-blue-100 text-blue-700 dark:bg-blue-900 dark:text-blue-200;
}

.mobile-nav-link-inactive {
  @apply text-gray-600 hover:text-gray-900 hover:bg-gray-100 dark:text-gray-300 dark:hover:text-white dark:hover:bg-gray-700;
}
</style>
