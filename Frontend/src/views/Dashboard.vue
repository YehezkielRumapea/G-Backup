<template>
  <div class="dashboard-layout">
    <NavBar />
    
    <main class="content-area">
      <router-view />
      
      </main>
  </div>
</template>

<script setup>
// Impor NavBar untuk digunakan di template
import NavBar from '@/components/NavBar.vue'
</script>

<style scoped>
.dashboard-layout {
  display: flex;
  height: 100vh;
  width: 100vw;
  background-color: #f4f7f6; /* Latar belakang area konten */
}
.content-area {
  flex-grow: 1;
  padding: 2rem;
  overflow-y: auto; /* Memungkinkan scrolling jika konten panjang */
}
</style>