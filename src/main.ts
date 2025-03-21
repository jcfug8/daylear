import './assets/main.css';

import { createApp } from 'vue';
import { createPinia } from 'pinia';

import { createVuetify } from 'vuetify';
import { VApp, VAppBar, VBtn, VList, VListGroup, VListItem, VMain, VNavigationDrawer } from 'vuetify/components';
import 'vuetify/styles'; // Import Vuetify styles
import '@mdi/font/css/materialdesignicons.css'; // Import Material Design Icons

import App from './App.vue';
import router from './router';

const vuetify = createVuetify({
    components: {
        VBtn,
        VApp,
        VAppBar,
        VMain,
        VNavigationDrawer,
        VList,
        VListItem,
        VListGroup
    },
});

const app = createApp(App);

app.use(createPinia());
app.use(vuetify);
app.use(router);

app.mount('#app');