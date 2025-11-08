import { createRouter, createWebHistory } from 'vue-router';
import { useAuthStore } from '@/stores/authStore';

import Login from '@/views/Login.vue';
import Dashboard from '@/views/Dashboard.vue';
import Remotes from '@/views/Remotes.vue';
import Jobs from '@/views/Jobs.vue';
import Logs from '@/views/Logs.vue';
import CreateJob from '@/views/CreateJob.vue';

const routes = [
    {
        path: '/login',
        name: 'Login',
        component: Login
    },
    {
        path: '/',
        name: 'Dashboard',
        component: Dashboard,
        meta: {requiresAuth: true},

        children: [
            { path: '', redirect: '/remotes' }, // Halaman default dashboard
            { path: 'remotes', name: 'Remotes', component: Remotes },
            { path: 'scheduled-jobs', name: 'ScheduledJobs', component: ScheduledJobs },
            { path: 'manual-jobs', name: 'ManualJobs', component: ManualJobs },
        ]
    },
    {
        path: '/',
        redirect: '/login'        
    },
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

router.beforeEach((to, from, next) => {
    const authStore = useAuthStore()

    // jik login perlu login tapi user belum terotentikasi
    if ( to.meta.requiresAuth && !authStore.isAuthenticated() ) {
        next({name: '/login'})
    } else {
        next()
    }
})

export default router