<template>
  <div class="setup-wizard">
    <h2>Selamat Datang! Lakukan Setup Awal</h2>
    <p>Buat akun administrator utama Anda.</p>
    
    <form @submit.prevent="handleSubmit">
      <div class="form-group">
        <label for="username">Username</label>
        <input type="text" id="username" v-model="username" required :disabled="loading">
      </div>
      <div class="form-group">
        <label for="password">Password</label>
        <input type="password" id="password" v-model="password" required :disabled="loading">
      </div>
      
      <p v-if="error" class="error-message">{{ error }}</p>
      
      <button type="submit" :disabled="loading">
        {{ loading ? 'Mendaftarkan...' : 'Daftar & Lanjutkan' }}
      </button>
    </form>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useAppStore } from '@/stores/app';
import { useRouter } from 'vue-router';

const appStore = useAppStore();
const router = useRouter();

const username = ref('');
const password = ref('');
const error = ref(null);
const loading = ref(false);

const handleSubmit = async () => {
  error.value = null;
  loading.value = true;
  
  try {
    const response = await fetch('/api/v1/setup/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username: username.value, password: password.value }),
    });

    if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Gagal mendaftar admin pertama.');
    }

    // Sukses: Ubah status di store dan redirect
    appStore.setSetupComplete(); // <-- Mengubah status di store
    alert('Admin berhasil terdaftar! Silakan Login.');
    router.push({ name: 'login' }); // Redirect ke halaman login standar
    
  } catch (err) {
    error.value = err.message;
  } finally {
    loading.value = false;
  }
};
</script>