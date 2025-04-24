import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import { createVuetify } from 'vuetify'
import {
  VApp,
  VAppBar,
  VAppBarNavIcon,
  VAppBarTitle,
  VBreadcrumbs,
  VBtn,
  VCard,
  VCardActions,
  VCardSubtitle,
  VCardText,
  VCardTitle,
  VCheckbox,
  VCol,
  VContainer,
  VDivider,
  VFab,
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
  VSelect,
  VSheet,
  VSpacer,
  VTab,
  VTabs,
  VTabsWindow,
  VTabsWindowItem,
  VTextarea,
  VTextField,
} from 'vuetify/components'
import 'vuetify/styles' // Import Vuetify styles
import '@mdi/font/css/materialdesignicons.css' // Import Material Design Icons

import App from './App.vue'
import router from './router'
import { VCalendar, VNumberInput } from 'vuetify/labs/components'

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
    VCardTitle,
    VCardActions,
    VCardText,
    VSpacer,
    VCardSubtitle,
    VImg,
    VBreadcrumbs,
    VTab,
    VTabs,
    VTabsWindow,
    VTabsWindowItem,
    VTextField,
    VTextarea,
    VSelect,
    VNumberInput,
    VCheckbox,
    VFab,
  },
})

const app = createApp(App)

app.use(createPinia())
app.use(vuetify)
app.use(router)

app.mount('#app')
