<template>
  <div class="app-layout">
    <!-- ==================== HEADER ==================== -->
    <header class="app-header">
      <div class="header-content">
        <div class="logo">
          <h1>☁️ G-Backup</h1>
        </div>
        
        <nav class="main-nav">
          <router-link to="/dashboard" class="nav-item">
            <span class="nav-icon">📊</span>
            <span class="nav-text">Dashboard</span>
          </router-link>
          
          <router-link to="/remotes" class="nav-item">
            <span class="nav-icon">☁️</span>
            <span class="nav-text">Remotes</span>
          </router-link>
          
          <router-link to="/scheduled" class="nav-item">
            <span class="nav-icon">⏰</span>
            <span class="nav-text">Scheduled Jobs</span>
          </router-link>
          
          <router-link to="/manual" class="nav-item">
            <span class="nav-icon">🔧</span>
            <span class="nav-text">Manual Jobs</span>
          </router-link>
          
          <router-link to="/create" class="nav-item">
            <span class="nav-icon">➕</span>
            <span class="nav-text">Create Job</span>
          </router-link>
          
          <router-link to="/logs" class="nav-item">
            <span class="nav-icon">📝</span>
            <span class="nav-text">Logs</span>
          </router-link>
        </nav>
        
        <div class="header-actions">
          <button @click="handleLogout" class="logout-btn">
            <span>🚪</span>
            <span>Logout</span>
          </button>
        </div>
      </div>
    </header>
    
    <!-- ==================== MAIN CONTENT ==================== -->
    <main class="app-main">
      <!-- ⭐ INI TEMPAT CHILD ROUTES DI-RENDER -->
      <transition name="fade" mode="out-in">
        <router-view />
      </transition>
    </main>
    
    <!-- ==================== FOOTER (Optional) ==================== -->
    <footer class="app-footer">
      <p>© 2024 G-Backup System | Powered by Vue.js & Golang</p>
    </footer>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/authStore';

const router = useRouter();
const authStore = useAuthStore();

const handleLogout = () => {
  if (confirm('Are you sure you want to logout?')) {
    authStore.logout();
    router.push('/login');
  }
};
</script>

<style scoped>
/* ==================== LAYOUT ==================== */
.app-layout {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  background-color: #f5f7fa;
}

/* ==================== HEADER ==================== */
.app-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-content {
  max-width: 1400px;
  margin: 0 auto;
  padding: 1rem 2rem;
  display: flex;
  align-items: center;
  gap: 2rem;
}

.logo h1 {
  margin: 0;
  font-size: 1.5rem;
  font-weight: 700;
  white-space: nowrap;
}

/* ==================== NAVIGATION ==================== */
.main-nav {
  display: flex;
  gap: 0.5rem;
  flex: 1;
  overflow-x: auto;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1rem;
  color: white;
  text-decoration: none;
  border-radius: 8px;
  transition: all 0.2s ease;
  white-space: nowrap;
  font-weight: 500;
}

.nav-item:hover {
  background: rgba(255, 255, 255, 0.15);
  transform: translateY(-2px);
}

.nav-item.router-link-active {
  background: rgba(255, 255, 255, 0.25);
  font-weight: 600;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.nav-icon {
  font-size: 1.2rem;
}

.nav-text {
  font-size: 0.9rem;
}

/* ==================== HEADER ACTIONS ==================== */
.header-actions {
  display: flex;
  align-items: center;
}

.logout-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.25rem;
  background: rgba(255, 255, 255, 0.2);
  color: white;
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 8px;
  cursor: pointer;
  font-weight: 600;
  transition: all 0.2s ease;
  white-space: nowrap;
}

.logout-btn:hover {
  background: rgba(255, 255, 255, 0.3);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.logout-btn:active {
  transform: translateY(0);
}

/* ==================== MAIN CONTENT ==================== */
.app-main {
  flex: 1;
  padding: 2rem;
  max-width: 1400px;
  width: 100%;
  margin: 0 auto;
}

/* ==================== FOOTER ==================== */
.app-footer {
  background: #2c3e50;
  color: white;
  padding: 1.5rem 2rem;
  text-align: center;
  border-top: 3px solid #667eea;
}

.app-footer p {
  margin: 0;
  font-size: 0.9rem;
  opacity: 0.8;
}

/* ==================== TRANSITIONS ==================== */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* ==================== RESPONSIVE ==================== */
@media (max-width: 1024px) {
  .header-content {
    flex-direction: column;
    gap: 1rem;
  }
  
  .main-nav {
    width: 100%;
    justify-content: center;
  }
  
  .nav-text {
    display: none;
  }
  
  .nav-item {
    padding: 0.75rem;
  }
}

@media (max-width: 768px) {
  .app-main {
    padding: 1rem;
  }
  
  .header-content {
    padding: 1rem;
  }
  
  .logo h1 {
    font-size: 1.2rem;
  }
  
  .main-nav {
    gap: 0.25rem;
  }
  
  .logout-btn span:last-child {
    display: none;
  }
}
</style>