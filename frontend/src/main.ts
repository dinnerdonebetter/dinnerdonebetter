import { createApp } from 'vue';
import App from './App.vue';
import router from './router';
import './assets/main.css';

import AppLayout from './layouts/App.vue';
import AdminLayout from './layouts/Admin.vue';
import EmptyLayout from './layouts/Empty.vue';

const app = createApp(App);

app.component('app-layout', AppLayout);
app.component('admin-layout', AdminLayout);
app.component('empty-layout', EmptyLayout);

app.use(router);
app.mount('#app');
