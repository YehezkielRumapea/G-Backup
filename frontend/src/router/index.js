import { createRouter, createWebHistory } from 'vue-router';
import LoginView from '../views/LoginView.vue';
import DashboardView from '../views/DashboardView.vue';
import RemoteView from '../views/RemoteView.vue';
import JobStatusView from '../views/JobStatusView.vue';
import SchedulerView from '../views/SchedulerView.vue';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: LoginView
    },
    {
      path: '/',
      name: 'Dashboard',
      component: DashboardView,
      meta: { requiresAuth: true }
    },
    {
      path: '/remote',
      name: 'Remote',
      component: RemoteView,
      meta: { requiresAuth: true }
    },
    {
      path: '/job-status',
      name: 'JobStatus',
      component: JobStatusView,
      meta: { requiresAuth: true }
    },
	{
      path: '/scheduler',
      name: 'Scheduler',
      component: SchedulerView,
      meta: { requiresAuth: true }
    },
  ]
});

// Navigation Guard untuk memeriksa status login
router.beforeEach((to, from, next) => {
  const isAuthenticated = localStorage.getItem('isAuthenticated') === 'true';

  if (to.meta.requiresAuth && !isAuthenticated) {
    // Jika rute memerlukan login dan pengguna belum login, arahkan ke halaman login
    next('/login');
  } else if (to.name === 'Login' && isAuthenticated) {
    // Jika pengguna sudah login dan mencoba mengakses halaman login, arahkan ke dashboard
    next('/');
  } else {
    // Lanjutkan navigasi
    next();
  }
});

export default router;

