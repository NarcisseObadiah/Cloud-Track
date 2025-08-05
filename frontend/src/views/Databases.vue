<template>
  <div class="space-y-8">
    <!-- Header -->
    <div class="bg-gradient-to-r from-blue-600 to-indigo-600 rounded-2xl shadow-xl p-8 text-white">
      <h1 class="text-4xl font-bold mb-2 flex items-center">
        <FontAwesomeIcon icon="database" class="mr-4" />
        Database Management
      </h1>
      <p class="text-blue-100 text-lg">Create and manage your PostgreSQL databases with ease</p>
    </div>

    <!-- Quick Stats -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
      <div class="bg-white rounded-xl shadow-lg p-6">
        <div class="flex items-center">
          <div class="bg-green-500 rounded-lg p-3 text-white mr-4">
            <FontAwesomeIcon icon="database" class="text-xl" />
          </div>
          <div>
            <div class="text-2xl font-bold text-gray-800">{{ clusterSummary.total || 0 }}</div>
            <div class="text-gray-600 text-sm">Total Clusters</div>
          </div>
        </div>
      </div>
      
      <div class="bg-white rounded-xl shadow-lg p-6">
        <div class="flex items-center">
          <div class="bg-blue-500 rounded-lg p-3 text-white mr-4">
            <FontAwesomeIcon icon="check-circle" class="text-xl" />
          </div>
          <div>
            <div class="text-2xl font-bold text-gray-800">{{ runningDatabases }}</div>
            <div class="text-gray-600 text-sm">Ready</div>
          </div>
        </div>
      </div>
      
      <div class="bg-white rounded-xl shadow-lg p-6">
        <div class="flex items-center">
          <div class="bg-yellow-500 rounded-lg p-3 text-white mr-4">
            <FontAwesomeIcon icon="clock" class="text-xl" />
          </div>
          <div>
            <div class="text-2xl font-bold text-gray-800">{{ pendingDatabases }}</div>
            <div class="text-gray-600 text-sm">Creating</div>
          </div>
        </div>
      </div>
      
      <div class="bg-white rounded-xl shadow-lg p-6">
        <div class="flex items-center">
          <div class="bg-red-500 rounded-lg p-3 text-white mr-4">
            <FontAwesomeIcon icon="exclamation-triangle" class="text-xl" />
          </div>
          <div>
            <div class="text-2xl font-bold text-gray-800">{{ failedDatabases }}</div>
            <div class="text-gray-600 text-sm">Failed</div>
          </div>
        </div>
      </div>
    </div>

    <!-- Create Database Section -->
    <div class="bg-white rounded-xl shadow-lg p-8">
      <h2 class="text-2xl font-bold mb-6 flex items-center text-gray-800">
        <FontAwesomeIcon icon="plus-circle" class="mr-3 text-blue-600" />
        Create New Database
      </h2>

      <!-- Display current user -->
      <div class="mb-6 p-4 bg-blue-50 rounded-lg border border-blue-200">
        <p class="text-sm text-blue-800 flex items-center">
          <FontAwesomeIcon icon="user" class="mr-2" />
          Creating database for tenant: <strong class="ml-1">{{ username }}</strong>
        </p>
      </div>

      <form @submit.prevent="create" class="space-y-6">
        <!-- Database Name (Optional) -->
        <div>
          <label for="dbName" class="block text-sm font-medium text-gray-700 mb-2">
            Database Name 
            <span class="text-xs text-gray-500">(optional - will auto-generate if empty)</span>
          </label>
          <input
            v-model="dbName"
            id="dbName"
            type="text"
            placeholder="e.g., my-app-db (leave empty for auto-generation)"
            class="w-full px-4 py-3 border border-gray-300 rounded-lg shadow-sm focus:ring-blue-500 focus:border-blue-500 transition-colors"
          />
          <p class="text-xs text-gray-500 mt-1">
            If empty, will create: <code class="bg-gray-100 px-2 py-1 rounded">{{ username }}-db</code>
          </p>
        </div>

        <!-- Replicas (Optional) -->
        <div>
          <label for="replicas" class="block text-sm font-medium text-gray-700 mb-2">
            High Availability 
            <span class="text-xs text-gray-500">(optional)</span>
          </label>
          <select
            v-model="replicas"
            id="replicas"
            class="w-full px-4 py-3 border border-gray-300 rounded-lg shadow-sm focus:ring-blue-500 focus:border-blue-500 transition-colors"
          >
            <option value="1">Single Instance (Default)</option>
            <option value="2">2 Replicas</option>
            <option value="3">3 Replicas (Recommended for Production)</option>
          </select>
        </div>

        <!-- Submit Button -->
        <button
          type="submit"
          :disabled="creating"
          class="w-full bg-blue-600 hover:bg-blue-700 disabled:opacity-50 text-white font-semibold py-3 rounded-lg flex items-center justify-center transition duration-150 ease-in-out"
        >
          <FontAwesomeIcon :icon="creating ? 'spinner' : 'plus-circle'" :class="creating ? 'animate-spin' : ''" class="mr-2" />
          {{ creating ? 'Creating Database (2-5 min)...' : 'Create Database' }}
        </button>
      </form>
    </div>

    <!-- Success message with credentials -->
    <transition name="fade">
      <div v-if="credentials" class="bg-green-50 border border-green-200 rounded-xl p-6">
        <h3 class="text-xl font-semibold text-green-800 mb-4 flex items-center">
          <FontAwesomeIcon icon="check-circle" class="mr-3" />
          Database Created Successfully!
        </h3>
        
        <div class="space-y-4">
          <div class="bg-white p-4 rounded-lg border">
            <h4 class="font-medium text-gray-700 mb-3">Connection Details:</h4>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
              <div>
                <strong class="text-gray-600">Database:</strong>
                <div class="font-mono bg-gray-100 p-2 rounded mt-1">{{ credentials.database_name }}</div>
              </div>
              <div>
                <strong class="text-gray-600">Host:</strong>
                <div class="font-mono bg-gray-100 p-2 rounded mt-1 text-xs">{{ credentials.host }}</div>
              </div>
              <div>
                <strong class="text-gray-600">Port:</strong>
                <div class="font-mono bg-gray-100 p-2 rounded mt-1">{{ credentials.port }}</div>
              </div>
              <div>
                <strong class="text-gray-600">Username:</strong>
                <div class="font-mono bg-gray-100 p-2 rounded mt-1">{{ credentials.primary_user?.username }}</div>
              </div>
              <div class="md:col-span-2">
                <strong class="text-gray-600">Password:</strong>
                <div class="flex items-center mt-1">
                  <div class="font-mono bg-gray-100 p-2 rounded flex-1" :class="showCredPassword ? '' : 'blur-sm'">
                    {{ credentials.primary_user?.password }}
                  </div>
                  <button @click="showCredPassword = !showCredPassword" class="ml-2 text-blue-600 text-sm px-2 py-1 hover:bg-blue-50 rounded">
                    {{ showCredPassword ? 'Hide' : 'Show' }}
                  </button>
                </div>
              </div>
            </div>
          </div>
          
          <div class="bg-white p-4 rounded-lg border">
            <h4 class="font-medium text-gray-700 mb-2">Ready-to-use Connection String:</h4>
            <div class="bg-gray-100 p-3 rounded font-mono text-xs break-all">
              {{ credentials.connection_string }}
            </div>
            <button 
              @click="copyToClipboard(credentials.connection_string)"
              class="mt-2 text-sm bg-blue-600 text-white px-3 py-1 rounded hover:bg-blue-700 transition-colors"
            >
              <FontAwesomeIcon icon="key" class="mr-1" />
              Copy Connection String
            </button>
          </div>
          
          <!-- External Connection Instructions -->
          <div class="bg-blue-50 p-4 rounded-lg border border-blue-200">
            <h4 class="font-medium text-blue-800 mb-2 flex items-center">
              <FontAwesomeIcon icon="info-circle" class="mr-2" />
              External Connection (from outside Kubernetes)
            </h4>
            <p class="text-sm text-blue-700 mb-2">To connect from outside the cluster, use port forwarding:</p>
            <div class="bg-white p-2 rounded font-mono text-xs border">
              kubectl port-forward -n {{ 'tenant-' + username }} pod/{{ credentials.database_name }}-0 5432:5432
            </div>
            <button 
              @click="copyToClipboard(`kubectl port-forward -n tenant-${username} pod/${credentials.database_name}-0 5432:5432`)"
              class="mt-2 text-sm bg-blue-600 text-white px-3 py-1 rounded hover:bg-blue-700 transition-colors"
            >
              <FontAwesomeIcon icon="terminal" class="mr-1" />
              Copy Port Forward Command
            </button>
            <p class="text-xs text-gray-600 mt-2">
              <FontAwesomeIcon icon="info-circle" class="mr-1" />
              Using pod/{{ credentials.database_name }}-0 because Zalando creates StatefulSets with pods
            </p>
            <p class="text-xs text-blue-600 mt-1">Then connect to: localhost:5432</p>
          </div>
          
          <!-- Connection Examples -->
          <div class="bg-gray-50 p-4 rounded-lg border">
            <h4 class="font-medium text-gray-700 mb-3">Connection Examples:</h4>
            <div class="space-y-3">
              <div>
                <h5 class="text-sm font-medium text-gray-600">psql (PostgreSQL CLI):</h5>
                <div class="bg-white p-2 rounded font-mono text-xs border mt-1">
                  psql "{{ credentials.connection_string }}"
                </div>
              </div>
              
              <div>
                <h5 class="text-sm font-medium text-gray-600">Python (psycopg2):</h5>
                <div class="bg-white p-2 rounded font-mono text-xs border mt-1">
                  conn = psycopg2.connect("{{ credentials.connection_string }}")
                </div>
              </div>
              
              <div>
                <h5 class="text-sm font-medium text-gray-600">Node.js (pg):</h5>
                <div class="bg-white p-2 rounded font-mono text-xs border mt-1">
                  const client = new Client({ connectionString: "{{ credentials.connection_string }}" })
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <button 
          @click="credentials = null; loadMyDatabases()"
          class="mt-4 text-green-600 hover:text-green-800 text-sm font-medium"
        >
          <FontAwesomeIcon icon="times" class="mr-1" />
          Dismiss and Refresh List
        </button>
      </div>
    </transition>

    <!-- Error or info message -->
    <transition name="fade">
      <div
        v-if="msg && !credentials"
        :class="['p-4 rounded-lg text-sm font-medium shadow-sm', msgClass]"
      >
        {{ msg }}
      </div>
    </transition>

    <!-- My Databases List -->
    <div class="bg-white rounded-xl shadow-lg overflow-hidden">
      <div class="px-6 py-4 border-b border-gray-200 bg-gray-50">
        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <h2 class="text-xl font-semibold text-gray-800 flex items-center">
            <FontAwesomeIcon icon="database" class="mr-3 text-blue-600" />
            My Databases
          </h2>
          
          <button
            @click="loadMyDatabases"
            :disabled="loadingDatabases"
            class="bg-blue-600 hover:bg-blue-700 disabled:opacity-50 text-white px-4 py-2 rounded-lg text-sm flex items-center transition-colors"
          >
            <FontAwesomeIcon :icon="loadingDatabases ? 'spinner' : 'sync-alt'" :class="loadingDatabases ? 'animate-spin' : ''" class="mr-2" />
            Refresh
          </button>
        </div>
      </div>

      <div v-if="loadingDatabases" class="flex justify-center py-12">
        <div class="text-center">
          <FontAwesomeIcon icon="spinner" class="text-4xl text-blue-600 animate-spin mb-4" />
          <p class="text-gray-600">Loading your databases...</p>
        </div>
      </div>

      <div v-else-if="myDatabases.length === 0" class="text-center py-12">
        <FontAwesomeIcon icon="database" class="text-4xl text-gray-400 mb-4" />
        <p class="text-gray-600 text-lg mb-2">No databases found</p>
        <p class="text-gray-500 text-sm">Create your first database using the form above.</p>
      </div>

      <div v-else class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Database Cluster
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Status & Health
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Resources
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Connection
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Actions
              </th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr
              v-for="(db, index) in myDatabases"
              :key="index"
              class="hover:bg-gray-50 transition-colors"
            >
              <td class="px-6 py-4">
                <div class="flex items-center">
                  <div class="flex-shrink-0 mr-4">
                    <div :class="[
                      'w-3 h-3 rounded-full',
                      db.status === 'Ready' ? 'bg-green-500' : 
                      db.status === 'Creating' || db.status === 'Pending' || db.status === 'Pod Running' || db.status === 'Credentials Ready' ? 'bg-yellow-500' : 
                      'bg-red-500'
                    ]"></div>
                  </div>
                  <div>
                    <div class="text-sm font-medium text-gray-900 flex items-center">
                      {{ db.name }}
                      <span v-if="db.creation_method === 'manual'" class="ml-2 inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-orange-100 text-orange-800">
                        Manual
                      </span>
                      <span v-else-if="db.creation_method === 'zalando'" class="ml-2 inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-blue-100 text-blue-800">
                        Zalando
                      </span>
                    </div>
                    <div class="text-xs text-gray-500">{{ db.namespace }}</div>
                    <div v-if="db.postgres_version" class="text-xs text-gray-400">{{ db.postgres_version }}</div>
                  </div>
                </div>
              </td>
              
              <td class="px-6 py-4">
                <div class="space-y-1">
                  <span :class="[
                    'inline-flex px-2 py-1 text-xs font-semibold rounded-full',
                    db.status === 'Ready' ? 'bg-green-100 text-green-800' :
                    db.status === 'Creating' || db.status === 'Pending' ? 'bg-yellow-100 text-yellow-800' :
                    db.status === 'Pod Running' ? 'bg-blue-100 text-blue-800' :
                    db.status === 'Credentials Ready' ? 'bg-purple-100 text-purple-800' :
                    'bg-red-100 text-red-800'
                  ]">
                    {{ db.status }}
                  </span>
                  <div class="text-xs text-gray-500">
                    {{ db.detailed_status }}
                  </div>
                  <div v-if="db.replicas" class="text-xs text-gray-400">
                    {{ db.running_replicas || 0 }}/{{ db.replicas }} pods running
                  </div>
                </div>
              </td>
              
              <td class="px-6 py-4 text-sm text-gray-600">
                <div class="space-y-1">
                  <div v-if="db.storage_size">Storage: {{ db.storage_size }}</div>
                  <div v-if="db.created_at">Created: {{ new Date(db.created_at).toLocaleDateString() }}</div>
                  <div class="flex space-x-2">
                    <span v-if="db.credentials_ready" class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-green-100 text-green-800">
                      <FontAwesomeIcon icon="key" class="mr-1" />
                      Credentials
                    </span>
                    <span v-if="db.connection_ready" class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-blue-100 text-blue-800">
                      <FontAwesomeIcon icon="link" class="mr-1" />
                      Connected
                    </span>
                  </div>
                </div>
              </td>
              
              <td class="px-6 py-4 text-sm text-gray-600">
                <div v-if="db.connection_ready && db.connection_info" class="space-y-1">
                  <div class="text-xs font-mono bg-gray-100 px-2 py-1 rounded">
                    {{ db.connection_info.host }}:{{ db.connection_info.port }}
                  </div>
                  <div class="text-xs text-gray-500">{{ db.connection_info.database }}</div>
                </div>
                <div v-else class="text-xs text-gray-400">
                  {{ db.credentials_ready ? 'Testing connection...' : 'Waiting for credentials...' }}
                </div>
              </td>
              
              <td class="px-6 py-4 text-sm">
                <div class="flex space-x-2">
                  <button
                    @click="checkDatabaseStatus(db.name)"
                    class="text-blue-600 hover:text-blue-800 transition-colors"
                    title="Check Status"
                  >
                    <FontAwesomeIcon icon="sync-alt" />
                  </button>
                  <button
                    @click="viewDatabaseDetails(db)"
                    class="text-green-600 hover:text-green-800 transition-colors"
                    title="View Details"
                  >
                    <FontAwesomeIcon icon="eye" />
                  </button>
                  <button
                    @click="confirmDelete(db)"
                    class="text-red-600 hover:text-red-800 transition-colors"
                    title="Delete Database"
                  >
                    <FontAwesomeIcon icon="trash" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <div v-if="deleteTarget" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div class="bg-white rounded-xl shadow-2xl max-w-md w-full">
        <div class="p-6">
          <div class="flex items-center mb-4">
            <FontAwesomeIcon icon="exclamation-triangle" class="text-red-500 text-2xl mr-3" />
            <h3 class="text-xl font-semibold text-gray-800">Confirm Deletion</h3>
          </div>
          <p class="text-gray-600 mb-6">
            Are you sure you want to delete the database <strong>{{ deleteTarget.name }}</strong>? 
            This action cannot be undone and all data will be lost.
          </p>
          <div class="flex space-x-3">
            <button
              @click="deleteDatabase"
              :disabled="deleting"
              class="flex-1 bg-red-600 hover:bg-red-700 disabled:opacity-50 text-white px-4 py-2 rounded-lg font-medium transition-colors"
            >
              <FontAwesomeIcon :icon="deleting ? 'spinner' : 'trash'" :class="deleting ? 'animate-spin' : ''" class="mr-2" />
              {{ deleting ? 'Deleting...' : 'Delete' }}
            </button>
            <button
              @click="deleteTarget = null"
              :disabled="deleting"
              class="flex-1 bg-gray-300 hover:bg-gray-400 disabled:opacity-50 text-gray-700 px-4 py-2 rounded-lg font-medium transition-colors"
            >
              Cancel
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Database Details Modal -->
    <div v-if="selectedDatabase" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div class="bg-white rounded-xl shadow-2xl max-w-4xl w-full max-h-[90vh] overflow-y-auto">
        <div class="flex items-center justify-between p-6 border-b border-gray-200">
          <h3 class="text-xl font-semibold text-gray-800 flex items-center">
            <FontAwesomeIcon icon="database" class="mr-3 text-blue-600" />
            Database Cluster Details
          </h3>
          <button
            @click="selectedDatabase = null"
            class="text-gray-400 hover:text-gray-600 transition-colors"
          >
            <FontAwesomeIcon icon="times" class="text-xl" />
          </button>
        </div>
        <div class="p-6 space-y-6">
          <!-- Overview -->
          <div class="bg-gray-50 rounded-lg p-4">
            <h4 class="text-lg font-medium text-gray-800 mb-3">Overview</h4>
            <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div>
                <label class="text-sm font-medium text-gray-700">Cluster Name</label> 
                <div class="mt-1 font-mono text-sm p-3 bg-white rounded border">{{ selectedDatabase.name }}</div>
              </div>
              <div>
                <label class="text-sm font-medium text-gray-700">Namespace</label>
                <div class="mt-1 text-sm p-3 bg-white rounded border">{{ selectedDatabase.namespace }}</div>
              </div>
              <div>
                <label class="text-sm font-medium text-gray-700">PostgreSQL Version</label>
                <div class="mt-1 text-sm p-3 bg-white rounded border">{{ selectedDatabase.postgres_version || 'N/A' }}</div>
              </div>
            </div>
          </div>

          <!-- Status & Health -->
          <div class="bg-gray-50 rounded-lg p-4">
            <h4 class="text-lg font-medium text-gray-800 mb-3">Status & Health</h4>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label class="text-sm font-medium text-gray-700">Current Status</label>
                <div class="mt-1">
                  <span :class="[
                    'inline-flex px-3 py-1 text-sm font-semibold rounded-full',
                    selectedDatabase.status === 'Ready' ? 'bg-green-100 text-green-800' :
                    selectedDatabase.status === 'Creating' || selectedDatabase.status === 'Pending' ? 'bg-yellow-100 text-yellow-800' :
                    selectedDatabase.status === 'Pod Running' ? 'bg-blue-100 text-blue-800' :
                    selectedDatabase.status === 'Credentials Ready' ? 'bg-purple-100 text-purple-800' :
                    'bg-red-100 text-red-800'
                  ]">
                    {{ selectedDatabase.status }}
                  </span>
                </div>
              </div>
              <div>
                <label class="text-sm font-medium text-gray-700">Detailed Status</label>
                <div class="mt-1 text-sm p-3 bg-white rounded border">{{ selectedDatabase.detailed_status }}</div>
              </div>
              <div>
                <label class="text-sm font-medium text-gray-700">Creation Method</label>
                <div class="mt-1">
                  <span :class="[
                    'inline-flex px-3 py-1 text-sm font-semibold rounded-full',
                    selectedDatabase.creation_method === 'zalando' ? 'bg-blue-100 text-blue-800' :
                    selectedDatabase.creation_method === 'manual' ? 'bg-orange-100 text-orange-800' :
                    'bg-gray-100 text-gray-800'
                  ]">
                    {{ selectedDatabase.creation_method || 'Unknown' }}
                  </span>
                </div>
              </div>
              <div>
                <label class="text-sm font-medium text-gray-700">Created</label>
                <div class="mt-1 text-sm p-3 bg-white rounded border">
                  {{ selectedDatabase.created_at ? new Date(selectedDatabase.created_at).toLocaleString() : 'N/A' }}
                </div>
              </div>
            </div>
          </div>

          <!-- Resources -->
          <div class="bg-gray-50 rounded-lg p-4">
            <h4 class="text-lg font-medium text-gray-800 mb-3">Resources</h4>
            <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div>
                <label class="text-sm font-medium text-gray-700">Replicas</label>
                <div class="mt-1 text-sm p-3 bg-white rounded border">
                  {{ selectedDatabase.running_replicas || 0 }} / {{ selectedDatabase.replicas || 1 }} running
                </div>
              </div>
              <div>
                <label class="text-sm font-medium text-gray-700">Storage</label>
                <div class="mt-1 text-sm p-3 bg-white rounded border">{{ selectedDatabase.storage_size || 'N/A' }}</div>
              </div>
              <div>
                <label class="text-sm font-medium text-gray-700">Pod Status</label>
                <div class="mt-1 text-sm p-3 bg-white rounded border">{{ selectedDatabase.pod_status || 'N/A' }}</div>
              </div>
            </div>
          </div>

          <!-- Connection Information -->
          <div v-if="selectedDatabase.connection_ready && selectedDatabase.connection_info" class="bg-green-50 rounded-lg p-4 border border-green-200">
            <h4 class="text-lg font-medium text-green-800 mb-3 flex items-center">
              <FontAwesomeIcon icon="link" class="mr-2" />
              Connection Information
            </h4>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label class="text-sm font-medium text-green-700">Host</label>
                <div class="mt-1 font-mono text-sm p-3 bg-white rounded border">{{ selectedDatabase.connection_info.host }}</div>
              </div>
              <div>
                <label class="text-sm font-medium text-green-700">Port</label>
                <div class="mt-1 font-mono text-sm p-3 bg-white rounded border">{{ selectedDatabase.connection_info.port }}</div>
              </div>
              <div>
                <label class="text-sm font-medium text-green-700">Database</label>
                <div class="mt-1 font-mono text-sm p-3 bg-white rounded border">{{ selectedDatabase.connection_info.database }}</div>
              </div>
              <div>
                <label class="text-sm font-medium text-green-700">Username</label>
                <div class="mt-1 font-mono text-sm p-3 bg-white rounded border">{{ selectedDatabase.connection_info.username }}</div>
              </div>
            </div>
            
            <!-- Port Forward Command -->
            <div class="mt-4">
              <label class="text-sm font-medium text-green-700">External Access Command</label>
              <div class="mt-1 bg-white p-3 rounded border font-mono text-xs">
                kubectl port-forward -n {{ selectedDatabase.namespace }} pod/{{ selectedDatabase.name }}-0 5432:5432
              </div>
              <button 
                @click="copyToClipboard(`kubectl port-forward -n ${selectedDatabase.namespace} pod/${selectedDatabase.name}-0 5432:5432`)"
                class="mt-2 text-sm bg-green-600 text-white px-3 py-1 rounded hover:bg-green-700 transition-colors"
              >
                <FontAwesomeIcon icon="terminal" class="mr-1" />
                Copy Port Forward Command
              </button>
            </div>
          </div>

          <!-- Credentials Status -->
          <div class="bg-gray-50 rounded-lg p-4">
            <h4 class="text-lg font-medium text-gray-800 mb-3">Credentials Status</h4>
            <div class="flex space-x-4">
              <div class="flex items-center">
                <span :class="[
                  'inline-flex items-center px-3 py-1 text-sm font-medium rounded-full',
                  selectedDatabase.credentials_ready ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
                ]">
                  <FontAwesomeIcon :icon="selectedDatabase.credentials_ready ? 'check-circle' : 'times-circle'" class="mr-1" />
                  {{ selectedDatabase.credentials_ready ? 'Credentials Available' : 'Credentials Pending' }}
                </span>
              </div>
              <div class="flex items-center">
                <span :class="[
                  'inline-flex items-center px-3 py-1 text-sm font-medium rounded-full',
                  selectedDatabase.connection_ready ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
                ]">
                  <FontAwesomeIcon :icon="selectedDatabase.connection_ready ? 'link' : 'unlink'" class="mr-1" />
                  {{ selectedDatabase.connection_ready ? 'Connection Ready' : 'Connection Not Ready' }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import axios from 'axios'
import { useAuth } from '../auth/useAuth.js'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import { API_BASE_URL, API_ENDPOINTS } from '../config/api.js'

const { user } = useAuth()

// Form data
const dbName = ref('')
const replicas = ref(1)

// State
const creating = ref(false)
const loadingDatabases = ref(false)
const deleting = ref(false)
const msg = ref('')
const msgClass = ref('')
const credentials = ref(null)
const showCredPassword = ref(false)
const myDatabases = ref([])
const clusterSummary = ref({})
const deleteTarget = ref(null)
const selectedDatabase = ref(null)

// Computed properties
const username = computed(() => {
  const raw = user.value?.profile?.preferred_username || user.value?.profile?.email || 'user'
  const base = raw.split('@')[0]
  return base.toLowerCase().replace(/[^a-z0-9-]/g, '')
})

const runningDatabases = computed(() => {
  return clusterSummary.value.ready || 0
})

const totalReplicas = computed(() => {
  return myDatabases.value.reduce((sum, db) => sum + (db.replicas || 1), 0)
})

const pendingDatabases = computed(() => {
  return clusterSummary.value.creating || 0
})

const failedDatabases = computed(() => {
  return clusterSummary.value.failed || 0
})

// Methods
function copyToClipboard(text) {
  navigator.clipboard.writeText(text).then(() => {
    msg.value = 'Connection string copied to clipboard!'
    msgClass.value = 'bg-blue-100 text-blue-800'
    setTimeout(() => { msg.value = '' }, 3000)
  }).catch(() => {
    msg.value = 'Failed to copy to clipboard'
    msgClass.value = 'bg-red-100 text-red-800'
  })
}

async function create() {
  credentials.value = null
  msg.value = ''
  creating.value = true

  try {
    const payload = {
      username: username.value
    }

    if (dbName.value.trim()) {
      payload.db_name = dbName.value.trim()
    }
    
    if (replicas.value > 1) {
      payload.replicas = parseInt(replicas.value)
    }

    console.log('Creating database with payload:', payload)
    
    // Show initial message
    msg.value = 'Creating database instance... This may take 2-5 minutes.'
    msgClass.value = 'bg-blue-100 text-blue-800'

    const response = await axios.post(`${API_BASE_URL}${API_ENDPOINTS.DATABASES}`, payload, {
      headers: {
        Authorization: `Bearer ${user.value.access_token}`
      },
      timeout: 6 * 60 * 1000 // 6 minutes timeout
    })

    console.log('Database creation response:', response.data)
    
    if (response.data && response.data.credentials) {
      credentials.value = response.data.credentials
      msg.value = response.data.message || 'Database created successfully with credentials!'
      msgClass.value = 'bg-green-100 text-green-800'
    } else {
      console.error('No credentials in response:', response.data)
      msg.value = 'Database created but no credentials returned'
      msgClass.value = 'bg-yellow-100 text-yellow-800'
    }
    
    dbName.value = ''
    replicas.value = 1
    
  } catch (err) {
    console.error('Database creation error:', err)
    console.error('Error response:', err.response?.data)
    
    if (err.code === 'ECONNABORTED') {
      msg.value = 'Request timed out. Database may still be creating in the background.'
      msgClass.value = 'bg-yellow-100 text-yellow-800'
    } else {
      msg.value = err.response?.data?.error || err.response?.data?.message || err.message || 'Failed to create database.'
      msgClass.value = 'bg-red-100 text-red-800'
    }
    credentials.value = null
  } finally {
    creating.value = false
  }
}

async function loadMyDatabases() {
  loadingDatabases.value = true
  try {
    const response = await axios.get(`${API_BASE_URL}/databases/${username.value}`, {
      headers: {
        Authorization: `Bearer ${user.value.access_token}`
      }
    })
    
    console.log('Database clusters response:', response.data)
    myDatabases.value = response.data.clusters || []
    clusterSummary.value = response.data.summary || {}
    
  } catch (err) {
    console.error('Failed to load database clusters:', err)
    myDatabases.value = []
    clusterSummary.value = {}
  } finally {
    loadingDatabases.value = false
  }
}

async function checkDatabaseStatus(dbName) {
  try {
    const response = await axios.get(`${API_BASE_URL}/databases/${username.value}/${dbName}`, {
      headers: {
        Authorization: `Bearer ${user.value.access_token}`
      }
    })
    
    const cluster = response.data.cluster
    msg.value = `Database ${dbName} status: ${cluster.status} - ${cluster.detailed_status}`
    msgClass.value = 'bg-blue-100 text-blue-800'
    setTimeout(() => { msg.value = '' }, 5000)
    
    // Refresh the list
    loadMyDatabases()
  } catch (err) {
    msg.value = err.response?.data?.error || 'Failed to check database status'
    msgClass.value = 'bg-red-100 text-red-800'
  }
}

function viewDatabaseDetails(database) {
  selectedDatabase.value = database
}

function confirmDelete(database) {
  deleteTarget.value = database
}

async function deleteDatabase() {
  if (!deleteTarget.value) return
  
  console.log('Deleting database:', deleteTarget.value.name)
  deleting.value = true
  
  try {
    const payload = {
      username: username.value,
      db_name: deleteTarget.value.name
    }
    
    console.log('Delete payload:', payload)
    
    const response = await axios.delete(`${API_BASE_URL}${API_ENDPOINTS.DATABASES}`, {
      data: payload,
      headers: {
        Authorization: `Bearer ${user.value.access_token}`
      }
    })
    
    console.log('Delete response:', response.data)
    
    msg.value = `Database ${deleteTarget.value.name} deleted successfully`
    msgClass.value = 'bg-green-100 text-green-800'
    
    deleteTarget.value = null
    
    // Reload the databases list after deletion
    await loadMyDatabases()
    
  } catch (err) {
    console.error('Delete error:', err)
    console.error('Error response:', err.response?.data)
    msg.value = err.response?.data?.error || 'Failed to delete database'
    msgClass.value = 'bg-red-100 text-red-800'
  } finally {
    deleting.value = false
  }
}

onMounted(() => {
  loadMyDatabases()
})
</script>

<style scoped>
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.3s ease;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}

.blur-sm {
  filter: blur(4px);
  transition: filter 0.2s ease;
}

.animate-spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
</style>Icon } from '@fortawesome/vue-fontawesome'

const { user } = useAuth()

const dbName = ref('')
const replicas = ref(1)

const msg = ref('')
const msgClass = ref('')
const credentials = ref(null)
const showCredPassword = ref(false)

// Automatically extract username from profile
const username = computed(() => {
  const raw = user.value?.profile?.preferred_username || user.value?.profile?.email || 'user'
  const base = raw.split('@')[0] // Remove domain if email
  const sanitized = base.toLowerCase().replace(/[^a-z0-9-]/g, '') // Remove invalid chars
  return `${sanitized}`
})


function copyToClipboard(text) {
  navigator.clipboard.writeText(text).then(() => {
    msg.value = 'Connection string copied to clipboard!'
    msgClass.value = 'bg-blue-100 text-blue-800'
    setTimeout(() => { msg.value = '' }, 3000)
  }).catch(() => {
    msg.value = 'Failed to copy to clipboard'
    msgClass.value = 'bg-red-100 text-red-800'
  })
}

onMounted(() => {
  loadMyDatabases()
})

<style scoped>
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.3s ease;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}

.blur-sm {
  filter: blur(4px);
  transition: filter 0.2s ease;
}

.animate-spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
</style>
