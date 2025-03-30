import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import { createVuetify } from 'vuetify'
import {
  VApp,
  VAppBar,
  VAppBarNavIcon,
  VAppBarTitle,
  VBtn,
  VCard,
  VCardActions,
  VCardSubtitle,
  VCardText,
  VCol,
  VContainer,
  VDivider,
  VForm,
  VIcon,
  VImg,
  VLayout,
  VList,
  VListGroup,
  VListItem,
  VMain,
  VMenu,
  VNavigationDrawer,
  VRow,
  VSheet,
  VSpacer,
} from 'vuetify/components'
import 'vuetify/styles' // Import Vuetify styles
import '@mdi/font/css/materialdesignicons.css' // Import Material Design Icons

import App from './App.vue'
import router from './router'
import { VCalendar } from 'vuetify/labs/components'

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
    VCalendar,
    VAppBarNavIcon,
    VRow,
    VCol,
    VCard,
    VCardActions,
    VCardText,
    VSpacer,
    VCardSubtitle,
    VImg,
  },
})

const app = createApp(App)

app.use(createPinia())
app.use(vuetify)
app.use(router)

app.mount('#app')
