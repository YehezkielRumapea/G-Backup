<template>
  <div class="manual-jobs-view">
    <h1>Manual Jobs (Templates)</h1>
    <p>Menjalankan Job (Backup/Restore) yang telah disimpan tanpa jadwal.</p>
    
    <div v-if="isLoading" class="loading">
      Memuat data manual jobs...
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
          <th>Action</th>
        </tr>
      </thead>
      <tbody>
        <JobRow
          v-for="job in jobs"
          :key="job.id"
          :job="job"
          :show-next-run="false" @trigger="handleTrigger"
          @view-script="handleViewScript"
        />
      </tbody>
    </table>
    
    <div v-else-if="!isLoading && jobs.length === 0">
      <p>Tidak ada template job manual yang ditemukan. Buat di halaman "Create Job".</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import jobService from '@/services/jobService' // Service untuk GET data dan Trigger
import JobRow from '@/components/JobRow.vue' // Komponen baris

const jobs = ref([])
const isLoading = ref(true)
const errorMessage = ref(null)

// Fungsi untuk memuat data dari backend
async function fetchData() {
  isLoading.value = true
  errorMessage.value = null
  try {
    // PENTING: Panggil API backend Golang (GET /api/v1/jobs/manual)
    // Anda perlu membuat endpoint ini di backend (JobHandler)
    // yang memanggil JobRepo.FindAllManualJobs()
    const data = await jobService.getManualJobs() // (Asumsi method ini ada)
    jobs.value = data
  } catch (error) {
    errorMessage.value = 'Gagal memuat data manual jobs.'
  } finally {
    isLoading.value = false
  }
}

// Panggil fetchData() saat komponen (halaman) dimuat
onMounted(fetchData)

// --- ACTIONS ---

// Dipanggil oleh tombol "Run Now" (Play) dari JobRow
async function handleTrigger(jobId) {
  if (!confirm(`Apakah Anda yakin ingin menjalankan job manual ID ${jobId} sekarang?`)) {
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
    alert(`Script untuk Job ${jobId}:\n\n${data.script_preview}`)
  } catch (error) {
    alert(`Gagal memuat script: ${error}`)
  }
}
</script>

<style scoped>
/* Styling (mirip dengan ScheduledJobs.vue) */
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
</style>