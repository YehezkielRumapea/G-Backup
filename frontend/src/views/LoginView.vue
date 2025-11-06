<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router'; // 1. Impor useRouter untuk redirect
import apiClient from '@/api'; // 2. Impor 'jembatan' api.js kita
import '@/assets/login.css'; // 3. Impor file CSS baru Anda

// Mendefinisikan event yang akan dikirim ke parent (App.vue)
const emit = defineEmits(['login-success']);
const router = useRouter(); // 4. Inisialisasi router

const username = ref('admin'); // (Default dari backend Anda)
const password = ref('admin123'); // (Default dari backend Anda di main.go)
const errorMsg = ref('');

// 5. Ini adalah fungsi 'async' yang terhubung ke backend
const submitLogin = async () => {
    errorMsg.value = ''; // Reset pesan error
    
    try {
        // 6. Panggil API login backend
        const response = await apiClient.post('/auth/login', {
            username: username.value,
            password: password.value
        });
        
        // 7. SUKSES! Simpan token JWT asli dari backend
        localStorage.setItem('jwt_token', response.data.token);
        localStorage.setItem('isAuthenticated', 'true'); // Tetap gunakan ini untuk router guard

        // 8. Kirim event dan arahkan ke dashboard
        emit('login-success');
        router.push('/'); // Arahkan ke halaman utama ('/')

    } catch (error) {
        // 9. Tangani error jika backend merespons 'Unauthorized'
        console.error("Login gagal:", error.response?.data || error.message);
        errorMsg.value = 'Username atau password salah!';
        password.value = ''; // Kosongkan password saja
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