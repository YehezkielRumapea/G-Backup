<script setup>
import { ref, computed } from 'vue';

// --- DATA LOG SIMULASI ---
const allLogs = ref([
  { id: 1, timestamp: '2025-10-07 09:10:05', level: 'INFO', job: 'Local-NAS-01', message: 'Connection test successful.' },
  { id: 2, timestamp: '2025-10-07 08:15:22', level: 'INFO', job: 'Google Drive Utama', message: 'Authentication successful.' },
  { id: 3, timestamp: '2025-10-07 03:08:45', level: 'SUCCESS', job: 'Backup_WebApp_Source', message: 'Job completed successfully. 1.2GB transferred.' },
  { id: 4, timestamp: '2025-10-07 03:00:10', level: 'INFO', job: 'Backup_WebApp_Source', message: 'Job started.' },
  { id: 5, timestamp: '2025-10-07 02:08:15', level: 'ERROR', job: 'Sync_Media_Weekly_Photos', message: 'Failed to connect to remote: Connection timed out.' },
  { id: 6, timestamp: '2025-10-07 02:05:01', level: 'INFO', job: 'Sync_Media_Weekly_Photos', message: 'Job started.' },
  { id: 7, timestamp: '2025-10-07 01:30:00', level: 'INFO', job: 'Archive_Docs_Monthly', message: 'Job started, scanning 15,000 files...' },
]);

// --- STATE UNTUK FILTER ---
const searchQuery = ref('');
const levelFilter = ref('All');

// --- COMPUTED PROPERTY UNTUK FILTER LOGS ---
const filteredLogs = computed(() => {
  return allLogs.value.filter(log => {
    const searchMatch = log.message.toLowerCase().includes(searchQuery.value.toLowerCase()) || log.job.toLowerCase().includes(searchQuery.value.toLowerCase());
    const levelMatch = levelFilter.value === 'All' || log.level === levelFilter.value;
    return searchMatch && levelMatch;
  });
});

// --- HELPER UNTUK STYLING LEVEL LOG ---
const getLevelClass = (level) => {
    switch (level) {
        case 'SUCCESS': return 'level-success';
        case 'INFO': return 'level-info';
        case 'WARNING': return 'level-warning';
        case 'ERROR': return 'level-error';
        default: return '';
    }
}
</script>

<template>
  <div class="logs-content">
    <header class="main-header">
      <h1>Log Activity</h1>
      <p>Review system, job, and error logs.</p>
    </header>

    <div class="card">
      <div class="toolbar">
        <input type="text" v-model="searchQuery" placeholder="Search in logs..." class="search-input">
        <select v-model="levelFilter" class="filter-select">
          <option>All</option>
          <option>SUCCESS</option>
          <option>INFO</option>
          <option>WARNING</option>
          <option>ERROR</option>
        </select>
      </div>

      <div class="table-responsive">
        <table class="data-table">
          <thead>
            <tr>
              <th>Timestamp</th>
              <th>Level</th>
              <th>Job/System</th>
              <th>Message</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="filteredLogs.length === 0">
              <td colspan="4" class="text-center">No logs found.</td>
            </tr>
            <tr v-for="log in filteredLogs" :key="log.id">
              <td>{{ log.timestamp }}</td>
              <td><span :class="['level-badge', getLevelClass(log.level)]">{{ log.level }}</span></td>
              <td>{{ log.job }}</td>
              <td class="log-message">{{ log.message }}</td>
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
.main-header h1 { font-size: 1.8rem; }
.main-header p { color: #6c757d; }

.card {
  background-color: var(--card-bg);
  border-radius: 8px;
  box-shadow: 0 4px 8px rgba(0,0,0,0.05);
  border: 1px solid var(--border-color);
  padding: 20px;
}

.toolbar {
  display: flex;
  gap: 15px;
  margin-bottom: 20px;
}

.search-input, .filter-select {
  padding: 8px 12px;
  border-radius: 6px;
  border: 1px solid var(--border-color);
  font-size: 0.9rem;
}
.search-input { flex-grow: 1; }

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
}

.data-table th {
  background-color: #f8f9fa;
  font-weight: 600;
}

.data-table tbody tr:hover { background-color: #f1f3f5; }
.text-center { text-align: center; color: #6c757d; padding: 20px; }

.level-badge {
  display: inline-block;
  padding: 4px 9px;
  border-radius: 6px;
  font-weight: 600;
  font-size: 0.75rem;
  color: white;
  text-align: center;
}

.level-success { background-color: #28a745; }
.level-info { background-color: #17a2b8; }
.level-warning { background-color: #ffc107; color: #333; }
.level-error { background-color: #dc3545; }

.log-message {
  white-space: normal;
  word-break: break-word;
}
</style>
