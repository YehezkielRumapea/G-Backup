<template>
  <div class="dashboard-container">
    <!-- Page Header -->
    <div class="page-header">
      <h1>📊 Dashboard</h1>
      <p>Overview sistem backup dan monitoring activity</p>
    </div>

    <!-- Stats Cards -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon">☁️</div>
        <div class="stat-content">
          <span class="stat-label">Total Remote</span>
          <span class="stat-value">{{ totalRemotes }}</span>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">⚙️</div>
        <div class="stat-content">
          <span class="stat-label">Total Job</span>
          <span class="stat-value">{{ totalJobs }}</span>
        </div>
      </div>
    </div>

    <!-- Quick Action Section -->
    <div class="quick-action-section">
      <h2>⚡ Quick Action</h2>
      <p>Start backup or restore job instantly</p>
      
      <div class="action-cards">
        <div class="action-card backup" @click="openBackupModal">
          <div class="action-icon">📦</div>
          <h3>Create Backup</h3>
          <p>Start manual backup job</p>
          <div class="action-arrow">→</div>
        </div>

        <div class="action-card restore" @click="openRestoreModal">
          <div class="action-icon">🔄</div>
          <h3>Start Restore</h3>
          <p>Restore from cloud backup</p>
          <div class="action-arrow">→</div>
        </div>
      </div>
    </div>

    <!-- Next Job Section -->
    <div class="next-job-section">
      <h2>⏰ Next Scheduled Job</h2>
      
      <div v-if="isLoadingJobs" class="loading-box">
        <div class="spinner"></div>
        <p>Loading jobs...</p>
      </div>

      <div v-else-if="nextJob" class="next-job-card">
        <div class="job-header">
          <span class="job-type-badge">{{ nextJob.jobType }}</span>
          <span class="job-time">{{ formatNextRun(nextJob.nextRun) }}</span>
        </div>
        <h3>{{ nextJob.name }}</h3>
        <div class="job-details">
          <div class="detail-item">
            <span class="detail-label">Remote:</span>
            <span class="detail-value">{{ nextJob.remoteName }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">Schedule:</span>
            <span class="detail-value">{{ nextJob.scheduleCron }}</span>
          </div>
        </div>
      </div>

      <div v-else class="empty-box">
        <span class="empty-icon">📅</span>
        <p>No scheduled jobs</p>
      </div>
    </div>

    <!-- Backup Modal -->
    <QuickBackupModal
      :is-visible="showBackupModal"
      @close="showBackupModal = false"
      @success="handleBackupSuccess"
    />

    <!-- Restore Modal -->
    <QuickRestoreModal
      :is-visible="showRestoreModal"
      @close="showRestoreModal = false"
      @success="handleRestoreSuccess"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import monitoringService from '@/services/monitoringService';
import QuickBackupModal from '@/components/CreateBackup.vue';
import QuickRestoreModal from '@/components/CreateRestore.vue';

// State
const remotes = ref([]);
const jobs = ref([]);
const isLoadingJobs = ref(true);
const showBackupModal = ref(false);
const showRestoreModal = ref(false);

// Computed
const totalRemotes = computed(() => remotes.value.length);
const totalJobs = computed(() => jobs.value.length);

const nextJob = computed(() => {
  if (!Array.isArray(jobs.value) || jobs.value.length === 0) return null;
  
  const futureJobs = jobs.value
    .map(j => ({ ...j, nextRunDate: new Date(j.nextRun) }))
    .filter(j => j.nextRunDate && j.nextRunDate > new Date())
    .sort((a, b) => a.nextRunDate - b.nextRunDate);
  
  return futureJobs.length > 0 ? futureJobs[0] : null;
});

// Lifecycle
onMounted(async () => {
  await fetchData();
});

// Functions
async function fetchData() {
  isLoadingJobs.value = true;
  
  try {
    const [remotesRes, jobsRes] = await Promise.all([
      monitoringService.getRemotes(),
      monitoringService.getScheduledJobs()
    ]);
    
    remotes.value = Array.isArray(remotesRes.data) ? remotesRes.data : [];
    jobs.value = Array.isArray(jobsRes.data) ? jobsRes.data : [];
  } catch (error) {
    console.error('Failed to fetch dashboard data:', error);
  } finally {
    isLoadingJobs.value = false;
  }
}

function openBackupModal() {
  showBackupModal.value = true;
}

function openRestoreModal() {
  showRestoreModal.value = true;
}

function handleBackupSuccess() {
  showBackupModal.value = false;
  fetchData(); // Refresh data
}

function handleRestoreSuccess() {
  showRestoreModal.value = false;
}

function formatNextRun(timestamp) {
  if (!timestamp) return '-';
  
  try {
    const date = new Date(timestamp);
    const now = new Date();
    const diff = date - now;
    
    if (diff < 0) return 'Overdue';
    
    const hours = Math.floor(diff / 3600000);
    const minutes = Math.floor((diff % 3600000) / 60000);
    
    if (hours < 1) return `in ${minutes}m`;
    if (hours < 24) return `in ${hours}h ${minutes}m`;
    
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
</script>

<style scoped>
.dashboard-container {
  max-width: 100%;
}

/* Page Header */
.page-header {
  margin-bottom: 2rem;
}

.page-header h1 {
  font-size: 2rem;
  font-weight: 700;
  color: #2c3e50;
  margin: 0 0 0.5rem 0;
}

.page-header p {
  color: #6c757d;
  margin: 0;
  font-size: 1rem;
}

/* Stats Grid */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.stat-card {
  background: white;
  padding: 1.5rem;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  display: flex;
  align-items: center;
  gap: 1rem;
  transition: all 0.3s;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.12);
}

.stat-icon {
  font-size: 2.5rem;
  line-height: 1;
}

.stat-content {
  display: flex;
  flex-direction: column;
}

.stat-label {
  font-size: 0.85rem;
  color: #6c757d;
  font-weight: 500;
  margin-bottom: 0.25rem;
}

.stat-value {
  font-size: 2rem;
  font-weight: 700;
  color: #2c3e50;
  line-height: 1;
}

/* Quick Action Section */
.quick-action-section {
  background: white;
  padding: 2rem;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  margin-bottom: 2rem;
}

.quick-action-section h2 {
  font-size: 1.5rem;
  font-weight: 700;
  color: #2c3e50;
  margin: 0 0 0.5rem 0;
}

.quick-action-section > p {
  color: #6c757d;
  margin: 0 0 1.5rem 0;
}

.action-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
}

.action-card {
  position: relative;
  padding: 2rem;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s;
  overflow: hidden;
}

.action-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  opacity: 0.1;
  transition: opacity 0.3s;
}

.action-card.backup {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.action-card.restore {
  background: linear-gradient(135deg, #3498db 0%, #2980b9 100%);
  color: white;
}

.action-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.2);
}

.action-card:hover::before {
  opacity: 0.2;
}

.action-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
  line-height: 1;
}

.action-card h3 {
  font-size: 1.25rem;
  font-weight: 700;
  margin: 0 0 0.5rem 0;
}

.action-card p {
  margin: 0;
  opacity: 0.9;
  font-size: 0.95rem;
}

.action-arrow {
  position: absolute;
  bottom: 1.5rem;
  right: 1.5rem;
  font-size: 1.5rem;
  transition: transform 0.3s;
}

.action-card:hover .action-arrow {
  transform: translateX(4px);
}

/* Next Job Section */
.next-job-section {
  background: white;
  padding: 2rem;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.next-job-section h2 {
  font-size: 1.5rem;
  font-weight: 700;
  color: #2c3e50;
  margin: 0 0 1.5rem 0;
}

.next-job-card {
  background: #f8f9fa;
  padding: 1.5rem;
  border-radius: 8px;
  border-left: 4px solid #667eea;
}

.job-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.job-type-badge {
  background: #667eea;
  color: white;
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 0.8rem;
  font-weight: 600;
  text-transform: uppercase;
}

.job-time {
  color: #6c757d;
  font-weight: 600;
  font-size: 0.9rem;
}

.next-job-card h3 {
  font-size: 1.25rem;
  font-weight: 700;
  color: #2c3e50;
  margin: 0 0 1rem 0;
}

.job-details {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
}

.detail-item {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.detail-label {
  font-size: 0.8rem;
  color: #6c757d;
  font-weight: 500;
}

.detail-value {
  font-size: 0.95rem;
  color: #2c3e50;
  font-weight: 600;
  font-family: monospace;
}

/* Loading & Empty States */
.loading-box,
.empty-box {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 3rem 2rem;
  gap: 1rem;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #f3f3f3;
  border-top: 4px solid #667eea;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.empty-icon {
  font-size: 3rem;
  opacity: 0.3;
}

.loading-box p,
.empty-box p {
  margin: 0;
  color: #6c757d;
}

/* Responsive */
@media (max-width: 768px) {
  .dashboard-container {
    padding: 1rem;
  }
  
  .page-header h1 {
    font-size: 1.5rem;
  }
  
  .stats-grid {
    grid-template-columns: 1fr;
  }
  
  .action-cards {
    grid-template-columns: 1fr;
  }
  
  .job-details {
    grid-template-columns: 1fr;
  }
}
</style>