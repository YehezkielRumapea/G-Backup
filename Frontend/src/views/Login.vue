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
          @input="clearError"
        />
      </div>
      
      <div class="form-group">
        <label for="password">Password</label>
        <input 
          type="password" 
          id="password" 
          v-model="password" 
          placeholder="Enter your password"
          @input="clearError"
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
const username = ref('') 
const password = ref('') 
const errorMessage = ref(null)
const isLoading = ref(false)

const authStore = useAuthStore()
const router = useRouter()

function clearError() {
  errorMessage.value = null
}

async function handleLogin() {
  errorMessage.value = null
  
  // Validasi field kosong
  if (!username.value.trim() && !password.value.trim()) {
    errorMessage.value = 'Plese enter your username and password!'
    return
  }
  
  if (!username.value.trim()) {
    errorMessage.value = 'Plese enter your username!'
    return
  }
  
  if (!password.value.trim()) {
    errorMessage.value = 'Plese enter your password!'
    return
  }
  
  isLoading.value = true
  
  try {
    const response = await authService.login(username.value, password.value)
    
    const token = response.data.token
    if (token) {
      authStore.setToken(token)
      router.push('/') 
    }
  } catch (error) {
    // Cek response dari server untuk menentukan pesan error yang tepat
    if (error.response) {
      const status = error.response.status
      const data = error.response.data
      
      // Jika backend mengirim pesan spesifik
      if (data.message) {
        if (data.message.includes('username') && data.message.includes('incorrect')) {
          errorMessage.value = 'Wrong username! Please check your username.'
        } else if (data.message.includes('password') && data.message.includes('incorrect')) {
          errorMessage.value = 'Wrong password! Please check your password.'
        } else if (data.message.includes('not found') || data.message.includes('does not exist')) {
          errorMessage.value = 'Wrong username or password! Account not found.'
        } else {
          errorMessage.value = data.message
        }
      } 
      // Jika backend mengirim field spesifik yang salah
      else if (data.field) {
        if (data.field === 'username') {
          errorMessage.value = 'Wrong username! Please check your username.'
        } else if (data.field === 'password') {
          errorMessage.value = 'Wrong password! Please check your password.'
        }
      }
      // Status code based error handling
      else if (status === 401) {
        errorMessage.value = 'Wrong username or password! Please try again.'
      } else if (status === 404) {
        errorMessage.value = 'Account not found in the database!'
      } else {
        errorMessage.value = 'Login failed. Please check your username and password.'
      }
    } else if (error.request) {
      errorMessage.value = 'Could not connect to the server. Please check your internet connection.'
    } else {
      errorMessage.value = 'An error occurred. Please try again.'
    }
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
  min-height: 100vh;
  width: 100%;
  background-color: #f3f4f6;
  padding: 20px;
  box-sizing: border-box;
}

/* 2. Kartu Login (Kotak Putih) */
.login-card {
  width: 100%;
  max-width: 400px;
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
  box-sizing: border-box;
}

input:focus {
  outline: none;
  border-color: #10b981;
  background-color: #ffffff;
  box-shadow: 0 0 0 3px rgba(16, 185, 129, 0.1);
}

/* 5. Tombol Submit */
.submit-btn {
  width: 100%;
  padding: 0.875rem;
  background-color: #10b981;
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
  animation: slideDown 0.3s ease-out;
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>