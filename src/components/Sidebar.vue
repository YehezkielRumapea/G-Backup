<script setup>
import { RouterLink } from 'vue-router';

// Mendefinisikan event 'logout' yang akan dikirim ke parent component (App.vue)
const emit = defineEmits(['logout']);

const handleLogoutClick = () => {
    emit('logout');
};
</script>

<template>
    <nav class="sidebar">
        <h2 class="sidebar-title">G-Backup</h2>
        <RouterLink to="/" class="sidebar-link">Dashboard</RouterLink>
        <RouterLink to="/remote" class="sidebar-link">Remote</RouterLink>
        <RouterLink to="/job-status" class="sidebar-link">Job Status</RouterLink>
	<RouterLink to="/scheduler" class="sidebar-link">Scheduler</RouterLink>
        
        <!-- Tombol Logout akan memanggil event 'logout' -->
        <a href="#" @click.prevent="handleLogoutClick" class="sidebar-link sidebar-logout">Logout</a>
    </nav>
</template>

<style scoped>
/* --- SIDEBAR STYLING --- */
.sidebar {
  width: 250px;
  background-color: var(--sidebar-bg);
  color: var(--sidebar-text);
  padding: 20px 0;
  display: flex;
  flex-direction: column;
  box-shadow: 2px 0 5px rgba(0, 0, 0, 0.1);
  flex-shrink: 0;
  height: 100vh; /* Memastikan sidebar memenuhi tinggi layar */
}

.sidebar-title {
  color: var(--sidebar-active);
  text-align: center;
  margin-bottom: 30px;
  font-size: 1.5rem;
  padding: 0 20px;
}

.sidebar-link {
  color: var(--sidebar-text);
  padding: 12px 20px;
  text-decoration: none;
  font-weight: 500;
  transition: background-color 0.3s, color 0.3s;
  display: block;
}

.sidebar-link:hover,
.router-link-exact-active { /* Menggunakan class yang lebih spesifik untuk link aktif */
  background-color: var(--primary-color);
  color: var(--sidebar-active);
}

.sidebar-logout {
    margin-top: auto; /* Mendorong tombol logout ke bagian bawah */
    border-top: 1px solid rgba(255, 255, 255, 0.1);
}

/* Penyesuaian Responsif */
@media (max-width: 768px) {
    .sidebar {
        width: 100%;
        height: auto;
        flex-direction: row;
        justify-content: space-around;
        padding: 10px 0;
    }
    .sidebar-title {
        display: none;
    }
    .sidebar-link {
        padding: 8px 10px;
        white-space: nowrap;
    }
    .sidebar-logout {
        margin-top: 0;
        border-top: none;
    }
}
</style>

