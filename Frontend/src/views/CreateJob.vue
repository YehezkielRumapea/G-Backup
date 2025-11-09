<template>
  <div class="config-view">
    
    <form @submit.prevent="handleBackupSubmit" class="config-form">
      <h2>Backup Config (Mode Script Runner)</h2>
      <p>Buat template Job baru (Manual atau Terjadwal). Logika (seperti mysqldump atau zip) harus disediakan di Pre-Script.</p>

      <div>
        <label for="backup-jobName">Job Name (Nama Pekerjaan)</label>
        <input type="text" id="backup-jobName" v-model="backupForm.job_name" required />
      </div>

      <div>
        <label for="backup-source">Source Path (Lokal)</label>
        <input type="text" id="backup-source" v-model="backupForm.source_path" required placeholder="Contoh: /tmp/file_hasil_prescript.zip" />
      </div>
      <div>
        <label for="backup-remote">Remote Name</label>
        <input type="text" id="backup-remote" v-model="backupForm.remote_name" required placeholder="Gdrive1" />
      </div>
      <div>
        <label for="backup-dest">Destination Path (Cloud)</label>
        <input type="text" id="backup-dest" v-model="backupForm.destination_path" required placeholder="folder_backup/sub_folder" />
      </div>

      <div>
        <label for="backup-cron">Schedule (CRON) - (Kosongkan untuk Job Manual)</label>
        <input type="text" id="backup-cron" v-model="backupForm.schedule_cron" placeholder="Contoh: */5 * * * * (Tiap 5 menit)" />
      </div>

      <div>
        <label for="backup-pre">Pre-Script (Dijalankan SEBELUM Rclone)</label>
        <textarea id="backup-pre" v-model="backupForm.pre_script" rows="5" placeholder="Contoh: mysqldump -uuser -ppass db > /tmp/file.sql.gz"></textarea>
      </div>

      <div>
        <label for="backup-post">Post-Script (Dijalankan SETELAH Rclone sukses)</label>
        <textarea id="backup-post" v-model="backupForm.post_script" rows="4" placeholder="Contoh: rm /tmp/file.sql.gz"></textarea>
      </div>

      <button type="submit" :disabled="isLoading" class="btn-submit">
        {{ backupForm.schedule_cron ? 'Simpan Job Terjadwal' : 'Simpan & Jalankan Job Manual' }}
      </button>
    </form>

    <hr class="divider" />

    <form @submit.prevent="handleRestoreSubmit" class="config-form restore-form">
      <h2>Restore Config (Simpel)</h2>
      <p>Pilih file dari Remote dan restore ke path lokal (Mode ini tidak menggunakan Pre/Post Script).</p>
      
      <div>
        <label for="restore-remote">Remote Name (Sumber)</label>
        <input type="text" id="restore-remote" v-model="restoreForm.remote_name" required placeholder="Gdrive1" />
      </div>

      <div>
        <label for="restore-source">Source Path (Path di Cloud)</label>
        <input type="text" id="restore-source" v-model="restoreForm.source_path" required placeholder="folder_backup/file.zip" />
      </div>

      <div>
        <label for="restore-dest">Destination Path (Path di Server Lokal)</label>
        <input type="text" id="restore-dest" v-model="restoreForm.destination_path" required placeholder="/home/user/restore_here" />
      </div>
      
      <button type="submit" :disabled="isLoading" class="btn-submit btn-restore">Mulai Restore</button>
    </form>

    <div vMonitor="message" class="message" v-if="message">{{ message }}</div>
    <div vMonitor="errorMessage" class="error" v-if="errorMessage">{{ errorMessage }}</div>

  </div>
</template>

<script setup>
import { ref } from 'vue'
import jobService from '@/services/jobService' // Import service API
import { useRouter } from 'vue-router'

const router = useRouter()
const isLoading = ref(false)
const errorMessage = ref(null)
const message = ref(null)

// State untuk Form Backup (DTO)
const backupForm = ref({
  job_name: '',
  rclone_mode: 'COPY',
  source_path: '',
  remote_name: '',
  destination_path: '',
  schedule_cron: '',
  pre_script: '',
  post_script: ''
})

// State untuk Form Restore (DTO)
const restoreForm = ref({
  remote_name: '',
  source_path: '',
  destination_path: ''
})

// --- ACTIONS ---

async function handleBackupSubmit() {
  isLoading.value = true
  errorMessage.value = null
  message.value = null
  
  try {
    // Panggil API backend Golang (POST /api/v1/jobs/new)
    const response = await jobService.createBackupJob(backupForm.value)
    message.value = response.message // "Job berhasil diterima..."
    
    // Arahkan user ke halaman yang relevan
    if (backupForm.value.schedule_cron === '') {
      router.push('/manual-jobs')
    } else {
      router.push('/scheduled-jobs')
    }
    
  } catch (error) {
    errorMessage.value = 'Gagal membuat job backup. Cek kembali input Anda.'
  } finally {
    isLoading.value = false
  }
}

async function handleRestoreSubmit() {
  isLoading.value = true
  errorMessage.value = null
  message.value = null

  try {
    // Panggil API backend Golang (POST /api/v1/jobs/restore)
    const response = await jobService.createRestoreJob(restoreForm.value)
    message.value = response.message // "Restore Job diterima..."
    // Arahkan ke halaman Logs untuk melihat progres
    router.push('/logs') 
  } catch (error) {
    errorMessage.value = 'Gagal memulai restore.'
  } finally {
    isLoading.value = false
  }
}
</script>

<style scoped>
.config-view {
  max-width: 800px;
}
.config-form {
  background: #fff;
  padding: 1.5rem;
  border-radius: 8px;
  border: 1px solid #ddd;
}
.restore-form {
  border-top: 2px solid #3498db; /* Biru */
}
.divider {
  margin: 2rem 0;
  border: 0;
  border-top: 1px solid #eee;
}
.config-form div {
  margin-bottom: 1rem;
}
.config-form label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 600;
}
.config-form input[type="text"],
.config-form select,
.config-form textarea {
  width: 100%;
  padding: 0.5rem;
  border: 1px solid #ccc;
  border-radius: 4px;
  box-sizing: border-box; /* Mencegah padding merusak layout */
}
.config-form textarea {
  font-family: monospace;
  font-size: 0.9rem;
}
.btn-submit {
  background-color: #1abc9c; /* Hijau */
  color: white;
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-weight: bold;
}
.btn-submit.btn-restore {
  background-color: #3498db; /* Biru */
}
.btn-submit:disabled {
  background-color: #bdc3c7; /* Abu-abu */
}
.message, .error {
  margin-top: 1rem;
  padding: 1rem;
  border-radius: 4px;
}
.message {
  background-color: #dff0d8;
  color: #3c763d;
}
.error {
  background-color: #f2dede;
  color: #a94442;
}
</style>