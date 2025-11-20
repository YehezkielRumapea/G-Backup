<template>
  <div class="app-layout">
    <!-- Header -->
    <header class="app-header">
      <div class="header-content">
        <div class="logo">
          <h1>G-Backup</h1>
        </div>
        
        <nav class="main-nav">
          <router-link to="/dashboard" class="nav-item">
            <span class="nav-text">Dashboard</span>
          </router-link>
          
          <router-link to="/remotes" class="nav-item">
            <span class="nav-text">Remotes</span>
          </router-link>
          
          <router-link to="/scheduled" class="nav-item">
            <span class="nav-text">Scheduled Jobs</span>
          </router-link>
          
          <router-link to="/manual" class="nav-item">
            <span class="nav-text">Jobs</span>
          </router-link>
        </nav>
      </div>
    </header>
    
    <!-- Main Content dengan Log Panel -->
    <div class="content-wrapper">
      <!-- Main Content -->
      <main class="app-main">
        <transition name="fade" mode="out-in">
          <router-view />
        </transition>
      </main>
      
      <!-- Log Panel (Sticky di kanan) -->
      <aside class="log-panel">
        <SimpleLiveLog />
      </aside>
    </div>
    
    <!-- Footer -->
    <footer class="app-footer">
      <p>© 2024 G-Backup System</p>
    </footer>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/authStore';
import SimpleLiveLog from '@/components/LiveLog.vue';

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
.app-layout {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  background-color: #fafafa;
}

/* Header */
.app-header {
  background: #fff;
  border-bottom: 1px solid #e5e5e5;
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-content {
  max-width: 100%;
  margin: 0 auto;
  padding: 1rem 2rem;
  display: flex;
  align-items: center;
  gap: 2rem;
}

.logo h1 {
  margin: 0;
  font-size: 1.25rem;
  font-weight: 600;
  color: #1a1a1a;
  white-space: nowrap;
  letter-spacing: -0.02em;
}

.main-nav {
  display: flex;
  gap: 0.25rem;
  flex: 1;
  overflow-x: auto;
}

.nav-item {
  display: flex;
  align-items: center;
  padding: 0.625rem 1rem;
  color: #666;
  text-decoration: none;
  border-radius: 6px;
  transition: all 0.2s ease;
  white-space: nowrap;
  font-weight: 500;
  font-size: 0.9375rem;
}

.nav-item:hover {
  background: #f5f5f5;
  color: #1a1a1a;
}

.nav-item.router-link-active {
  background: #1a1a1a;
  color: white;
  font-weight: 500;
}

.nav-text {
  font-size: 0.9375rem;
}

.header-actions {
  display: flex;
  align-items: center;
}

.logout-btn {
  display: flex;
  align-items: center;
  padding: 0.625rem 1rem;
  background: transparent;
  color: #666;
  border: 1px solid #e5e5e5;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 500;
  font-size: 0.9375rem;
  transition: all 0.2s ease;
  white-space: nowrap;
}

.logout-btn:hover {
  background: #f5f5f5;
  color: #1a1a1a;
  border-color: #1a1a1a;
}

/* Content Wrapper dengan Log Panel */
.content-wrapper {
  flex: 1;
  display: grid;
  grid-template-columns: 1fr 320px;
  gap: 1.5rem;
  padding: 1.5rem 2rem;
  max-width: 100%;
}

/* Main Content */
.app-main {
  min-width: 0;
}

/* Log Panel (Sticky) */
.log-panel {
  position: sticky;
  top: calc(60px + 1.5rem);
  height: calc(100vh - 60px - 3rem - 60px);
  min-width: 0;
}

/* Footer */
.app-footer {
  background: #fff;
  border-top: 1px solid #e5e5e5;
  padding: 1.25rem 2rem;
  text-align: center;
}

.app-footer p {
  margin: 0;
  font-size: 0.875rem;
  color: #666;
}

/* Transitions */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* Responsive - Hide log panel on smaller screens */
@media (max-width: 1200px) {
  .content-wrapper {
    grid-template-columns: 1fr;
  }
  
  .log-panel {
    display: none;
  }
}

@media (max-width: 1024px) {
  .header-content {
    flex-direction: column;
    gap: 1rem;
    align-items: stretch;
  }
  
  .main-nav {
    width: 100%;
    justify-content: center;
  }
}

@media (max-width: 768px) {
  .content-wrapper {
    padding: 1rem;
  }
  
  .header-content {
    padding: 1rem;
  }
  
  .logo h1 {
    font-size: 1.125rem;
  }
  
  .main-nav {
    gap: 0.25rem;
    overflow-x: auto;
  }
  
  .nav-item {
    font-size: 0.875rem;
    padding: 0.5rem 0.875rem;
  }
  
  .logout-btn {
    font-size: 0.875rem;
    padding: 0.5rem 0.875rem;
  }
}

@media (max-width: 640px) {
  .nav-text {
    font-size: 0.8125rem;
  }
  
  .main-nav {
    justify-content: flex-start;
  }
}
</style>