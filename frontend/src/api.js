import axios from 'axios';
import router from './router'; // Pastikan file router Anda ada di 'src/router/index.js'

// --- KONFIGURASI UTAMA ---
// URL relatif ke backend. Nginx akan meneruskan permintaan yang dimulai dengan /api/v1
// ke server Go Anda yang berjalan (misal: di port 8080).
const API_BASE_URL = '/api/v1';

// Buat instance Axios dengan konfigurasi dasar
const apiClient = axios.create({
    baseURL: API_BASE_URL,
    headers: {
        'Content-Type': 'application/json',
    },
});

/**
 * Interceptor Permintaan (Request Interceptor):
 * Berjalan SEBELUM setiap permintaan API dikirim.
 * Fungsinya: Mengambil token JWT dari local storage dan menyisipkannya ke header 'Authorization'.
 */
apiClient.interceptors.request.use(config => {
    // Ambil token dari local storage
    const token = localStorage.getItem('jwt_token');
    
    // Jika token ada, tambahkan ke header
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    
    // Lanjutkan permintaan
    return config;
}, error => {
    // Tangani error saat persiapan request (jarang terjadi)
    console.error("Axios request interceptor error:", error);
    return Promise.reject(error);
});

/**
 * Interceptor Respons (Response Interceptor):
 * Berjalan SETELAH respons dari backend diterima.
 * Fungsinya: Memeriksa jika ada error 401 (Unauthorized), yang berarti token
 * tidak valid atau kedaluwarsa. Jika ya, hapus token lama dan paksa kembali ke login.
 */
apiClient.interceptors.response.use(
    // Jika respons sukses (status 2xx), langsung teruskan
    response => response,
    
    // Jika respons error
    error => {
        if (error.response && error.response.status === 401) {
            console.warn("Sesi berakhir atau token tidak valid. Anda akan diarahkan ke halaman login.");
            
            // Hapus token lama yang tidak valid
            localStorage.removeItem('jwt_token');
            
            // Arahkan ke halaman login menggunakan Vue Router
            // Gunakan 'replace' agar pengguna tidak bisa kembali ke halaman sebelumnya dengan tombol back
            router.replace({ name: 'Login' });
        }
        
        // Lanjutkan error agar bisa ditangani lebih lanjut di komponen Vue (misal: menampilkan pesan error)
        return Promise.reject(error);
    }
);

// Ekspor instance apiClient yang sudah dikonfigurasi agar bisa diimpor dan dipakai di file lain
export default apiClient;
