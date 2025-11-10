import { defineStore } from "pinia";
import { ref } from "vue";

export const useAuthStore = defineStore('auth', () => {
    const token = ref(localStorage.getItem('token') || null)

    function isAuthenticated() {
            return token.value !== null
    }

    function setToken(newToken) {
        token.value = newToken
        localStorage.setItem('token', newToken) 
    }

    function clearToken() {
        token.value = null
        localStorage.removeItem('token')
    }

    return {token, isAuthenticated, setToken, clearToken}
})