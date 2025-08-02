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
  scope: 'openid profile email',
  userStore: new WebStorageStateStore({ store: window.localStorage }),

  // ✅ Silent renew config
  automaticSilentRenew: true,
  silent_redirect_uri: window.location.origin + '/silent-renew.html',
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


// 🔁 Fetch user or login if not found
async function fetchUser() {
  try {
    const storedUser = await userManager.getUser()
    if (!storedUser || storedUser.expired) {
      await userManager.removeUser()
      await userManager.clearStaleState()
      userManager.signinRedirect()
      return null
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


// 🔐 Get fresh token
export async function getAccessToken() {
  let currentUser = await userManager.getUser()

  if (!currentUser || currentUser.expired) {
    console.warn('Access token expired or missing. Clearing session.')
    await userManager.removeUser()
    await userManager.clearStaleState()
    userManager.signinRedirect()
    throw new Error('Token expired, redirecting to login.')
  }

  return currentUser.access_token
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


// 🧠 Auth composable
export function useAuth() {
  return {
    user,
    login,
    logout,
    fetchUser,
    handleCallback,
    getAccessToken,
  }
}
