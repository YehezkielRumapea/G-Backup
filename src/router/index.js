import { createRouter, createWebHistory } from 'vue-router';
import DashboardView from '../views/DashboardView.vue';
import LoginView from '../views/LoginView.vue';
import RemoteView from '../views/RemoteView.vue';
import JobStatusView from '../views/JobStatusView.vue';
import ConfigurationView from '../views/ConfigurationView.vue';
import LogsView from '../views/LogsView.vue';


const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: LoginView
    },
    {
      path: '/',
      name: 'dashboard',
      component: DashboardView,
      meta: { requiresAuth: true }
    },
    {
      path: '/remote',
      name: 'remote',
      component: RemoteView,
      meta: { requiresAuth: true }
    },
    {
      path: '/job-status',
      name: 'job-status',
      component: JobStatusView,
      meta: { requiresAuth: true }
    },
    {
      path: '/configuration',
      name: 'configuration',
      component: ConfigurationView,
      meta: { requiresAuth: true }
    },
    {
      path: '/logs',
      name: 'logs',
      component: LogsView,
      meta: { requiresAuth: true }
    }
  ]
});

// Navigation Guard untuk memeriksa status login sebelum mengakses halaman
router.beforeEach((to, from, next) => {
  const isAuthenticated = localStorage.getItem('isAuthenticated') === 'true';

  if (to.meta.requiresAuth && !isAuthenticated) {
    next({ name: 'login' });
  } else if (to.name === 'login' && isAuthenticated) {
    next({ name: 'dashboard' });
  } else {
    next();
  }
});

export default router;

