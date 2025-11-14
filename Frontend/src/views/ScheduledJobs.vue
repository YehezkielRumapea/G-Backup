<template>
  <div class="scheduled-jobs-view">
    <div class="header">
      <div>
        <h1>Scheduled Jobs</h1>
        <p class="subtitle">Memantau semua Job yang berjalan otomatis</p>
      </div>
    </div>
    
    <ScriptPreview
      :is-visible="isModalVisible"
      :job-id="currentJobId"
      :script-content="currentScript"
      @close="closeModal"
    />

    <div v-if="isLoading" class="status-message">
      <span class="loading-dot"></span>
      Memuat data...
    </div>
    
    <div v-if="errorMessage" class="status-message error">
      {{ errorMessage }}
      <button @click="fetchData" class="retry-btn">Coba Lagi</button>
    </div>

    <div v-if="!isLoading && jobs.length === 0 && !errorMessage" class="empty-state">
      <p>Tidak ada job terjadwal</p>
      <router-link to="/create" class="btn-create">
        Buat Job Terjadwal
      </router-link>
    </div>

    <div v-if="!isLoading && jobs.length > 0" class="table-container">
      <table class="jobs-table">
        <thead>
          <tr>
            <th>Nama Job</th>
            <th>Object</th>
            <th>Target GDrive</th>
            <th>Last Run</th>
            <th>Status</th>
            <th>Next Run</th>
            <th>Action</th>
          </tr>
        </thead>
        <tbody>
          <JobRow
            v-for="job in jobs"
            :key="job.id"
            :job="job"
            @trigger="handleTrigger"
            @view-script="handleViewScript"
            @delete="handleDeleteJob"
          />
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import monitoringService from '@/services/monitoringService';
import jobService from '@/services/jobService';
import JobRow from '@/components/SchedJobRow.vue';
import ScriptPreview from '@/components/ScriptPreview.vue';

const jobs = ref([]);
const isLoading = ref(true);
const errorMessage = ref(null);

const isModalVisible = ref(false);
const currentScript = ref('');
const currentJobId = ref(null);

async function fetchData() {
  isLoading.value = true;
  errorMessage.value = null;
  
  try {
    const data = await monitoringService.getScheduledJobs();
    jobs.value = data;
  } catch (error) {
    errorMessage.value = error.response?.data?.error || 'Gagal memuat data scheduled jobs.';
    console.error("Fetch Scheduled Jobs Error:", error);
  } finally {
    isLoading.value = false;
  }
}

onMounted(fetchData);

async function handleTrigger(jobId) {
  if (!confirm(`Apakah Anda yakin ingin menjalankan job ID ${jobId} sekarang?`)) {
    return;
  }
  
  try {
    await jobService.triggerJob(jobId);
    alert(`Job ${jobId} berhasil di-trigger! Status akan diperbarui.`);
    fetchData();
  } catch (error) {
    const errorMsg = error.response?.data?.error || error.message || String(error);
    alert(`Gagal men-trigger job: ${errorMsg}`);
  }
}

async function handleViewScript(jobId) {
  try {
    isModalVisible.value = true;
    currentScript.value = 'Loading script...';
    currentJobId.value = jobId;
    
    const data = await jobService.getJobScript(jobId);
    currentScript.value = data.script_preview || data.script || 'No script available';
  } catch (error) {
    isModalVisible.value = false;
    const errorMsg = error.response?.data?.error || error.message;
    alert(`Gagal memuat script: ${errorMsg}`);
  }
}

async function handleDeleteJob(jobId) {
    try {
        await jobService.deleteJob(jobId); 
        alert(`Job ID ${jobId} berhasil dihapus.`);
        fetchData();
    } catch (error) {
        const errorMsg = error.response?.data?.error || "Gagal menghapus job.";
        alert(`Gagal menghapus job: ${errorMsg}`);
    }
}

function closeModal() {
  isModalVisible.value = false;
  currentScript.value = '';
  currentJobId.value = null;
}
</script>

<style scoped>
.scheduled-jobs-view {
  max-width: 1400px;
  margin: 0 auto;
  padding: 2rem 1.5rem;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 2.5rem;
  gap: 1rem;
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

.btn-add {
  display: inline-flex;
  align-items: center;
  background: #1a1a1a;
  color: white;
  padding: 0.625rem 1.125rem;
  border-radius: 6px;
  text-decoration: none;
  font-weight: 500;
  font-size: 0.9375rem;
  transition: all 0.2s;
  white-space: nowrap;
}

.btn-add:hover {
  background: #333;
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
  justify-content: space-between;
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

.retry-btn {
  background: transparent;
  border: 1px solid #dc2626;
  color: #dc2626;
  padding: 0.375rem 0.875rem;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.875rem;
  transition: all 0.2s;
}

.retry-btn:hover {
  background: #dc2626;
  color: white;
}

.empty-state {
  text-align: center;
  padding: 3rem 1rem;
}

.empty-state p {
  color: #999;
  font-size: 0.9375rem;
  margin-bottom: 1.5rem;
}

.btn-create {
  display: inline-block;
  background: #1a1a1a;
  color: white;
  padding: 0.625rem 1.125rem;
  border-radius: 6px;
  text-decoration: none;
  font-weight: 500;
  font-size: 0.9375rem;
  transition: all 0.2s;
}

.btn-create:hover {
  background: #333;
}

.table-container {
  background: #fff;
  border: 1px solid #e5e5e5;
  border-radius: 8px;
  overflow: hidden;
}

.jobs-table {
  width: 100%;
  border-collapse: collapse;
}

.jobs-table thead {
  background: #fafafa;
}

.jobs-table th {
  padding: 0.875rem 1rem;
  text-align: left;
  font-size: 0.8125rem;
  font-weight: bold;
  color: #000000;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  border-bottom: 1px solid #e5e5e5;
}

.jobs-table td {
  padding: 1rem;
  border-bottom: 1px solid #f0f0f0;
  font-size: 0.9375rem;
  color: #333;
}

.jobs-table tbody tr:last-child td {
  border-bottom: none;
}

.jobs-table tbody tr:hover {
  background: #fafafa;
}

@media (max-width: 768px) {
  .scheduled-jobs-view {
    padding: 1.5rem 1rem;
  }
  
  .header {
    flex-direction: column;
    align-items: stretch;
  }
  
  h1 {
    font-size: 1.5rem;
  }
  
  .btn-add {
    justify-content: center;
  }
  
  .table-container {
    overflow-x: auto;
  }
  
  .jobs-table {
    min-width: 800px;
  }
}
</style>