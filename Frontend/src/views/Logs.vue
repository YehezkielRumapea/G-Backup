<template>
  <div class="logs-view">
    <h1>Execution Logs</h1>
    <p>Riwayat semua Job (Auto, Manual, Restore) yang telah dieksekusi.</p>

    <div v-if="isLoading" class="loading">
      Memuat data logs...
    </div>
    
    <div v-if="errorMessage" class="error">
      {{ errorMessage }}
    </div>

    <table v-if="logs.length > 0" class="logs-table">
      <thead>
        <tr>
          <th>Job ID</th>
          <th>Nama Job (dari Snapshot/Relasi)</th>
          <th>Status</th>
          <th>Pesan Output (CLI)</th>
          <th>Durasi (Detik)</th>
          <th>Timestamp</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="log in logs" :key="log.ID">
          <td>{{ log.JobID || 'N/A (Manual)' }}</td>
          <td>
            <strong>{{ getJobName(log) }}</strong>
          </td>
          <td>
            <span class="status" :class="log.Status.toLowerCase()">
              {{ log.Status }}
            </span>
          </td>
          <td class="message-cell">
            <pre>{{ log.Message }}</pre> </td>
          <td>{{ log.DurationSec }} detik</td>
          <td>{{ formatTimestamp(log.Timestamp) }}</td>
        </tr>
      </tbody>
    </table>
    
    <div v-else-if="!isLoading && logs.length === 0">
      <p>Belum ada eksekusi Job yang tercatat.</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import monitoringService from '@/services/monitoringService' // Service untuk GET data

const logs = ref([])
const isLoading = ref(true)
const errorMessage = ref(null)

// Fungsi untuk memuat data dari backend
async function fetchData() {
  isLoading.value = true
  errorMessage.value = null
  try {
    // Panggil API backend Golang (GET /api/v1/monitoring/logs)
    const data = await monitoringService.getLogs()
    logs.value = data
  } catch (error) {
    errorMessage.value = 'Gagal memuat data logs.'
  } finally {
    isLoading.value = false
  }
}

// Panggil fetchData() saat komponen (halaman) dimuat
onMounted(fetchData)

// --- Helper Functions ---

// Mengambil nama Job dari relasi (jika ada) atau snapshot (jika manual)
function getJobName(log) {
  if (log.ScheduledJob && log.ScheduledJob.JobName) {
    return log.ScheduledJob.JobName
  }
  if (log.ConfigSnapshot) {
    // Asumsi ConfigSnapshot adalah JSON string, kita perlu parse
    try {
      const config = JSON.parse(log.ConfigSnapshot)
      return config.job_name || 'Manual Job (No Name)'
    } catch (e) {
      return 'Manual Job (Invalid Snapshot)'
    }
  }
  return 'N/A'
}

// Memformat timestamp
function formatTimestamp(isoString) {
  const date = new Date(isoString)
  return date.toLocaleString('id-ID') // Format ke waktu lokal Indonesia
}
</script>

<style scoped>
/* Styling untuk tabel (mirip dengan Jobs.vue) */
.logs-table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 1.5rem;
  font-size: 0.9rem;
}
.logs-table th, .logs-table td {
  border-bottom: 1px solid #ddd;
  padding: 10px 12px;
  text-align: left;
  vertical-align: top;
}
.logs-table th {
  background-color: #f4f4f4;
}
.loading, .error {
  margin-top: 1rem;
}
.error {
  color: red;
}
/* Style untuk output CLI */
.message-cell pre {
  white-space: pre-wrap; /* Jaga format newline */
  word-break: break-all;
  background: #fdfdfd;
  padding: 5px;
  border: 1px solid #eee;
  max-height: 100px;
  overflow-y: auto;
  font-family: monospace;
  font-size: 0.8rem;
  color: #555;
}

/* Style untuk Status */
.status {
  padding: 4px 8px;
  border-radius: 12px;
  font-weight: bold;
  font-size: 0.8rem;
  color: white;
  text-transform: uppercase;
}
.status.pending { background-color: #aaa; }
.status.running { background-color: #3498db; }
.status.completed, .status.success { background-color: #2ecc71; }
.status.failed, .status.fail, .status.fail_pre_script, .status.fail_rclone, .status.fail_post_script { background-color: #e74c3c; }
</style>