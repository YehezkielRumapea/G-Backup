<script setup>
import { ref, computed, onMounted } from 'vue';
import apiClient from '@/api'; // 1. Impor apiClient

// --- STATE ---
const remotes = ref([]); // 2. Mulai dengan array kosong
const isLoading = ref(true);
const errorMsg = ref('');
const searchQuery = ref('');

// --- FUNGSI UNTUK MENGAMBIL DATA DARI BACKEND ---
const fetchRemotes = async () => {
    isLoading.value = true;
    errorMsg.value = '';
    try {
        // 3. Panggil API monitoring/remotes
        const response = await apiClient.get('/monitoring/remotes');
        // 4. Isi state 'remotes' dengan data asli dari backend
        remotes.value = response.data || [];
    } catch (error) {
        console.error("Gagal memuat data remote:", error);
        if (error.response?.status !== 401) {
            errorMsg.value = "Gagal memuat data remote dari server.";
        }
    } finally {
        isLoading.value = false;
    }
};

// 5. Panggil API saat komponen pertama kali dimuat
onMounted(fetchRemotes);

// --- FUNGSI-FUNGSI HELPER (Disesuaikan untuk data backend) ---
const filteredRemotes = computed(() => {
  if (!searchQuery.value) {
    return remotes.value;
  }
  return remotes.value.filter(remote =>
    remote.remote_name.toLowerCase().includes(searchQuery.value.toLowerCase())
  );
});

const getStatusClass = (status) => {
  if (status === 'CONNECTED') return 'status-connected';
  if (status === 'DISCONNECTED') return 'status-disconnected';
  return '';
};

const getStoragePercentage = (used, total) => {
  if (total === 0) return 0;
  return (used / total) * 100;
};

const formatTimestamp = (timestamp) => {
    if (!timestamp || timestamp === '0001-01-01T00:00:00Z') return '-';
    try { return new Date(timestamp).toLocaleString('id-ID'); } catch(e) { return timestamp; }
};
</script>

<template>
  <div class="remote-container">
    <header class="main-header">
      <div class="header-content">
        <div>
          <h1>Remote</h1>
          <p>Manage All your Gdrive Remote</p>
        </div>
        <button class="add-remote-btn">
          <span>Add Remote</span>
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M5 12h14"/><path d="M12 5v14"/></svg>
        </button>
      </div>
    </header>

    <div class="card">
      <div class="toolbar">
        <div class="search-wrapper">
           <svg class="search-icon" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line></svg>
           <input type="text" v-model="searchQuery" placeholder="Search Remote Name" class="search-input">
        </div>
      </div>

      <div class="table-responsive">
        <table class="data-table">
          <thead>
            <tr>
              <th>Remote Name</th>
              <th>Status</th>
              <th>Storage</th>
              <th>Last Checked</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="isLoading">
              <td colspan="4" class="text-center">Memuat data remote...</td>
            </tr>
            <tr v-else-if="errorMsg">
              <td colspan="4" class="text-center" style="color: red;">{{ errorMsg }}</td>
            </tr>
            <tr v-else-if="filteredRemotes.length === 0">
              <td colspan="4" class="text-center">Tidak ada remote ditemukan.</td>
            </tr>
            <tr v-for="remote in filteredRemotes" :key="remote.remote_name">
              <td>{{ remote.remote_name }}</td>
              <td>
                <span :class="['status-badge', getStatusClass(remote.status_connect)]">
                  <span class="status-dot"></span>
                  {{ remote.status_connect }}
                </span>
              </td>
              <td>
                <div class="storage-cell">
                  <span>{{ remote.used_storage_gb.toFixed(2) }} of {{ remote.total_storage_gb.toFixed(2) }} Gb ({{ getStoragePercentage(remote.used_storage_gb, remote.total_storage_gb).toFixed(0) }}%) Used</span>
                  <div class="progress-bar-bg">
                    <div class="progress-bar-fg" :style="{ width: getStoragePercentage(remote.used_storage_gb, remote.total_storage_gb) + '%' }"></div>
                  </div>
                </div>
              </td>
              <td>{{ formatTimestamp(remote.last_checked_at) }}</td>
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
.add-remote-btn { background-color: var(--primary-color); color: white; border: none; border-radius: 8px; padding: 10px 20px; font-size: 1rem; font-weight: 500; cursor: pointer; display: flex; align-items: center; gap: 8px; transition: background-color 0.3s; }
.add-remote-btn:hover { background-color: var(--primary-hover); }
.add-remote-btn svg { background-color: rgba(255,255,255,0.2); border-radius: 50%; padding: 2px; }
.card { background-color: #fff; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.05); border: 1px solid var(--border-color); padding: 20px; }
.toolbar { margin-bottom: 20px; }
.search-wrapper { position: relative; max-width: 300px; }
.search-input { width: 100%; padding: 10px 15px 10px 40px; border: 1px solid var(--border-color); border-radius: 6px; font-size: 1rem; }
.search-icon { position: absolute; left: 12px; top: 50%; transform: translateY(-50%); color: #6c757d; }
.table-responsive { overflow-x: auto; }
.data-table { width: 100%; border-collapse: collapse; }
.data-table th, .data-table td { padding: 15px; text-align: left; border-bottom: 1px solid var(--border-color); vertical-align: middle; }
.data-table th { background-color: #f8f9fa; font-weight: 600; color: #495057; }
.status-badge { display: inline-flex; align-items: center; gap: 8px; padding: 5px 12px; border-radius: 16px; font-weight: 500; font-size: 0.9rem; }
.status-dot { width: 8px; height: 8px; border-radius: 50%; }
.status-connected { background-color: #e7f5e9; color: #28a745; }
.status-connected .status-dot { background-color: #28a745; }
.status-disconnected { background-color: #fbe9e7; color: #dc3545; }
.status-disconnected .status-dot { background-color: #dc3545; }
.status-warning { background-color: #fff8e1; color: #ffc107; }
.status-warning .status-dot { background-color: #ffc107; }
.storage-cell { display: flex; flex-direction: column; gap: 5px; }
.progress-bar-bg { width: 100%; height: 6px; background-color: #e9ecef; border-radius: 3px; overflow: hidden; }
.progress-bar-fg { height: 100%; background-color: var(--primary-color); border-radius: 3px; }
</style>