// src/auth/useAuth.js
import { ref } from 'vue'
import { UserManager, WebStorageStateStore } from 'oidc-client-ts'

// 🔧 Zitadel OIDC configuration
const oidcConfig = {
  authority: 'https://openstack-integration-3vzdfy.us1.zitadel.cloud',
  client_id: '328409700146201818',
  redirect_uri: window.location.origin + '/callback',
  post_logout_redirect_uri: window.location.origin,
  response_type: 'code',
  scope: 'openid profile email offline_access', // Added offline_access for refresh tokens
  userStore: new WebStorageStateStore({ store: window.localStorage }),

  // ✅ Silent renew config - more aggressive settings
  automaticSilentRenew: true,
  silent_redirect_uri: window.location.origin + '/silent-renew.html',
  accessTokenExpiringNotificationTime: 300, // Start renewing 5 minutes before expiry
  silentRequestTimeout: 10000, // 10 seconds timeout for silent requests
  
  // ✅ Additional settings for better token management
  monitorSession: false, // Disable session monitoring to avoid conflicts
  checkSessionInterval: 2000, // Check every 2 seconds
}

export const userManager = new UserManager(oidcConfig)
const user = ref(null)


// ✅ Error logging for silent renew
userManager.events.addSilentRenewError(e => {
  console.error('Silent renew error:', e)
})

// ✅ Error logging for token expired
userManager.events.addUserSignedOut(() => {
  console.warn('User session ended. Signing out.')
  logout()
})


// 🔁 Fetch user or refresh token if expired
async function fetchUser() {
  try {
    let storedUser = await userManager.getUser()
    
    if (!storedUser) {
      // No user found, redirect to login
      userManager.signinRedirect()
      return null
    }
    
    if (storedUser.expired) {
      console.log('Token expired, attempting silent renew...')
      try {
        // Try to refresh the token silently
        storedUser = await userManager.signinSilent()
        console.log('Silent renew successful')
      } catch (error) {
        console.error('Silent renew failed:', error)
        // If silent renew fails, clear session and redirect to login
        await userManager.removeUser()
        await userManager.clearStaleState()
        userManager.signinRedirect()
        return null
      }
    }
    
    user.value = storedUser
    return storedUser
  } catch (err) {
    console.warn('Fetch user failed:', err)
    await userManager.removeUser()
    await userManager.clearStaleState()
    userManager.signinRedirect()
    return null
  }
}


// 🔐 Get fresh token with automatic refresh
export async function getAccessToken() {
  try {
    let currentUser = await userManager.getUser()

    // If no user, don't try to refresh, just return null
    if (!currentUser) {
      console.log('No user found')
      return null
    }

    // If token is expired, try to get a new one
    if (currentUser.expired) {
      console.log('Access token expired, attempting refresh...')
      
      try {
        // Try silent signin (uses refresh token if available)
        currentUser = await userManager.signinSilent()
        user.value = currentUser
        console.log('Token refreshed successfully via silent signin')
        return currentUser.access_token
      } catch (silentError) {
        console.error('Silent signin failed:', silentError)
        // Don't auto-redirect to avoid infinite loops
        return null
      }
    }

    return currentUser.access_token
  } catch (error) {
    console.error('Token access failed:', error)
    return null
  }
}


// 🔓 Login
function login() {
  return userManager.signinRedirect()
}


// 🔒 Logout (with full session clear)
async function logout() {
  try {
    await userManager.removeUser()
    await userManager.clearStaleState()
    await userManager.signoutRedirect()
  } catch (e) {
    console.error('Logout failed:', e)
  }
}


// ✅ Callback handler after login redirect
async function handleCallback() {
  try {
    await userManager.signinRedirectCallback()
    await fetchUser()
  } catch (e) {
    console.error('Callback error:', e)
  }
}


// 🧹 Clear all authentication data and restart
export async function clearAuthAndRestart() {
  try {
    console.log('Clearing all authentication data...')
    
    // Clear all stored tokens and state
    await userManager.removeUser()
    await userManager.clearStaleState()
    
    // Clear localStorage
    localStorage.removeItem('oidc.user:' + oidcConfig.authority + ':' + oidcConfig.client_id)
    
    // Reset user state
    user.value = null
    
    // Force a fresh login
    console.log('Starting fresh authentication...')
    await userManager.signinRedirect()
    
  } catch (error) {
    console.error('Error clearing auth:', error)
    // Force page reload as last resort
    window.location.reload()
  }
}
export async function refreshTokenOrLogin() {
  try {
    console.log('Attempting to refresh authentication...')
    
    // First, clear any existing expired/invalid tokens
    await userManager.removeUser()
    await userManager.clearStaleState()
    
    // Try to get a fresh token by redirecting to auth
    console.log('Redirecting to authentication...')
    await userManager.signinRedirect({
      prompt: 'none' // Try silent authentication first
    })
    
  } catch (error) {
    console.error('Token refresh failed, forcing login:', error)
    // If silent auth fails, force full login
    await userManager.signinRedirect()
  }
}

// 🔄 Check and refresh token if needed (call this before API requests)
export async function ensureValidToken() {
  try {
    let currentUser = await userManager.getUser()
    
    if (!currentUser || currentUser.expired) {
      console.log('No valid token found')
      return null // Don't auto-redirect to avoid loops
    }
    
    return currentUser.access_token
  } catch (error) {
    console.error('Token validation failed:', error)
    return null
  }
}


// 🧠 Auth composable
export function useAuth() {
  return {
    user,
    login,
    logout,
    fetchUser,
    handleCallback,
    getAccessToken,
    ensureValidToken,
    refreshTokenOrLogin,
    clearAuthAndRestart,
  }
}
