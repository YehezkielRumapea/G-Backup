<script setup>
import { ref, computed } from 'vue';

// --- DATA SIMULASI ---
const allJobs = ref([
  { id: 1, name: 'Backup_Daily_DB_01', remote: 'Azure-West', status: 'Completed', schedule: 'Daily at 23:00', lastRun: '2025-10-06 23:00', duration: '0h 12m', nextRun: '2025-10-07 23:00' },
  { id: 2, name: 'Archive_Docs_Monthly', remote: 'S3-Archive', status: 'Running', schedule: '1st of Month', lastRun: '2025-10-07 01:30', duration: '0h 45m', nextRun: '2025-11-01 01:00' },
  { id: 3, name: 'Sync_Media_Weekly_Photos', remote: 'Local-NAS', status: 'Failed', schedule: 'Weekly (Sun)', lastRun: '2025-10-07 02:05', duration: '0h 03m', nextRun: '2025-10-12 02:00' },
  { id: 4, name: 'VM_Snapshot_Hourly_AppSrv', remote: 'GCS-East', status: 'Pending', schedule: 'Every 1 Hour', lastRun: '-', duration: '-', nextRun: '2025-10-07 10:00' },
  { id: 5, name: 'Backup_WebApp_Source', remote: 'Google Drive Utama', status: 'Completed', schedule: 'Daily at 03:00', lastRun: '2025-10-07 03:00', duration: '0h 08m', nextRun: '2025-10-08 03:00' },
  { id: 6, name: 'Backup_SQL_SalesDB', remote: 'Azure-West', status: 'Pending', schedule: 'Every 6 Hours', lastRun: '-', duration: '-', nextRun: '2025-10-07 12:00' },
]);

// --- STATE UNTUK FILTER DAN PENCARIAN ---
const searchQuery = ref('');
const statusFilter = ref('All');
const remoteFilter = ref('All');

// --- COMPUTED PROPERTY UNTUK MENAMPILKAN DATA SESUAI FILTER ---
const filteredJobs = computed(() => {
  return allJobs.value.filter(job => {
    const searchMatch = job.name.toLowerCase().includes(searchQuery.value.toLowerCase());
    const statusMatch = statusFilter.value === 'All' || job.status === statusFilter.value;
    const remoteMatch = remoteFilter.value === 'All' || job.remote === remoteFilter.value;
    return searchMatch && statusMatch && remoteMatch;
  });
});

// --- HELPER UNTUK STYLING ---
const getStatusClass = (status) => {
    switch (status) {
        case 'Completed': return 'status-completed';
        case 'Running': return 'status-running';
        case 'Failed': return 'status-failed';
        case 'Pending': return 'status-pending';
        default: return '';
    }
};

// --- FUNGSI AKSI (PLACEHOLDER) ---
const runJob = (id) => alert(`Running job with ID: ${id}`);
const editJob = (id) => alert(`Editing job with ID: ${id}`);
const deleteJob = (id) => {
    if (confirm('Are you sure you want to delete this job?')) {
        allJobs.value = allJobs.value.filter(job => job.id !== id);
    }
};

</script>

<template>
  <div class="job-status-content">
    <header class="main-header">
      <h1>Jobs</h1>
      <p>Monitor and manage all backup and synchronization jobs.</p>
    </header>

    <div class="card">
      <!-- Toolbar untuk filter dan aksi -->
      <div class="toolbar">
        <div class="filters">
          <input type="text" v-model="searchQuery" placeholder="Search job name..." class="search-input">
          <select v-model="statusFilter" class="filter-select">
            <option>All</option>
            <option>Completed</option>
            <option>Running</option>
            <option>Failed</option>
            <option>Pending</option>
          </select>
          <select v-model="remoteFilter" class="filter-select">
            <option>All</option>
            <option>Azure-West</option>
            <option>S3-Archive</option>
            <option>Local-NAS</option>
            <option>GCS-East</option>
             <option>Google Drive Utama</option>
          </select>
        </div>
        <div class="batch-actions">
          <button class="btn btn-secondary">Run Selected</button>
          <button class="btn btn-secondary">Stop Selected</button>
          <button class="btn btn-danger">Delete Selected</button>
        </div>
      </div>

      <!-- Tabel Job Status -->
      <div class="table-responsive">
        <table class="data-table">
          <thead>
            <tr>
              <th><input type="checkbox"></th>
              <th>Job Name</th>
              <th>Remote</th>
              <th>Status</th>
              <th>Schedule</th>
              <th>Last Run</th>
              <th>Duration</th>
              <th>Next Run</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="filteredJobs.length === 0">
                <td colspan="9" class="text-center">No jobs match the current filters.</td>
            </tr>
            <tr v-for="job in filteredJobs" :key="job.id">
              <td><input type="checkbox"></td>
              <td>{{ job.name }}</td>
              <td>{{ job.remote }}</td>
              <td>
                <span :class="['status-badge', getStatusClass(job.status)]">{{ job.status }}</span>
              </td>
              <td>{{ job.schedule }}</td>
              <td>{{ job.lastRun }}</td>
              <td>{{ job.duration }}</td>
              <td>{{ job.nextRun }}</td>
              <td class="action-buttons">
                <button @click="runJob(job.id)" class="btn-icon btn-run" title="Run Now">
                  <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16"><path d="M10.804 8 5 4.633v6.734L10.804 8zM5 3.37a.5.5 0 0 1 .79-.407l7 5a.5.5 0 0 1 0 .814l-7 5A.5.5 0 0 1 5 13.63V3.37z"/></svg>
                </button>
                <button @click="editJob(job.id)" class="btn-icon btn-edit" title="Edit">
                  <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16"><path d="M12.854.146a.5.5 0 0 0-.707 0L10.5 1.793 14.207 5.5l1.647-1.646a.5.5 0 0 0 0-.708l-3-3zm.646 6.061L9.793 2.5 3.293 9H3.5a.5.5 0 0 1 .5.5v.5h.5a.5.5 0 0 1 .5.5v.5h.5a.5.5 0 0 1 .5.5v.5h.5a.5.5 0 0 1 .5.5v.207l6.5-6.5zm-7.468 7.468A.5.5 0 0 1 6 13.5V13h-.5a.5.5 0 0 1-.5-.5V12h-.5a.5.5 0 0 1-.5-.5V11h-.5a.5.5 0 0 1-.5-.5V10h-.5a.499.499 0 0 1-.175-.032l-.179.178a.5.5 0 0 0-.11.168l-2 5a.5.5 0 0 0 .65.65l5-2a.5.5 0 0 0 .168-.11l.178-.178z"/></svg>
                </button>
                <button @click="deleteJob(job.id)" class="btn-icon btn-delete" title="Delete">
                  <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16"><path d="M5.5 5.5A.5.5 0 0 1 6 6v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5zm2.5 0a.5.5 0 0 1 .5.5v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5zm3 .5a.5.5 0 0 0-1 0v6a.5.5 0 0 0 1 0V6z"/><path fill-rule="evenodd" d="M14.5 3a1 1 0 0 1-1 1H13v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V4h-.5a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1H6a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1h3.5a1 1 0 0 1 1 1v1zM4.118 4 4 4.059V13a1 1 0 0 0 1 1h6a1 1 0 0 0 1-1V4.059L11.882 4H4.118zM2.5 3V2h11v1h-11z"/></svg>
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<style scoped>
.main-header {
  padding-bottom: 10px;
  border-bottom: 1px solid var(--border-color);
  margin-bottom: 20px;
}
.main-header h1 {
  font-size: 1.8rem;
}
.main-header p {
  color: #6c757d;
}

.card {
  background-color: var(--card-bg);
  border-radius: 8px;
  box-shadow: 0 4px 8px rgba(0,0,0,0.05);
  border: 1px solid var(--border-color);
  padding: 20px;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  flex-wrap: wrap;
  gap: 15px;
}

.filters {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.search-input, .filter-select {
  padding: 8px 12px;
  border-radius: 6px;
  border: 1px solid var(--border-color);
  font-size: 0.9rem;
}

.search-input {
  width: 250px;
}

.batch-actions {
  display: flex;
  gap: 10px;
}

.btn {
  padding: 8px 16px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.9rem;
  font-weight: 500;
  transition: all 0.2s;
}
.btn-secondary { background-color: #6c757d; color: white; }
.btn-secondary:hover { background-color: #5a6268; }
.btn-danger { background-color: #dc3545; color: white; }
.btn-danger:hover { background-color: #c82333; }

.table-responsive { overflow-x: auto; }

.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th, .data-table td {
  padding: 12px 15px;
  text-align: left;
  border-bottom: 1px solid var(--border-color);
  vertical-align: middle;
  white-space: nowrap;
}

.data-table th {
  background-color: #f8f9fa;
  font-weight: 600;
  color: #495057;
}

.data-table tbody tr:hover { background-color: #f1f3f5; }

.text-center { text-align: center; color: #6c757d; padding: 20px; }

.status-badge {
  display: inline-block;
  padding: 5px 10px;
  border-radius: 15px;
  font-weight: 500;
  font-size: 0.8rem;
  color: white;
  min-width: 80px;
  text-align: center;
}
.status-completed { background-color: #28a745; }
.status-running { background-color: #ffc107; color: #333; }
.status-failed { background-color: #dc3545; }
.status-pending { background-color: #6c757d; }

.action-buttons { display: flex; gap: 10px; }
.btn-icon { background: none; border: none; cursor: pointer; padding: 5px; border-radius: 50%; display: flex; transition: background-color 0.2s; }
.btn-icon svg { width: 16px; height: 16px; }
.btn-run { color: #28a745; }
.btn-run:hover { background-color: rgba(40, 167, 69, 0.1); }
.btn-edit { color: #007bff; }
.btn-edit:hover { background-color: rgba(0, 123, 255, 0.1); }
.btn-delete { color: #dc3545; }
.btn-delete:hover { background-color: rgba(220, 53, 69, 0.1); }
</style>
