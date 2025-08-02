<template>
  <div class="bg-white p-8 rounded-xl shadow-lg max-w-lg mx-auto mt-10 transition-all">
    <h2 class="text-3xl font-extrabold mb-6 flex items-center text-gray-800">
      <FontAwesomeIcon icon="database" class="mr-3 text-blue-600" />
      Create a New Database
    </h2>

    <form @submit.prevent="create" class="space-y-5">
      <!-- Username -->
      <div>
        <label for="username" class="block text-sm font-medium text-gray-700 mb-1">Username</label>
        <input
          v-model="username"
          id="username"
          type="text"
          placeholder="e.g., francis"
          class="w-full px-4 py-2 border border-gray-300 rounded-lg shadow-sm focus:ring-blue-500 focus:border-blue-500"
          required
        />
      </div>

      <!-- Database Name -->
      <div>
        <label for="dbName" class="block text-sm font-medium text-gray-700 mb-1">Database Name</label>
        <input
          v-model="dbName"
          id="dbName"
          type="text"
          placeholder="e.g., francisdb"
          class="w-full px-4 py-2 border border-gray-300 rounded-lg shadow-sm focus:ring-blue-500 focus:border-blue-500"
          required
        />
      </div>


      <!-- Submit Button -->
      <button
        type="submit"
        class="w-full bg-blue-600 hover:bg-blue-700 text-white font-semibold py-2.5 rounded-lg flex items-center justify-center transition duration-150 ease-in-out"
      >
        <FontAwesomeIcon icon="plus-circle" class="mr-2" />
        Create Database
      </button>
    </form>

    <transition name="fade">
      <div
        v-if="msg"
        :class="['mt-6 px-4 py-3 rounded-lg text-sm font-medium shadow-sm', msgClass]"
      >
        {{ msg }}
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import axios from 'axios'
import { useAuth } from '../auth/useAuth.js'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

const { user } = useAuth()

const username = ref('')
const dbName = ref('')
const password = ref('')
const showPassword = ref(false)

const msg = ref('')
const msgClass = ref('')

function togglePassword() {
  showPassword.value = !showPassword.value
}

async function create() {
  if (!username.value || !dbName.value || !password.value) {
    msg.value = 'All fields are required.'
    msgClass.value = 'bg-yellow-100 text-yellow-800'
    return
  }

  try {
    await axios.post('http://91.99.215.178:2001/databases', {
      username: username.value,
      db_name: dbName.value,
      password: password.value
    }, {
      headers: {
        Authorization: `Bearer ${user.value.access_token}`
      }
    })
    msg.value = `✅ Database "${dbName.value}" created successfully!`
    msgClass.value = 'bg-green-100 text-green-800'

    // Reset form
    username.value = ''
    dbName.value = ''
    password.value = ''
  } catch (err) {
    msg.value = err.response?.data?.error || err.message || 'Failed to create database.'
    msgClass.value = 'bg-red-100 text-red-800'
  }
}
</script>

<style scoped>
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.3s ease;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}
</style>






