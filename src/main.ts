import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import { createVuetify } from 'vuetify'
import {
  VApp,
  VAppBar,
  VAppBarTitle,
  VBtn,
  VContainer,
  VDivider,
  VForm,
  VIcon,
  VLayout,
  VList,
  VListGroup,
  VListItem,
  VMain,
  VMenu,
  VNavigationDrawer,
  VSheet,
} from 'vuetify/components'
import 'vuetify/styles' // Import Vuetify styles
import '@mdi/font/css/materialdesignicons.css' // Import Material Design Icons

import App from './App.vue'
import router from './router'

const vuetify = createVuetify({
  components: {
    VBtn,
    VAppBar,
    VMain,
    VNavigationDrawer,
    VList,
    VListItem,
    VListGroup,
    VAppBarTitle,
    VMenu,
    VDivider,
    VIcon,
    VContainer,
    VLayout,
    VApp,
    VSheet,
    VForm,
  },
})

const app = createApp(App)

app.use(createPinia())
app.use(vuetify)
app.use(router)

app.mount('#app')
