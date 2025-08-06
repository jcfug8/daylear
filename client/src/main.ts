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
  VListItemTitle,
  VListItemSubtitle,
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
  VFileInput,
  VRadioGroup,
  VRadio,
  VDialog,
  VWindow,
  VWindowItem,
  VAutocomplete,
  VSpeedDial,
  VChip,
  VChipGroup,
  VAlert,
  VProgressLinear,
  VProgressCircular,      
  VCombobox,
  VAvatar,
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
    VListItemTitle,
    VListItemSubtitle,
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
    VFileInput,
    VRadioGroup,
    VRadio,
    VDialog,
    VWindow,
    VWindowItem,
    VAutocomplete,
    VSpeedDial,
    VChip,
    VChipGroup,
    VAlert,
    VProgressLinear,
    VProgressCircular,
    VCombobox,
    VAvatar,
  },
})

const app = createApp(App)

app.use(createPinia())
app.use(vuetify)
app.use(router)

app.mount('#app')
