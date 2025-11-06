import { createApp } from 'vue';
import App from './App.vue';
import router from './router'; // PERBAIKI: Impor router dari file terpisah

import './assets/main.css'; // TAMBAHKAN: Impor file CSS global

const app = createApp(App);

app.use(router);

app.mount('#app');
