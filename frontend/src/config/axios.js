// src/config/axios.js
import axios from 'axios'
import { getAccessToken } from '../auth/useAuth.js'

// Create axios instance
const apiClient = axios.create()

// Request interceptor to add auth token
apiClient.interceptors.request.use(
  async (config) => {
    try {
      const token = await getAccessToken()
      if (token) {
        config.headers.Authorization = `Bearer ${token}`
      }
    } catch (error) {
      console.error('Failed to get access token:', error)
      // Let the request continue without token - auth middleware will handle it
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor to handle token refresh on 401
apiClient.interceptors.response.use(
  (response) => {
    return response
  },
  async (error) => {
    const originalRequest = error.config

    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true

      try {
        console.log('Received 401, attempting to refresh token...')
        const newToken = await getAccessToken()
        originalRequest.headers.Authorization = `Bearer ${newToken}`
        return apiClient(originalRequest)
      } catch (refreshError) {
        console.error('Token refresh failed:', refreshError)
        // Token refresh failed, user will be redirected to login
        return Promise.reject(error)
      }
    }

    return Promise.reject(error)
  }
)

export default apiClient
