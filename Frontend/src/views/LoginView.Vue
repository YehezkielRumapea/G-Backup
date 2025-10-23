<script setup>
import { ref } from 'vue';

// Mendefinisikan event yang akan dikirim ke parent (App.vue)
const emit = defineEmits(['login-success']);

const username = ref('');
const password = ref('');
const errorMsg = ref('');

const submitLogin = () => {
    errorMsg.value = ''; // Reset pesan error
    // SIMULASI proses login
    if (username.value === 'admin' && password.value === 'admin') {
        localStorage.setItem('isAuthenticated', 'true');
        // Kirim event ke App.vue untuk memberitahu login berhasil
        emit('login-success');
    } else {
        errorMsg.value = 'Username atau password salah!';
        username.value = '';
        password.value = '';
    }
};
</script>

<template>
    <div class="login-page">
        <div class="login-card">
            <h1 class="card-title">G-Backup</h1>
            <p>Please sign in to continue</p>
            <form @submit.prevent="submitLogin">
                <input type="text" v-model="username" placeholder="Username (admin)" required class="input-field" />
                <input type="password" v-model="password" placeholder="Password (admin)" required class="input-field" />
                <p v-if="errorMsg" class="error-message">{{ errorMsg }}</p>
                <button type="submit" class="btn btn-primary">Login</button>
            </form>
        </div>
    </div>
</template>

<style scoped>
/* Style spesifik untuk halaman login */
.login-page {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: var(--body-bg);
}

.login-card {
  background: var(--card-bg);
  padding: 40px;
  border-radius: 12px;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 400px;
  text-align: center;
}

.card-title {
    margin-bottom: 10px;
    color: var(--primary-color);
    font-weight: 700;
}

.login-card p {
    margin-bottom: 25px;
    color: #6c757d;
}

.input-field {
  width: 100%;
  padding: 12px;
  margin-bottom: 15px;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  transition: border-color 0.3s;
}

.input-field:focus {
  border-color: var(--primary-color);
  outline: none;
}

.btn {
  padding: 12px 20px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 1rem;
  transition: background-color 0.3s;
  width: 100%;
  font-weight: 500;
  margin-top: 10px;
}

.btn-primary {
  background-color: var(--primary-color);
  color: white;
}

.btn-primary:hover {
  background-color: var(--primary-hover);
}

.error-message {
    color: #dc3545; /* Red */
    margin-bottom: 15px;
    font-size: 0.9rem;
}
</style>
