import { createRouter, createWebHistory } from 'vue-router';
import { useAuthStore } from '@/stores/authStore.js';

// ============================================
// IMPORT LAYOUT
// ============================================
import AppLayout from '@/views/AppLayout.vue';

// ============================================
// IMPORT PAGES
// ============================================
import Login from '@/views/Login.vue';
import Dashboard from '@/views/Dashboard.vue';
import CreateJob from '@/views/CreateJob.vue';
import Logs from '@/views/Logs.vue';
import ManualJobs from '@/views/ManualJobs.vue';
import Remotes from '@/views/Remotes.vue';
import ScheduledJobs from '@/views/ScheduledJobs.vue';

// ============================================
// ROUTE CONFIGURATION
// ============================================
const routes = [
    // ----------------------------------------
    // PUBLIC ROUTES (No Authentication)
    // ----------------------------------------
    {
        path: '/login',
        name: 'Login',
        component: Login,
        meta: { 
            requiresAuth: false,
            title: 'Login - G-Backup'
        }
    },
    
    // ----------------------------------------
    // PROTECTED ROUTES (With Authentication)
    // ----------------------------------------
    {
        path: '/',
        component: AppLayout,  // ⭐ LAYOUT WRAPPER
        meta: { requiresAuth: true },
        children: [
            // Default redirect
            {
                path: '',
                redirect: { name: 'Dashboard' }
            },
            
            // Dashboard
            {
                path: 'dashboard',
                name: 'Dashboard',
                component: Dashboard,
                meta: { 
                    title: 'Dashboard - G-Backup',
                    icon: '📊'
                }
            },
            
            // Remotes
            {
                path: 'remotes',
                name: 'Remotes',
                component: Remotes,
                meta: { 
                    title: 'Remotes - G-Backup',
                    icon: '☁️'
                }
            },
            
            // Scheduled Jobs
            {
                path: 'scheduled',
                name: 'ScheduledJobs',
                component: ScheduledJobs,
                meta: { 
                    title: 'Scheduled Jobs - G-Backup',
                    icon: '⏰'
                }
            },
            
            // Manual Jobs
            {
                path: 'manual',
                name: 'ManualJobs',
                component: ManualJobs,
                meta: { 
                    title: 'Manual Jobs - G-Backup',
                    icon: '🔧'
                }
            },
            
            // Create Job
            {
                path: 'create',
                name: 'CreateJob',
                component: CreateJob,
                meta: { 
                    title: 'Create Job - G-Backup',
                    icon: '➕'
                }
            },
            
            // Logs
            {
                path: 'logs',
                name: 'Logs',
                component: Logs,
                meta: { 
                    title: 'Logs - G-Backup',
                    icon: '📝'
                }
            }
        ]
    },
    
    // ----------------------------------------
    // CATCH ALL / 404
    // ----------------------------------------
    {
        path: '/:pathMatch(.*)*',
        name: 'NotFound',
        redirect: { name: 'Login' }
    }
];

// ============================================
// CREATE ROUTER INSTANCE
// ============================================
const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes,
    // Scroll behavior
    scrollBehavior(to, from, savedPosition) {
        if (savedPosition) {
            return savedPosition;
        } else {
            return { top: 0 };
        }
    }
});

// ============================================
// NAVIGATION GUARDS
// ============================================
router.beforeEach((to, from, next) => {
    // Debug logs (bisa dihapus di production)
    console.log('🔍 Navigation:', {
        from: from.path,
        to: to.path,
        requiresAuth: to.meta.requiresAuth
    });
    
    const authStore = useAuthStore();
    const isAuthenticated = authStore.isAuthenticated();
    
    console.log('🔐 Auth Status:', isAuthenticated);
    
    // Set document title
    if (to.meta.title) {
        document.title = to.meta.title;
    }
    
    // Check authentication
    if (to.meta.requiresAuth && !isAuthenticated) {
        console.log('❌ Not authenticated, redirect to login');
        next({ name: 'Login', query: { redirect: to.fullPath } });
        return;
    }
    
    // Prevent authenticated users from accessing login
    if (to.name === 'Login' && isAuthenticated) {
        console.log('✅ Already authenticated, redirect to dashboard');
        next({ name: 'Dashboard' });
        return;
    }
    
    console.log('✅ Navigation allowed');
    next();
});

// Optional: After navigation
router.afterEach((to, from) => {
    console.log('✅ Navigated to:', to.path);
});

export default router;