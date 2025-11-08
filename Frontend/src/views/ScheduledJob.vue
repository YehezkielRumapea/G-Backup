<template>
  <div class="scheduled-jobs-view">
    <h1>Scheduled Jobs</h1>
    <p>Memantau semua Job yang berjalan otomatis (CRON).</p>
    
    <router-link to="/create" class="btn-add-job">
      + Buat Job Baru
    </router-link>

    <div v-if="isLoading" class="loading">
      Memuat data scheduled jobs...
    </div>
    
    <div v-if="errorMessage" class="error">
      {{ errorMessage }}
    </div>

    <table v-if="jobs.length > 0" class="jobs-table">
      <thead>
        <tr>
          <th>Nama Job</th>
          <th>Tipe</th>
          <th>Target GDrive</th>
          <th>Mode</th>
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
        />
      </tbody>
    </table>
    
    <div v-else-if="!isLoading && jobs.length === 0">
      <p>Tidak ada job terjadwal yang ditemukan.</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import monitoringService from '@/services/monitoringService' // Service untuk GET data
import jobService from '@/services/jobService' // Service untuk Trigger (Action)
import JobRow from '@/components/JobRow.vue' // Komponen baris

const jobs = ref([])
const isLoading = ref(true)
const errorMessage = ref(null)

// Fungsi untuk memuat data dari backend
async function fetchData() {
  isLoading.value = true
  errorMessage.value = null
  try {
    // Panggil API backend Golang (GET /api/v1/monitoring/jobs)
    const data = await monitoringService.getScheduledJobs()
    jobs.value = data
  } catch (error) {
    errorMessage.value = 'Gagal memuat data scheduled jobs.'
  } finally {
    isLoading.value = false
  }
}

// Panggil fetchData() saat komponen (halaman) dimuat
onMounted(fetchData)

// --- ACTIONS ---

// Dipanggil oleh tombol "Run Now" (Play) dari JobRow
async function handleTrigger(jobId) {
  if (!confirm(`Apakah Anda yakin ingin menjalankan job ID ${jobId} sekarang?`)) {
    return
  }
  
  try {
    // Panggil API backend Golang (POST /api/v1/jobs/trigger/:id)
    await jobService.triggerJob(jobId)
    alert(`Job ${jobId} berhasil di-trigger!`)
    // Muat ulang data untuk melihat status "RUNNING"
    fetchData() 
  } catch (error) {
    alert(`Gagal men-trigger job: ${error}`)
  }
}

// Dipanggil oleh tombol "View Script"
async function handleViewScript(jobId) {
  try {
    // Panggil API backend Golang (GET /api/v1/jobs/script/:id)
    const data = await jobService.getJobScript(jobId)
    // Tampilkan script di alert (atau modal)
    alert(`Script untuk Job ${jobId}:\n\n${data.script_preview}`)
  } catch (error) {
    alert(`Gagal memuat script: ${error}`)
  }
}
</script>

<style scoped>
/* Styling untuk tabel (mirip dengan Remotes.vue) */
.jobs-table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 1.5rem;
  font-size: 0.9rem;
}
.jobs-table th, .jobs-table td {
  border-bottom: 1px solid #ddd;
  padding: 10px 12px;
  text-align: left;
}
.jobs-table th {
  background-color: #f4f4f4;
}
.loading, .error {
  margin-top: 1rem;
}
.error {
  color: red;
}
.btn-add-job {
  /* ... styling (mirip btn-add-remote) ... */
  display: inline-block;
  background-color: #1abc9c;
  color: white;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  text-decoration: none;
  font-weight: bold;
}
</style>