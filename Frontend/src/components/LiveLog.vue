<template>
  <div class="live-log-panel">
    <!-- Header dengan Stats -->
    <div class="panel-header">
      <h3>📝 Recent Logs</h3>
      <button 
        @click="fetchLogs" 
        class="refresh-btn"
        :class="{ spinning: isRefreshing }"
        title="Refresh"
      >
        🔄
      </button>
    </div>

    <!-- Stats Bar -->
    <div class="stats-bar">
      <div class="stat success">
        <span class="stat-icon">✅</span>
        <div class="stat-info">
          <span class="stat-value">{{ successCount }}</span>
          <span class="stat-label">Success</span>
        </div>
      </div>
      <div class="stat failed">
        <span class="stat-icon">❌</span>
        <div class="stat-info">
          <span class="stat-value">{{ failedCount }}</span>
          <span class="stat-label">Failed</span>
        </div>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="isLoading" class="loading">
      <div class="spinner"></div>
      <p>Loading...</p>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="error">
      <p>{{ error }}</p>
      <button @click="fetchLogs" class="retry-btn">Retry</button>
    </div>

    <!-- Logs List -->
    <div v-else class="logs-list">
      <div v-if="logs.length === 0" class="empty">
        <span>📋</span>
        <p>No logs yet</p>
      </div>

      <div 
        v-for="log in logs" 
        :key="log.ID"
        class="log-item"
        :class="getStatusClass(log.Status)"
      >
        <!-- Status Icon -->
        <div class="log-status">
          {{ getStatusIcon(log.Status) }}
        </div>

        <!-- Log Info -->
        <div class="log-info">
          <div class="log-name">{{ getJobName(log) }}</div>
          <div class="log-meta">
            <span class="duration">{{ formatDuration(log.DurationSec) }}</span>
            <span class="time">{{ formatTime(log.Timestamp) }}</span>
          </div>
        </div>

        <!-- Detail Button -->
        <button 
          @click="viewDetail(log)" 
          class="detail-btn"
          title="View detail"
        >
          ...
        </button>
      </div>
    </div>

    <!-- Detail Modal -->
    <SimpleLogModal
      :is-visible="showModal"
      :log="selectedLog"
      @close="closeModal"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue';
import monitoringService from '@/services/monitoringService';
import SimpleLogModal from './LogModel.vue';

const logs = ref([]);
const isLoading = ref(true);
const isRefreshing = ref(false);
const error = ref(null);
const showModal = ref(false);
const selectedLog = ref(null);

let refreshInterval = null;

// Stats
const successCount = computed(() => {
  return logs.value.filter(l => 
    ['SUCCESS', 'COMPLETED'].includes(l.Status.toUpperCase())
  ).length;
});

const failedCount = computed(() => {
  return logs.value.filter(l => 
    l.Status.toUpperCase().includes('FAIL') || l.Status.toUpperCase() === 'ERROR'
  ).length;
});

// Fetch logs
async function fetchLogs() {
  if (isRefreshing.value) return;
  
  isRefreshing.value = true;
  error.value = null;
  
  try {
    const data = await monitoringService.getLogs();
    
    // Take latest 10 logs
    logs.value = (Array.isArray(data) ? data : [])
      .sort((a, b) => new Date(b.Timestamp) - new Date(a.Timestamp))
      .slice(0, 20);
  } catch (err) {
    console.error('Failed to fetch logs:', err);
    error.value = 'Failed to load';
  } finally {
    isLoading.value = false;
    isRefreshing.value = false;
  }
}

// Auto refresh every 10 seconds
function startAutoRefresh() {
  refreshInterval = setInterval(fetchLogs, 10000);
}

function stopAutoRefresh() {
  if (refreshInterval) {
    clearInterval(refreshInterval);
    refreshInterval = null;
  }
}

// View detail
function viewDetail(log) {
  selectedLog.value = log;
  showModal.value = true;
}

function closeModal() {
  showModal.value = false;
  selectedLog.value = null;
}

// Helper functions
function getJobName(log) {
  if (log.ScheduledJob?.JobName) return log.ScheduledJob.JobName;
  if (log.ConfigSnapshot) {
    try {
      const config = JSON.parse(log.ConfigSnapshot);
      return config.job_name || 'Manual Job';
    } catch (e) {
      return 'Manual Job';
    }
  }
  return 'Unknown';
}

function getStatusClass(status) {
  const s = status.toUpperCase();
  if (['SUCCESS', 'COMPLETED'].includes(s)) return 'success';
  if (s.includes('FAIL') || s === 'ERROR') return 'failed';
  if (s === 'RUNNING') return 'running';
  return 'pending';
}

function getStatusIcon(status) {
  const s = status.toUpperCase();
  if (['SUCCESS', 'COMPLETED'].includes(s)) return '✅';
  if (s.includes('FAIL') || s === 'ERROR') return '❌';
  if (s === 'RUNNING') return '⏳';
  return '⏱️';
}

function formatDuration(seconds) {
  if (!seconds) return '-';
  if (seconds < 60) return `${seconds}s`;
  const m = Math.floor(seconds / 60);
  const s = seconds % 60;
  return `${m}m ${s}s`;
}

function formatTime(timestamp) {
  if (!timestamp) return '-';
  
  try {
    const date = new Date(timestamp);
    const now = new Date();
    const diff = now - date;
    
    if (diff < 60000) return 'Just now';
    if (diff < 3600000) return `${Math.floor(diff / 60000)}m ago`;
    if (diff < 86400000) return `${Math.floor(diff / 3600000)}h ago`;
    
    return date.toLocaleString('id-ID', { 
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  } catch (e) {
    return timestamp;
  }
}

// Lifecycle
onMounted(() => {
  fetchLogs();
  startAutoRefresh();
});

onUnmounted(() => {
  stopAutoRefresh();
});
</script>

<style scoped>
.live-log-panel {
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
}

/* Header */
.panel-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 1rem 1.25rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.panel-header h3 {
  margin: 0;
  font-size: 1rem;
  font-weight: 600;
}

.refresh-btn {
 background: rgba(255, 255, 255, 0.2);
 border: none;
 color: white;
 width: 32px;
 height: 32px;
 border-radius: 6px;
 cursor: pointer;
 transition: all 0.2s;
 font-size: 2rem;
 display: flex;
 align-items: center;
 justify-content: center;
}

.refresh-btn:hover {
  background: rgba(255, 255, 255, 0.3);
}

.refresh-btn.spinning {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* Stats Bar */
.stats-bar {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
  padding: 1rem 1.25rem;
  background: #f8f9fa;
  border-bottom: 1px solid #e9ecef;
}

.stat {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem;
  background: white;
  border-radius: 8px;
  border-left: 3px solid;
}

.stat.success {
  border-left-color: #28a745;
}

.stat.failed {
  border-left-color: #dc3545;
}

.stat-icon {
  font-size: 1.5rem;
  line-height: 1;
}

.stat-info {
  display: flex;
  flex-direction: column;
}

.stat-value {
  font-size: 1.5rem;
  font-weight: 700;
  color: #2c3e50;
  line-height: 1;
}

.stat-label {
  font-size: 0.75rem;
  color: #6c757d;
  font-weight: 500;
}

/* Loading, Error, Empty */
.loading,
.error,
.empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 2rem 1rem;
  gap: 0.75rem;
}

.spinner {
  width: 32px;
  height: 32px;
  border: 3px solid #f0f0f0;
  border-top: 3px solid #667eea;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

.loading p,
.error p,
.empty p {
  margin: 0;
  color: #6c757d;
  font-size: 0.9rem;
}

.retry-btn {
  background: #dc3545;
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.85rem;
  transition: all 0.2s;
}

.retry-btn:hover {
  background: #c82333;
}

.empty span {
  font-size: 2.5rem;
  opacity: 0.3;
}

/* Logs List */
.logs-list {
  flex: 1;
  overflow-y: auto;
  padding: 0.75rem;
}

.logs-list::-webkit-scrollbar {
  width: 6px;
}

.logs-list::-webkit-scrollbar-track {
  background: #f1f1f1;
}

.logs-list::-webkit-scrollbar-thumb {
  background: #cbd5e0;
  border-radius: 3px;
}

.logs-list::-webkit-scrollbar-thumb:hover {
  background: #a0aec0;
}

/* Log Item */
.log-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem;
  background: #f8f9fa;
  border-radius: 8px;
  margin-bottom: 0.5rem;
  border-left: 3px solid #cbd5e0;
  transition: all 0.2s;
}

.log-item:hover {
  background: #e9ecef;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.log-item:last-child {
  margin-bottom: 0;
}

.log-item.success {
  border-left-color: #28a745;
}

.log-item.failed {
  border-left-color: #dc3545;
}

.log-item.running {
  border-left-color: #ffc107;
}

.log-item.pending {
  border-left-color: #6c757d;
}

/* Log Status Icon */
.log-status {
  font-size: 1.25rem;
  line-height: 1;
  flex-shrink: 0;
}

/* Log Info */
.log-info {
  flex: 1;
  min-width: 0;
}

.log-name {
  font-weight: 600;
  font-size: 0.85rem;
  color: #2c3e50;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-bottom: 0.25rem;
}

.log-meta {
  display: flex;
  gap: 0.75rem;
  font-size: 0.75rem;
  color: #6c757d;
}

.duration,
.time {
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

/* Detail Button */
/* Detail Button */
.detail-btn {
background: transparent; 
border: none;
color: #6c757d; 
width: 32px;
height: 32px;
border-radius: 6px;
cursor: pointer;
font-size: 2rem;
transition: all 0.2s;
flex-shrink: 0;
display: flex;
align-items: center;
justify-content: center;
padding: 0;
}

.detail-btn:hover {
/* background: #5568d3; <--- HAPUS/UBAH ini */
background: #dcdfe3; /* Ganti dengan warna latar belakang hover yang lembut */
transform: scale(1.0); /* Scale(1.1) terlalu agresif, ganti atau hapus */
}
</style>