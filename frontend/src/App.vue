<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import Sidebar from './components/Sidebar.vue'; // PERBAIKI: Impor komponen Sidebar

const router = useRouter();

// State untuk autentikasi, membaca dari localStorage
const isAuthenticated = ref(localStorage.getItem('isAuthenticated') === 'true');

// Fungsi untuk logout
const logout = () => {
  isAuthenticated.value = false;
  localStorage.removeItem('isAuthenticated');
  router.push('/login'); // Arahkan ke halaman login setelah logout
};

// Fungsi untuk menandai bahwa user sudah login (dipanggil dari LoginView)
const handleLoginSuccess = () => {
  isAuthenticated.value = true;
  router.push('/');
};
</script>

<template>
  <div v-if="isAuthenticated" class="app-layout">
    <!-- Gunakan komponen Sidebar yang sudah diimpor -->
    <Sidebar @logout="logout" />
    <main class="main-content-area">
      <RouterView />
    </main>
  </div>
  <div v-else class="auth-layout">
    <!-- 
      RouterView akan menampilkan LoginView. 
      Kita tambahkan event 'login-success' untuk memberitahu App.vue
    -->
    <RouterView @login-success="handleLoginSuccess" />
  </div>
</template>

<style>
/* Pindahkan semua style global (seperti :root, body, app-layout) ke src/assets/main.css 
  agar tidak mengotori komponen App.vue. 
  Style yang spesifik untuk App.vue bisa diletakkan di sini dengan tag <style scoped>.
*/
.app-layout {
  display: flex;
  min-height: 100vh;
}

.main-content-area {
  flex-grow: 1;
  padding: 20px;
  background-color: var(--body-bg);
}

/* Penyesuaian Responsif untuk layout utama */
@media (max-width: 600px) {
    .app-layout {
        flex-direction: column;
    }
}
</style>
