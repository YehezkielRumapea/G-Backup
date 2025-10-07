<script setup>
import { ref, computed } from 'vue';

// --- DATA SIMULASI ---
const remotes = ref([
  { id: 1, name: 'Data Duit', status: 'Connected', storageUsed: 7.17, storageTotal: 15, jobRuns: 4, lastChecked: '07:19 PM' },
  { id: 2, name: 'Azure-West-Archive', status: 'Disconnected', storageUsed: 102.5, storageTotal: 500, jobRuns: 12, lastChecked: 'Yesterday' },
  { id: 3, name: 'S3-Primary-Bucket', status: 'Connected', storageUsed: 230.8, storageTotal: 1024, jobRuns: 25, lastChecked: '08:00 AM' },
  { id: 4, name: 'Local-NAS-Media', status: 'Warning', storageUsed: 1800, storageTotal: 2048, jobRuns: 8, lastChecked: '06:00 AM' },
]);

const searchQuery = ref('');

// --- FUNGSI-FUNGSI ---
const filteredRemotes = computed(() => {
  if (!searchQuery.value) {
    return remotes.value;
  }
  return remotes.value.filter(remote =>
    remote.name.toLowerCase().includes(searchQuery.value.toLowerCase())
  );
});

const getStatusClass = (status) => {
  if (status === 'Connected') return 'status-connected';
  if (status === 'Disconnected') return 'status-disconnected';
  if (status === 'Warning') return 'status-warning';
  return '';
};

const getStoragePercentage = (used, total) => {
  if (total === 0) return 0;
  return (used / total) * 100;
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
              <th>Job Runs</th>
              <th>Last Checked</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="filteredRemotes.length === 0">
              <td colspan="5" class="text-center">No remotes found.</td>
            </tr>
            <tr v-for="remote in filteredRemotes" :key="remote.id">
              <td>{{ remote.name }}</td>
              <td>
                <span :class="['status-badge', getStatusClass(remote.status)]">
                  <span class="status-dot"></span>
                  {{ remote.status }}
                </span>
              </td>
              <td>
                <div class="storage-cell">
                  <span>{{ remote.storageUsed }} of {{ remote.storageTotal }} Gb ({{ getStoragePercentage(remote.storageUsed, remote.storageTotal).toFixed(0) }}%) Used</span>
                  <div class="progress-bar-bg">
                    <div class="progress-bar-fg" :style="{ width: getStoragePercentage(remote.storageUsed, remote.storageTotal) + '%' }"></div>
                  </div>
                </div>
              </td>
              <td>{{ remote.jobRuns }}</td>
              <td>{{ remote.lastChecked }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<style scoped>
.main-header {
  margin-bottom: 20px;
}
.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.main-header h1 { font-size: 2rem; font-weight: 700; color: #333; }
.main-header p { color: #6c757d; }

.add-remote-btn {
  background-color: var(--primary-color);
  color: white;
  border: none;
  border-radius: 8px;
  padding: 10px 20px;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
  transition: background-color 0.3s;
}
.add-remote-btn:hover {
  background-color: var(--primary-hover);
}
.add-remote-btn svg {
  background-color: rgba(255,255,255,0.2);
  border-radius: 50%;
  padding: 2px;
}

.card {
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.05);
  border: 1px solid var(--border-color);
  padding: 20px;
}

.toolbar { margin-bottom: 20px; }
.search-wrapper {
  position: relative;
  max-width: 300px;
}
.search-input {
  width: 100%;
  padding: 10px 15px 10px 40px; /* Ruang untuk ikon */
  border: 1px solid var(--border-color);
  border-radius: 6px;
  font-size: 1rem;
}
.search-icon {
  position: absolute;
  left: 12px;
  top: 50%;
  transform: translateY(-50%);
  color: #6c757d;
}

.table-responsive { overflow-x: auto; }
.data-table {
  width: 100%;
  border-collapse: collapse;
}
.data-table th, .data-table td {
  padding: 15px;
  text-align: left;
  border-bottom: 1px solid var(--border-color);
  vertical-align: middle;
}
.data-table th {
  background-color: #f8f9fa;
  font-weight: 600;
  color: #495057;
}

.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 5px 12px;
  border-radius: 16px;
  font-weight: 500;
  font-size: 0.9rem;
}
.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.status-connected { background-color: #e7f5e9; color: #28a745; }
.status-connected .status-dot { background-color: #28a745; }

.status-disconnected { background-color: #fbe9e7; color: #dc3545; }
.status-disconnected .status-dot { background-color: #dc3545; }

.status-warning { background-color: #fff8e1; color: #ffc107; }
.status-warning .status-dot { background-color: #ffc107; }

.storage-cell { display: flex; flex-direction: column; gap: 5px; }
.progress-bar-bg {
  width: 100%;
  height: 6px;
  background-color: #e9ecef;
  border-radius: 3px;
  overflow: hidden;
}
.progress-bar-fg {
  height: 100%;
  background-color: var(--primary-color);
  border-radius: 3px;
}
</style>

