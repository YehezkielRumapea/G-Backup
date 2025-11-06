<script setup>
import { ref, computed, onMounted } from 'vue';
import apiClient from '@/api'; // 1. Impor apiClient

// --- STATE ---
const jobs = ref([]); // 2. Mulai dengan array kosong
const isLoading = ref(true);
const errorMsg = ref('');
const searchQuery = ref('');

// --- FUNGSI FETCH DATA ---
const fetchJobs = async () => {
    isLoading.value = true;
    errorMsg.value = '';
    try {
        // 3. Panggil API jobs/scheduled
        const response = await apiClient.get('/jobs/scheduled');
        jobs.value = response.data || [];
    } catch (error) {
        console.error("Gagal memuat data jobs:", error);
         if (error.response?.status !== 401) {
            errorMsg.value = "Gagal memuat data jobs dari server.";
         }
    } finally {
        isLoading.value = false;
    }
};

// 4. Panggil API saat komponen dimuat
onMounted(fetchJobs);

// --- FUNGSI AKSI ---
const runJob = async (jobId) => {
  if (!confirm(`Yakin ingin menjalankan job ID: ${jobId} sekarang?`)) return;
  
  try {
    // 5. Panggil API 'trigger' yang benar
    const response = await apiClient.post(`/jobs/trigger/${jobId}`);
    alert(response.data.message || `Job ${jobId} berhasil dipicu!`);
    // Tunggu 2 detik lalu refresh data untuk melihat status 'RUNNING'
    setTimeout(fetchJobs, 2000); 
  } catch (error) {
     console.error("Gagal menjalankan job:", error.response?.data || error.message);
     alert("Gagal menjalankan job: " + (error.response?.data?.error || "Error tidak diketahui"));
  }
}

const deleteJob = (jobId) => {
  alert(`Fungsi "Delete Job" (ID: ${jobId}) belum terhubung ke backend.`);
  // (Anda perlu membuat endpoint DELETE /api/v1/jobs/:id di backend nanti)
}

// --- FUNGSI HELPER ---
const filteredJobs = computed(() => {
  if (!searchQuery.value) return jobs.value;
  return jobs.value.filter(job =>
    job.job_name.toLowerCase().includes(searchQuery.value.toLowerCase())
  );
});

const getStatusClass = (status) => {
  status = status?.toUpperCase();
  if (status === 'SUCCESS' || status === 'COMPLETED') return 'status-success';
  if (status === 'FAIL' || status === 'FAILED' || status.startsWith('FAIL_')) return 'status-failed';
  if (status === 'RUNNING') return 'status-running';
  return 'status-pending'; // Default untuk 'PENDING' dll
};

const formatTimestamp = (timestamp) => {
    if (!timestamp || timestamp === '-') return '-';
    try { return new Date(timestamp).toLocaleString('id-ID'); } catch(e) { return timestamp; }
};
</script>

<template>
  <div class="jobs-container">
    <header class="main-header">
      <div class="header-content">
        <div>
          <h1>Jobs</h1>
          <p>See your job status</p>
        </div>
      </div>
    </header>

    <div class="card">
      <div class="toolbar">
        <div class="search-wrapper">
           <svg class="search-icon" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line></svg>
           <input type="text" v-model="searchQuery" placeholder="Search Job" class="search-input">
        </div>
      </div>

      <div class="table-responsive">
        <table class="data-table">
          <thead>
            <tr>
              <th>Job Name</th>
              <th>Type</th>
              <th>Gdrive Target</th>
              <th>Mode</th>
              <th>Last Run</th>
              <th>Status</th>
              <th>Next Run</th>
              <th>Action</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="isLoading">
              <td colspan="8" class="text-center">Memuat data jobs...</td>
            </tr>
            <tr v-else-if="errorMsg">
              <td colspan="8" class="text-center" style="color: red;">{{ errorMsg }}</td>
            </tr>
            <tr v-else-if="filteredJobs.length === 0">
              <td colspan="8" class="text-center">No jobs found.</td>
            </tr>
            <tr v-for="job in filteredJobs" :key="job.id">
              <td><strong>{{ job.job_name }}</strong><br><small>{{ job.type }}</small></td>
              <td>{{ job.type }}</td>
              <td>{{ job.gdrive_target }}</td>
              <td>{{ job.mode }}</td>
              <td>{{ formatTimestamp(job.last_run) }}</td>
              <td>
                <span :class="['status-badge', getStatusClass(job.status)]">
                  {{ job.status }}
                </span>
              </td>
              <td>{{ formatTimestamp(job.next_run) }}</td>
              <td>
                <div class="action-buttons">
                  <button @click="runJob(job.id)" class="action-btn run-btn" title="Run Job">
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polygon points="5 3 19 12 5 21 5 3"></polygon></svg>
                  </button>
                  <button @click="deleteJob(job.id)" class="action-btn delete-btn" title="Delete Job">
                     <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path><line x1="10" y1="11" x2="10" y2="17"></line><line x1="14" y1="11" x2="14" y2="17"></line></svg>
                  </button>
                </div>
              </td>
            </tr>
            </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<style scoped>
.main-header { margin-bottom: 20px; }
.header-content { display: flex; justify-content: space-between; align-items: center; }
.main-header h1 { font-size: 2rem; font-weight: 700; color: #333; }
.main-header p { color: #6c757d; }
.card { background-color: #fff; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.05); border: 1px solid var(--border-color); padding: 20px; }
.toolbar { margin-bottom: 20px; }
.search-wrapper { position: relative; max-width: 300px; }
.search-input { width: 100%; padding: 10px 15px 10px 40px; border: 1px solid var(--border-color); border-radius: 6px; font-size: 1rem; }
.search-icon { position: absolute; left: 12px; top: 50%; transform: translateY(-50%); color: #6c757d; }
.table-responsive { overflow-x: auto; }
.data-table { width: 100%; border-collapse: collapse; }
.data-table th, .data-table td { padding: 12px 15px; text-align: left; border-bottom: 1px solid var(--border-color); vertical-align: middle; white-space: nowrap; }
.data-table th { background-color: #f8f9fa; font-weight: 600; color: #495057; }
.data-table td small { color: #6c757d; }
.status-badge { display: inline-block; padding: 4px 10px; border-radius: 4px; font-weight: 500; font-size: 0.85rem; color: white; min-width: 70px; text-align: center; }
.status-success { background-color: #28a745; }
.status-failed { background-color: #dc3545; }
.status-running { background-color: #007bff; }
.status-pending { background-color: #ffc107; color: #333;}
.action-buttons { display: flex; gap: 8px; }
.action-btn { border: none; background: none; cursor: pointer; padding: 5px; border-radius: 4px; display: flex; align-items: center; justify-content: center; transition: background-color 0.2s; }
.action-btn:hover { background-color: #e9ecef; }
.run-btn { color: #28a745; }
.delete-btn { color: #dc3545; }
.text-center { text-align: center; }
</style>