<script setup>
import { ref, computed, onMounted } from 'vue';
import apiClient from '@/api'; // Gunakan jembatan API

// State untuk menampung data dari backend
const jobs = ref([]);
const isLoading = ref(true);
const errorMsg = ref('');
const searchQuery = ref('');

// Fungsi untuk mengambil data job dari backend
const fetchScheduledJobs = async () => {
    isLoading.value = true;
    errorMsg.value = '';
    try {
        // Panggil endpoint yang sama dengan Job Status
        const response = await apiClient.get('/jobs/scheduled');
        jobs.value = response.data || [];
    } catch (error) {
        console.error("Gagal memuat data scheduler:", error);
         if (error.response?.status !== 401) {
            errorMsg.value = "Gagal memuat data dari server.";
         }
    } finally {
        isLoading.value = false;
    }
};

// Panggil fungsi ini saat komponen dimuat
onMounted(fetchScheduledJobs);

// --- FUNGSI UTAMA: MENJALANKAN JOB YANG SUDAH ADA ---
const runJobNow = async (jobId) => {
  if (!confirm(`Yakin ingin menjalankan ulang job ID: ${jobId} sekarang?`)) {
    return;
  }
  
  try {
    // Panggil endpoint API untuk menjalankan job (Trigger)
    // PERBAIKAN: Menggunakan endpoint /trigger/ yang benar
    const response = await apiClient.post(`/jobs/trigger/${jobId}`);
    
    alert(response.data.message || `Job ${jobId} berhasil dipicu!`);
    // Tunggu sejenak sebelum refresh agar backend punya waktu update status
    setTimeout(() => {
        fetchScheduledJobs();
    }, 2000); // Refresh setelah 2 detik
    
  } catch (error) {
    console.error("Gagal menjalankan job:", error.response?.data || error.message);
    alert("Gagal menjalankan job: " + (error.response?.data?.error || "Error tidak diketahui"));
  }
};

// --- COMPUTED & HELPERS (DENGAN PERBAIKAN) ---
const filteredJobs = computed(() => {
  if (!searchQuery.value) return jobs.value;
  // PERBAIKAN: Menggunakan 'job_name'
  return jobs.value.filter(job =>
    job.job_name?.toLowerCase().includes(searchQuery.value.toLowerCase())
  );
});

const getStatusClass = (status) => {
  status = status?.toUpperCase();
  if (status === 'SUCCESS' || status === 'COMPLETED') return 'status-success';
  if (status === 'FAIL' || status === 'FAILED' || status.startsWith('FAIL_')) return 'status-failed';
  if (status === 'RUNNING') return 'status-running';
  if (status === 'PENDING') return 'status-pending';
  return '';
};

const formatTimestamp = (timestamp) => {
    if (!timestamp || timestamp === '-' || timestamp.startsWith('0001-01-01')) return '-';
    try { 
      // API Anda mengembalikan format yang berbeda di 'next_run'
      // jadi kita tangani keduanya
      return new Date(timestamp).toLocaleString('id-ID', { hour12: false });
    } catch(e) { 
      return timestamp; 
    }
};
</script>

<template>
  <div class="scheduler-container">
    <header class="main-header">
      <div class="header-content">
        <div>
          <h1>Scheduler</h1>
          <p>Monitor and manage job execution schedules.</p>
        </div>
      </div>
    </header>

    <div class="card">
      <div class="toolbar">
        <div class="search-wrapper">
           <svg class="search-icon" xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16"><path d="M11.742 10.344a6.5 6.5 0 1 0-1.397 1.398h-.001c.03.04.062.078.098.115l3.85 3.85a1 1 0 0 0 1.415-1.414l-3.85-3.85a1.007 1.007 0 0 0-.115-.1zM12 6.5a5.5 5.5 0 1 1-11 0 5.5 5.5 0 0 1 11 0z"/></svg>
           <input type="text" v-model="searchQuery" placeholder="Search Job" class="search-input">
        </div>
      </div>

      <div v-if="isLoading" class="text-center p-8">Memuat data jadwal...</div>
      <div v-else-if="errorMsg" class="error-message p-4">{{ errorMsg }}</div>

      <div v-else class="table-responsive">
        <table class="data-table">
          <thead>
            <tr>
              <th>Job Name</th>
              <th>Last Run</th>
              <th>Status</th>
              <th>Next Run</th>
              <th>Action</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="!filteredJobs || filteredJobs.length === 0">
              <td colspan="5" class="text-center">Tidak ada job terjadwal ditemukan.</td>
            </tr>
            <tr v-for="job in filteredJobs" :key="job.id">
              <td><strong>{{ job.job_name }}</strong><br><small>{{ job.type }} to {{ job.gdrive_target }}</small></td>
              <td>{{ formatTimestamp(job.last_run) }}</td>
              <td>
                <span :class="['status-badge', getStatusClass(job.status)]">
                  {{ job.status }}
                </span>
              </td>
              <td>{{ formatTimestamp(job.next_run) }}</td>
              <td>
                <div class="action-buttons">
                  <button @click="runJobNow(job.id)" class="action-btn run-btn" title="Run Job Now">
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polygon points="5 3 19 12 5 21 5 3"></polygon></svg>
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
.data-table td small { color: #6c757d; font-size: 0.8rem; }
.status-badge { display: inline-block; padding: 4px 10px; border-radius: 4px; font-weight: 500; font-size: 0.85rem; color: white; min-width: 70px; text-align: center; text-transform: capitalize; }
.status-success { background-color: #28a745; }
.status-failed { background-color: #dc3545; }
.status-running { background-color: #0d6efd; }
.status-pending { background-color: #ffc107; color: #333;}
.action-buttons { display: flex; gap: 8px; }
.action-btn { border: none; background: none; cursor: pointer; padding: 5px; border-radius: 4px; display: flex; align-items: center; justify-content: center; }
.run-btn { color: #28a745; }
.text-center { text-align: center; }
.error-message { color: #dc3545; }
</style>