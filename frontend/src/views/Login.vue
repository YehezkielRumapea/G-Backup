<template>
  <div class="login-container">
    <form @submit.prevent="handleLogin" class="login-card">
      <div class="form-header">
        <h2>G-Backup</h2>
        <p>Please sign in to continue</p>
      </div>
      
      <div class="form-group">
        <label for="username">Username</label>
        <input 
          type="text" 
          id="username" 
          v-model="username" 
          placeholder="Enter your username"
          required 
        />
      </div>
      
      <div class="form-group">
        <label for="password">Password</label>
        <input 
          type="password" 
          id="password" 
          v-model="password" 
          placeholder="Enter your password"
          required 
        />
      </div>
      
      <button type="submit" :disabled="isLoading" class="submit-btn">
        <span v-if="isLoading">Loading...</span>
        <span v-else>Sign In</span>
      </button>
      
      <div v-if="errorMessage" class="error-alert">
        {{ errorMessage }}
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useAuthStore } from '@/stores/authStore' 
import { useRouter } from 'vue-router'
import authService from '@/services/authService'

// State
const username = ref('admin') 
const password = ref('admin123') 
const errorMessage = ref(null)
const isLoading = ref(false)

const authStore = useAuthStore()
const router = useRouter()

async function handleLogin() {
  errorMessage.value = null
  isLoading.value = true
  
  try {
    const response = await authService.login(username.value, password.value)
    
    const token = response.data.token
    if (token) {
      authStore.setToken(token)
      router.push('/') 
    }
  } catch (error) {
    errorMessage.value = 'Login failed. Please check your credentials.'
  } finally {
    isLoading.value = false
  }
}
</script>

<style scoped>
/* 1. Layout Utama (Agar di tengah) */
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh; /* Tinggi minimal setinggi layar browser */
  width: 100%;
  background-color: #f3f4f6; /* Warna background abu-abu muda modern */
  padding: 20px; /* Padding agar tidak nempel tepi di layar HP */
  box-sizing: border-box;
}

/* 2. Kartu Login (Kotak Putih) */
.login-card {
  width: 100%;
  max-width: 400px; /* Lebar maksimal agar rapi */
  background: #ffffff;
  padding: 2.5rem;
  border-radius: 12px;
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
}

/* 3. Header Teks */
.form-header {
  text-align: center;
  margin-bottom: 2rem;
}

.form-header h2 {
  margin: 0 0 0.5rem 0;
  color: #111827;
  font-size: 1.75rem;
  font-weight: 700;
}

.form-header p {
  margin: 0;
  color: #6b7280;
  font-size: 0.95rem;
}

/* 4. Input Fields */
.form-group {
  margin-bottom: 1.5rem;
}

label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 600;
  color: #374151;
  font-size: 0.9rem;
}

input {
  width: 100%;
  padding: 0.75rem 1rem;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 1rem;
  color: #1f2937;
  background-color: #f9fafb;
  transition: all 0.2s;
  box-sizing: border-box; /* Penting agar padding tidak merusak layout */
}

input:focus {
  outline: none;
  border-color: #10b981; /* Warna Hijau saat aktif */
  background-color: #ffffff;
  box-shadow: 0 0 0 3px rgba(16, 185, 129, 0.1);
}

/* 5. Tombol Submit */
.submit-btn {
  width: 100%;
  padding: 0.875rem;
  background-color: #10b981; /* Warna G-Backup */
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.2s;
}

.submit-btn:hover {
  background-color: #059669;
}

.submit-btn:disabled {
  background-color: #9ca3af;
  cursor: not-allowed;
}

/* 6. Pesan Error */
.error-alert {
  margin-top: 1.5rem;
  padding: 0.75rem;
  background-color: #fef2f2;
  border: 1px solid #fecaca;
  color: #991b1b;
  border-radius: 6px;
  text-align: center;
  font-size: 0.9rem;
}
</style>