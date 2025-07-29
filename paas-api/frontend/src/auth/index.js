import { UserManager, WebStorageStateStore } from 'oidc-client-ts'

const oidcConfig = {
  authority: 'https://openstack-integration-3vzdfy.us1.zitadel.cloud',
  client_id: '328409700146201818',
  redirect_uri: window.location.origin + '/callback',
  post_logout_redirect_uri: window.location.origin,
  response_type: 'code',
  scope: 'openid profile email',
  userStore: new WebStorageStateStore({ store: window.localStorage }),
}

const userManager = new UserManager(oidcConfig)

export const login = () => userManager.signinRedirect()
export const logout = () => userManager.signoutRedirect()
export const getUser = () => userManager.getUser()
export const handleCallback = () => userManager.signinRedirectCallback()
