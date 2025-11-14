<template>
  <div class="dashboard-container">
    <!-- Page Header -->
    <div class="header">
      <h1>Dashboard</h1>
      <p class="subtitle">Overview sistem backup dan monitoring activity</p>
    </div>

    <!-- Stats Cards -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-content">
          <span class="stat-label">Total Remote</span>
          <span class="stat-value">{{ totalRemotes }}</span>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-content">
          <span class="stat-label">Total Job</span>
          <span class="stat-value">{{ totalJobs }}</span>
        </div>
      </div>
    </div>

    <!-- Quick Action Section -->
    <div class="section">
      <h2>Quick Action</h2>
      <p class="section-subtitle">Start backup or restore job instantly</p>
      
      <div class="action-cards">
        <button class="action-card" @click="openBackupModal">
          <div class="action-content">
            <h3>Create Backup</h3>
            <p>Start manual backup job</p>
          </div>
          <span class="action-arrow">→</span>
        </button>

        <button class="action-card" @click="openRestoreModal">
          <div class="action-content">
            <h3>Start Restore</h3>
            <p>Restore from cloud backup</p>
          </div>
          <span class="action-arrow">→</span>
        </button>
      </div>
    </div>

    <!-- Next Job Section -->
    <div class="section">
      <h2>Next Scheduled Job</h2>
      
      <div v-if="isLoadingJobs" class="status-message">
        <span class="loading-dot"></span>
        Loading jobs...
      </div>

      <div v-else-if="nextJob" class="next-job-card">
        <div class="job-header">
          <span class="job-type-badge">{{ nextJob.jobType }}</span>
          <span class="job-time">{{ formatNextRun(nextJob.nextRun) }}</span>
        </div>
        <h3>{{ nextJob.name }}</h3>
        <div class="job-details">
          <div class="detail-item">
            <span class="detail-label">Remote</span>
            <span class="detail-value">{{ nextJob.remoteName }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">Schedule</span>
            <span class="detail-value">{{ nextJob.scheduleCron }}</span>
          </div>
        </div>
      </div>

      <div v-else class="empty-state">
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
    const [remotesData, jobsData] = await Promise.all([
      monitoringService.getRemoteStatus(),
      monitoringService.getScheduledJobs()
    ]);
    
    remotes.value = Array.isArray(remotesData) ? remotesData : [];
    jobs.value = Array.isArray(jobsData) ? jobsData : [];
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
  fetchData();
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
  max-width: 1400px;
  margin: 0 auto;
  padding: 2rem 1.5rem;
}

/* Page Header */
.header {
  margin-bottom: 2.5rem;
}

h1 {
  font-size: 1.75rem;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0 0 0.5rem 0;
  letter-spacing: -0.02em;
}

.subtitle {
  font-size: 0.95rem;
  color: #666;
  margin: 0;
  font-weight: 400;
}

/* Stats Grid */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
  margin-bottom: 2rem;
}

.stat-card {
  background: white;
  padding: 1.5rem;
  border-radius: 8px;
  border: 1px solid #e5e5e5;
}

.stat-content {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.stat-label {
  font-size: 0.8125rem;
  color: #666;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.stat-value {
  font-size: 2rem;
  font-weight: 600;
  color: #1a1a1a;
  line-height: 1;
}

/* Section */
.section {
  background: white;
  padding: 1.5rem;
  border-radius: 8px;
  border: 1px solid #e5e5e5;
  margin-bottom: 1.5rem;
}

.section h2 {
  font-size: 1.125rem;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0 0 0.25rem 0;
}

.section-subtitle {
  font-size: 0.875rem;
  color: #666;
  margin: 0 0 1.5rem 0;
}

/* Action Cards */
.action-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1rem;
}

.action-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1.25rem;
  border: 1px solid #e5e5e5;
  border-radius: 6px;
  background: white;
  cursor: pointer;
  transition: all 0.2s;
  text-align: left;
  width: 100%;
}

.action-card:hover {
  background: #fafafa;
  border-color: #1a1a1a;
}

.action-content h3 {
  font-size: 1rem;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0 0 0.25rem 0;
}

.action-content p {
  margin: 0;
  font-size: 0.875rem;
  color: #666;
}

.action-arrow {
  font-size: 1.25rem;
  color: #666;
  transition: transform 0.2s;
}

.action-card:hover .action-arrow {
  transform: translateX(4px);
  color: #1a1a1a;
}

/* Next Job Card */
.next-job-card {
  background: #fafafa;
  padding: 1.25rem;
  border-radius: 6px;
  border-left: 3px solid #1a1a1a;
}

.job-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.875rem;
}

.job-type-badge {
  background: #1a1a1a;
  color: white;
  padding: 0.25rem 0.625rem;
  border-radius: 4px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.job-time {
  color: #666;
  font-weight: 500;
  font-size: 0.875rem;
}

.next-job-card h3 {
  font-size: 1.125rem;
  font-weight: 600;
  color: #1a1a1a;
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
  font-size: 0.75rem;
  color: #666;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.detail-value {
  font-size: 0.9375rem;
  color: #1a1a1a;
  font-weight: 500;
  font-family: monospace;
}

/* Status Message */
.status-message {
  padding: 1rem;
  border-radius: 6px;
  font-size: 0.9375rem;
  background: #f8f8f8;
  color: #666;
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.loading-dot {
  width: 8px;
  height: 8px;
  background: #666;
  border-radius: 50%;
  animation: pulse 1.5s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.3; }
}

/* Empty State */
.empty-state {
  text-align: center;
  padding: 2rem 1rem;
}

.empty-state p {
  margin: 0;
  color: #999;
  font-size: 0.9375rem;
}

/* Responsive */
@media (max-width: 768px) {
  .dashboard-container {
    padding: 1.5rem 1rem;
  }
  
  h1 {
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