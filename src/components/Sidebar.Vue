<script setup>
import { RouterLink } from 'vue-router';

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
        <RouterLink to="/configuration" class="sidebar-link">Configuration</RouterLink>
        <RouterLink to="/logs" class="sidebar-link">Log Activity</RouterLink>
        <!-- Tombol Logout akan memanggil event 'logout' ke App.vue -->
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
  flex-shrink: 0; /* Pastikan tidak menyusut */
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
  padding: 10px 20px;
  text-decoration: none;
  font-weight: 400;
  transition: background-color 0.3s, color 0.3s;
}

.sidebar-link:hover,
.sidebar-link.router-link-active {
  background-color: var(--sidebar-link-hover);
  color: var(--sidebar-active);
}

.sidebar-logout {
    margin-top: auto; /* Dorong ke bawah */
    border-top: 1px solid rgba(255, 255, 255, 0.1);
}

/* Responsif */
@media (max-width: 600px) {
    .sidebar {
        width: 100%;
        height: auto;
        flex-direction: row;
        justify-content: space-around;
        padding: 10px 0;
        overflow-x: auto;
    }
    .sidebar-title {
        display: none; /* Sembunyikan judul di mobile */
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
