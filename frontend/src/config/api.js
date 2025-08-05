// API Configuration
export const API_BASE_URL = 'http://91.99.215.178:8080'

// API Endpoints
export const API_ENDPOINTS = {
  // Database endpoints
  DATABASES: '/databases',
  DATABASE_STATUS: '/databases/:username/:dbName/status',
  
  // Pod endpoints  
  PODS: '/pods/:namespace',
  
  // Admin endpoints
  ADMIN_TENANT_PODS: '/admin/tenants/pods'
}
