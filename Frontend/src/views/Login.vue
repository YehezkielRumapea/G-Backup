<template>
  <div class="login-container">
    <form @submit.prevent="handleLogin">
      <h2>G-Backup</h2>
      <p>Please sign in</p>
      
      <div>
        <label for="username">Username</label>
        <input 
          type="text" 
          id="username" 
          v-model="username" 
          required 
        />
      </div>
      
      <div>
        <label for="password">Password</label>
        <input 
          type="password" 
          id="password" 
          v-model="password" 
          required 
        />
      </div>
      
      <button type="submit" :disabled="isLoading">
        {{ isLoading ? 'Loading...' : 'Sign In' }}
      </button>
      
      <p v-if="errorMessage" class="error">{{ errorMessage }}</p>
    </form>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useAuthStore } from '@/stores/authStore' // Import Pinia Store
import { useRouter } from 'vue-router'
import authService from '@/services/authService' // Import service API

// State lokal untuk form
const username = ref('admin') // Default (opsional)
const password = ref('admin123') // Default (opsional)
const errorMessage = ref(null)
const isLoading = ref(false)

const authStore = useAuthStore()
const router = useRouter()

// Fungsi yang dipanggil saat form di-submit
async function handleLogin() {
  errorMessage.value = null
  isLoading.value = true
  
  try {
    // 1. Panggil authService (yang memanggil Axios)
    const response = await authService.login(username.value, password.value)
    
    // 2. Jika sukses, ambil token dari respons
    const token = response.data.token
    if (token) {
      // 3. Simpan token ke Pinia Store
      authStore.setToken(token)
      
      // 4. Arahkan user ke halaman dashboard
      // (Router akan otomatis mengarahkan ke /remotes karena itu default child)
      router.push('/') 
    }
  } catch (error) {
    // Tampilkan error jika backend Golang menolak login (401)
    errorMessage.value = 'Login failed. Please check username or password.'
  } finally {
    isLoading.value = false
  }
}
</script>

<style scoped>
/* Tambahkan sedikit styling untuk form login */
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
}
form {
  padding: 2rem;
  border: 1px solid #ccc;
  border-radius: 8px;
  background: #f9f9f9;
}
div {
  margin-bottom: 1rem;
}
.error {
  color: red;
}
button:disabled {
  background-color: #ccc;
}
</style>    