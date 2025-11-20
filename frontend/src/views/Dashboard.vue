<template>
  <div class="dashboard-container">
    <div class="header">
      <h1>Dashboard</h1>
      <p class="subtitle">Overview sistem backup dan monitoring activity</p>
    </div>

    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-content">
          <span class="stat-label">Total Remote</span>
          <span class="stat-value">{{ totalRemotes }}</span>
        </div>

        <div v-if="!isLoadingData" class="drive-breakdown-list">
          <!-- Rincian Kapasitas per Drive (sesuai permintaan Anda) -->
          <div v-for="remote in remotes" :key="remote.id" class="drive-item">
            <span class="drive-name">☁️ {{ remote.remote_name }}</span>
            <span class="drive-size">
              {{ (remote.used_storage_gb || 0).toFixed(2) }} GB / {{ (remote.total_storage_gb || 0).toFixed(0) }} GB
            </span>
          </div>
          <div v-if="remotes.length === 0" class="drive-empty">
            Belum ada remote
          </div>
        </div>
        <div v-else class="status-message small">
          <span class="loading-dot"></span> Loading...
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-content">
          <span class="stat-label">Total Job</span>
          <span class="stat-value">{{ totalJobs }}</span>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-content chart-card-content">
          <span class="stat-label">Total Storage (Used vs Free)</span>
          
          <div v-if="isLoadingData" class="status-message">
              <span class="loading-dot"></span>
              Loading storage...
          </div>
          
          <!-- Mengirim total kapasitas ke Donut Chart -->
          <StorageDonutChart
            v-else
            :series="storageStats.series"
            :labels="storageStats.labels"
            :totalAvailable="storageStats.totalAvailable"
          />
        </div>
      </div>
    </div>

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

    <div class="section">
      <h2>Next Run</h2>
      
      <div v-if="isLoadingData" class="status-message">
        <span class="loading-dot"></span>
        Loading jobs...
      </div>
      <div v-else-if="nextJob" class="next-job-card">
        <div class="job-header">
          <span class="job-type-badge">{{ nextJob.jobType || 'JOB' }}</span>
          <span class="job-time">{{ formatNextRun(nextJob.nextRun) }}</span>
        </div>
        <h3>{{ nextJob.name }}</h3>
        <div class="job-details">
          <div class="detail-item">
            <span class="detail-label">Drive Target</span>
            <span class="detail-value">{{ nextJob.remoteName }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">Time</span>
            <span class="detail-value">{{ nextJob.scheduleCron }}</span>
          </div>
        </div>
      </div>
      <div v-else class="empty-state">
        <p>No scheduled jobs</p>
      </div>
    </div>

    <QuickBackupModal :is-visible="showBackupModal" @close="showBackupModal = false" @success="handleBackupSuccess" />
    <QuickRestoreModal :is-visible="showRestoreModal" @close="showRestoreModal = false" @success="handleRestoreSuccess" />
  </div>
</template>

<script setup>
import StorageDonutChart from '@/components/StorageDonutChart.vue'; // <-- IMPORT CHART
import { ref, computed, onMounted } from 'vue';
import monitoringService from '@/services/monitoringService';
import QuickBackupModal from '@/components/CreateBackup.vue';
import QuickRestoreModal from '@/components/CreateRestore.vue';

// State
const remotes = ref([]);
const jobs = ref([]);
const isLoadingData = ref(true); 
const showBackupModal = ref(false);
const showRestoreModal = ref(false);

// Computed
const totalRemotes = computed(() => remotes.value.length);
const totalJobs = computed(() => jobs.value.length);

// Hitung "Used vs Free"
const storageStats = computed(() => {
  if (!Array.isArray(remotes.value) || remotes.value.length === 0) {
    return { series: [0, 100], labels: ['Terpakai', 'Kosong'], totalAvailable: 0 };
  }

  let totalAvailableGB = 0;
  let totalUsedGB = 0;
  
  remotes.value.forEach(remote => {
    // Hanya hitung remote yang terhubung
    if (remote.status_connect === 'CONNECTED') { 
      totalAvailableGB += remote.total_storage_gb || 0;
      totalUsedGB += remote.used_storage_gb || 0;
    }
  });
  
  // Hitung sisa (kosong)
  const totalFreeGB = totalAvailableGB - totalUsedGB;

  return { 
    series: [totalUsedGB, totalFreeGB > 0 ? totalFreeGB : 0], // [Terpakai, Kosong]
    labels: ['Terpakai', 'Kosong'],
    totalAvailable: totalAvailableGB // Total kapasitas dari semua remote
  };
});

const nextJob = computed(() => {
    if (!Array.isArray(jobs.value) || jobs.value.length === 0) return null;
    const now = new Date();
    const futureJobs = jobs.value
        .filter(j => j.next_run && j.next_run !== 'N/A' && j.next_run !== '')
        .map(j => {
            const nextRunDate = new Date(j.next_run);
            return { ...j, nextRunDate: nextRunDate, timeDiff: nextRunDate - now };
        })
        .filter(j => !isNaN(j.nextRunDate.getTime()) && j.timeDiff > 0)
        .sort((a, b) => a.timeDiff - b.timeDiff);
    if (futureJobs.length > 0) {
        const next = futureJobs[0];
        return {
            jobType: next.type,
            nextRun: next.nextRunDate,
            name: next.job_name,
            remoteName: next.gdrive_target,
            scheduleCron: next.next_run
        };
    }
    return null;
});

// Lifecycle
onMounted(async () => {
  await fetchData();
});

// Functions
async function fetchData() {
  isLoadingData.value = true;
  try {
    // Ambil data remotes (yang sudah ada rinciannya)
    const [remotesData, jobsData] = await Promise.all([
      monitoringService.getRemoteStatus(), 
      monitoringService.getScheduledJobs()
    ]);
    remotes.value = Array.isArray(remotesData) ? remotesData : [];
    jobs.value = Array.isArray(jobsData) ? jobsData : [];
  } catch (error) {
    console.error('Failed to fetch dashboard data:', error);
  } finally {
    isLoadingData.value = false;
  }
}

function openBackupModal() { showBackupModal.value = true; }
function openRestoreModal() { showRestoreModal.value = true; }
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
    return date.toLocaleString('id-ID', { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' });
  } catch (e) { return timestamp; }
}
</script>

<style scoped>
/* ========================================= */
/* MENGGUNAKAN CSS VARIABLES DARI APPLAYOUT */
/* ========================================= */
.dashboard-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 2rem 1.5rem;
}
.header {
  margin-bottom: 2.5rem;
}
h1 {
  font-size: 1.75rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 0.5rem 0;
  letter-spacing: -0.02em;
}
.subtitle {
  font-size: 0.95rem;
  color: var(--text-secondary);
  margin: 0;
  font-weight: 400;
}
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
  margin-bottom: 2rem;
}
.stat-card {
  background: var(--bg-card);
  padding: 1.5rem;
  border-radius: 8px;
  border: 1px solid var(--border-color);
}
.stat-content {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}
.chart-card-content {
  min-height: 250px; 
  justify-content: center;
}
.stat-label {
  font-size: 0.8125rem;
  color: var(--text-secondary);
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}
.stat-value {
  font-size: 2rem;
  font-weight: 600;
  color: var(--text-primary);
  line-height: 1;
}

/* CSS BARU UNTUK RINCIAN DRIVE */
.drive-breakdown-list {
  margin-top: 1.25rem;
  padding-top: 1.25rem;
  border-top: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}
.drive-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.9rem;
}
.drive-name {
  color: var(--text-primary);
  font-weight: 500;
}
.drive-size {
  color: var(--text-secondary);
  font-family: monospace;
}
.drive-empty {
  font-size: 0.9rem;
  color: var(--text-secondary);
  font-style: italic;
}
.status-message.small {
  padding: 0;
  font-size: 0.9rem;
  background: transparent;
}


/* Section */
.section {
  background: var(--bg-card);
  padding: 1.5rem;
  border-radius: 8px;
  border: 1px solid var(--border-color);
  margin-bottom: 1.5rem;
}
.section h2 {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 0.25rem 0;
}
.section-subtitle {
  font-size: 0.875rem;
  color: var(--text-secondary);
  margin: 0 0 1.5rem 0;
}
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
  border: 1px solid var(--border-color);
  border-radius: 6px;
  background: var(--bg-card);
  cursor: pointer;
  transition: all 0.2s;
  text-align: left;
  width: 100%;
}
.action-card:hover {
  background: var(--bg-main);
  border-color: var(--text-primary);
}
.action-content h3 {
  font-size: 1rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 0.25rem 0;
}
.action-content p {
  margin: 0;
  font-size: 0.875rem;
  color: var(--text-secondary);
}
.action-arrow {
  font-size: 1.25rem;
  color: var(--text-secondary);
  transition: transform 0.2s;
}
.action-card:hover .action-arrow {
  transform: translateX(4px);
  color: var(--text-primary);
}
.next-job-card {
  background: var(--bg-main);
  padding: 1.25rem;
  border-radius: 6px;
  border-left: 3px solid var(--text-primary);
}
.job-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.875rem;
}
.job-type-badge {
  background: var(--text-primary);
  color: var(--bg-card);
  padding: 0.25rem 0.625rem;
  border-radius: 4px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}
.job-time {
  color: var(--text-secondary);
  font-weight: 500;
  font-size: 0.875rem;
}
.next-job-card h3 {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--text-primary);
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
  color: var(--text-secondary);
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}
.detail-value {
  font-size: 0.9375rem;
  color: var(--text-primary);
  font-weight: 500;
  font-family: monospace;
}
.status-message {
  padding: 1rem;
  border-radius: 6px;
  font-size: 0.9375rem;
  background: var(--bg-main);
  color: var(--text-secondary);
  display: flex;
  align-items: center;
  gap: 0.75rem;
}
.loading-dot {
  width: 8px;
  height: 8px;
  background: var(--text-secondary);
  border-radius: 50%;
  animation: pulse 1.5s ease-in-out infinite;
}
@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.3; }
}
.empty-state {
  text-align: center;
  padding: 2rem 1rem;
}
.empty-state p {
  margin: 0;
  color: var(--text-secondary);
  font-size: 0.9375rem;
}
</style>