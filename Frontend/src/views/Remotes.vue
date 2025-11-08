<template>
  <div class="remotes-view">
    <h1>Monitoring Remote</h1>
    <p>Kelola dan pantau status GDrive Remote Anda.</p>
    
    <router-link to="/create" class="btn-add-remote">
      + Tambah Remote Baru
    </router-link>

    <div v-if="isLoading" class="loading">
      Memuat data remote...
    </div>
    
    <div v-if="errorMessage" class="error">
      {{ errorMessage }}
    </div>

    <table v-if="remotes.length > 0" class="remotes-table">
      <thead>
        <tr>
          <th>Nama Remote</th>
          <th>Status</th>
          <th>Storage</th>
          <th>Job Runs</th>
          <th>Last Checked</th>
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
</template>

<script setup>
import { ref, onMounted } from 'vue'
import monitoringService from '@/services/monitoringService' // Import service API
import RemoteRow from '@/components/RemoteRow.vue' // Import komponen baris

// State lokal
const remotes = ref([])
const isLoading = ref(true)
const errorMessage = ref(null)

// onMounted dipanggil saat komponen (halaman) dimuat
onMounted(async () => {
  try {
    // Panggil API backend Golang untuk mengambil data remote
    const data = await monitoringService.getRemoteStatus()
    remotes.value = data
  } catch (error) {
    // Tangani jika API gagal (misalnya token salah)
    errorMessage.value = 'Gagal memuat data remote. Coba login ulang.'
  } finally {
    isLoading.value = false
  }
})
</script>

<style scoped>
.remotes-table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 1.5rem;
}
.remotes-table th, .remotes-table td {
  border-bottom: 1px solid #ddd;
  padding: 12px 15px;
  text-align: left;
}
.remotes-table th {
  background-color: #f4f4f4;
}
.loading, .error {
  margin-top: 1rem;
  font-style: italic;
}
.error {
  color: red;
}
.btn-add-remote {
  display: inline-block;
  background-color: #1abc9c;
  color: white;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  text-decoration: none;
  font-weight: bold;
}
</style>