import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import './assets/index.css'

import { library } from '@fortawesome/fontawesome-svg-core'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

// ✅ MOVE ALL ICON IMPORTS TOGETHER
import {
  faSyncAlt, faSpinner, faServer, faSearch,
  faPlusCircle, faEye, faEyeSlash,
  faHome, faDatabase, faCubes, faUserShield,
  faSignOutAlt, faSignInAlt, faKey
} from '@fortawesome/free-solid-svg-icons'

// ✅ THEN ADD TO LIBRARY
library.add(
  faSyncAlt, faSpinner, faServer, faSearch, faPlusCircle, faEye, faEyeSlash,
  faHome, faDatabase, faCubes, faUserShield, faSignOutAlt, faSignInAlt, faKey
)

const app = createApp(App)
app.component('FontAwesomeIcon', FontAwesomeIcon)
app.use(router).mount('#app')
