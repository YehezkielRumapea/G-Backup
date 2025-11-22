<template>
 <div class="remotes-view">
  <div class="header">
   <div>
    <h1>Gdrive Monitoring</h1>
    <p class="subtitle">
Manage and monitor your Gdrive status</p>
   </div>
  </div>
  
  <div v-if="isLoading" class="status-message">
   <span class="loading-dot"></span>
   Memuat data...
  </div>
  
  <div v-if="errorMessage" class="status-message error">
   {{ errorMessage }}
  </div>

  <div v-if="!isLoading && remotes.length > 0" class="table-container">
   <table class="remotes-table">
    <thead>
     <tr>
      <th>Gdrive</th>
      <th>Status</th>
      <th>Storage</th>
      <th>Job Runs</th>
      <th>Last Check</th>
     </tr>
    </thead>
    <tbody>
     <RemoteRow
      v-for="remote in remotes"
      :key="remote.remote_name"
      :remote="remote"
     />
    </tbody>
   </table>
  </div>

  <div v-if="!isLoading && remotes.length === 0" class="empty-state">
   <p>Belum ada remote yang terdaftar</p>
  </div>
 </div>
</template>

<script setup>
import { ref, onMounted } from 'vue' 
import monitoringService from '@/services/monitoringService'
import RemoteRow from '@/components/RemoteRow.vue'

const remotes = ref([])
const isLoading = ref(true)
const errorMessage = ref(null)

onMounted(async () => {
 try {
  const data = await monitoringService.getRemoteStatus()
  remotes.value = data
 } catch (error) {
  errorMessage.value = 'Gagal memuat data. Silakan login ulang.'
 } finally {
  isLoading.value = false
 }
})
</script>

<style scoped>
.remotes-view {
 max-width: 1200px;
 margin: 0 auto;
 padding: 2rem 1.5rem;
}

.header {
 margin-bottom: 2.5rem;
}

h1 {
 font-size: 1.75rem;
 font-weight: 600;
 color: #1a1a1a;
 margin: 0 0 0.5rem 0;
 letter-spacing: -0.02em;
}

.subtitle {
 font-size: 0.95rem;
 color: #666;
 margin: 0;
 font-weight: 400;
}

.table-container {
 background: #fff;
 border: 1px solid #e5e5e5;
 border-radius: 8px;
 overflow: hidden;
}

.remotes-table {
 width: 100%;
 border-collapse: collapse;
}

.remotes-table th {
 background: #fafafa;
 padding: 0.875rem 1rem;
 text-align: left;
 font-size: 0.8125rem;
 font-weight: bold;
 color: #000000;
 text-transform: uppercase;
 letter-spacing: 0.05em;
 border-bottom: 1px solid #e5e5e5;
}

.remotes-table td {
 padding: 1rem;
 border-bottom: 1px solid #f0f0f0;
 font-size: 0.9375rem;
 color: #333;
}

.remotes-table tbody tr:last-child td {
 border-bottom: none;
}

.remotes-table tbody tr:hover {
 background: #fafafa;
}

.status-message {
 padding: 1rem;
 border-radius: 6px;
 font-size: 0.9375rem;
 background: #f8f8f8;
 color: #666;
 display: flex;
 align-items: center;
 gap: 0.75rem;
}

.status-message.error {
 background: #fef2f2;
 color: #dc2626;
 border: 1px solid #fee2e2;
}

.loading-dot {
 width: 8px;
 height: 8px;
 background: #666;
 border-radius: 50%;
 animation: pulse 1.5s ease-in-out infinite;
}

@keyframes pulse {
 0%, 100% { opacity: 1; }
 50% { opacity: 0.3; }
}

.empty-state {
 text-align: center;
 padding: 3rem 1rem;
 color: #999;
 font-size: 0.9375rem;
}

/* Responsive */
@media (max-width: 768px) {
 .remotes-view {
  padding: 1.5rem 1rem;
 }
 
 h1 {
  font-size: 1.5rem;
 }
 
 .table-container {
  overflow-x: auto;
 }
 
 .remotes-table {
  min-width: 600px;
 }
}
</style>