<template>
  <div class="manual-jobs-view">
    <div class="header">
      <div>
        <h1>Manual Jobs</h1>
        <p class="subtitle">Menjalankan Job Tanpa Jadwal yang telah disimpan</p>
      </div>
    </div>
    
    <ScriptPreview
      :is-visible="isModalVisible"
      :job-id="currentJobId"
      :script-content="currentScript"
      @close="closeModal"
    />

    <!-- ✅ Edit Job Modal - FIXED -->
    <EditJobModal
      :job-id="editingJobId"
      :is-open="isEditModalOpen"
      @close="closeEditModal"
      @success="handleEditSuccess"
    />

    <div v-if="isLoading" class="status-message">
      <span class="loading-dot"></span>
      Memuat data...
    </div>
    
    <div v-if="errorMessage" class="status-message error">
      {{ errorMessage }}
      <button @click="fetchData" class="retry-btn">Coba Lagi</button>
    </div>

    <div v-if="!isLoading && jobs.length === 0 && !errorMessage" class="empty-state">
      <p>Tidak ada template job manual yang ditemukan</p>
      <router-link to="/create" class="btn-create">
        Buat Job Baru
      </router-link>
    </div>

    <div v-if="!isLoading && jobs.length > 0" class="table-container">
      <table class="jobs-table">
        <thead>
          <tr>
            <th>Nama Job</th>
            <th>Object</th>
            <th>GDrive</th>
            <th>Last Run</th>
            <th>Status</th>
            <th>Action</th>
          </tr>
        </thead>
        <tbody>
          <JobRow
            v-for="job in jobs"
            :key="job.id"
            :job="job"
            :show-next-run="false"
            @trigger="handleTrigger"
            @view-script="handleViewScript"
            @edit="handleEdit"
            @delete="handleDelete"
          />
        </tbody>
      </table>
    </div>

    <!-- ✅ Toast Notification -->
    <Transition name="toast">
      <div
        v-if="toast.show"
        :class="getToastClass(toast.type)"
        class="toast-notification"
      >
        <div class="toast-content">
          <svg v-if="toast.type === 'success'" class="toast-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
          </svg>
          <svg v-else-if="toast.type === 'error'" class="toast-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
          <svg v-else class="toast-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <p class="toast-message">{{ toast.message }}</p>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import jobService from '@/services/jobService';
import JobRow from '@/components/ManualJobRow.vue';
import ScriptPreview from '@/components/ScriptPreview.vue';
import EditJobModal from '@/components/EditJobModal.vue';

const jobs = ref([]);
const isLoading = ref(true);
const errorMessage = ref(null);

const isModalVisible = ref(false);
const currentScript = ref('');
const currentJobId = ref(null);

// ✅ Edit modal state
const editingJobId = ref(null);
const isEditModalOpen = ref(false);

// ✅ Toast notification state
const toast = ref({
  show: false,
  message: '',
  type: 'success', // 'success' | 'error' | 'info'
}); 

// ✅ Fetch all manual jobs
async function fetchData() {
  isLoading.value = true;
  errorMessage.value = null;
  
  try {
    const data = await jobService.getManualJobs();
    jobs.value = data;
  } catch (error) {
    errorMessage.value = error.response?.data?.error || 'Gagal memuat data manual jobs.';
    console.error("Fetch Manual Jobs Error:", error);
  } finally {
    isLoading.value = false;
  }
}

onMounted(fetchData);

// ✅ Handle edit job - Open modal with job ID
function handleEdit(jobId) {
  console.log('✅ handleEdit dipanggil dengan jobId:', jobId);
  editingJobId.value = jobId;
  isEditModalOpen.value = true;
}

// ✅ Close edit modal
function closeEditModal() {
  console.log('✅ closeEditModal dipanggil');
  isEditModalOpen.value = false;
  editingJobId.value = null;
}

// ✅ Handle edit success - Reload jobs and show toast
function handleEditSuccess() {
  console.log('✅ handleEditSuccess dipanggil');
  fetchData(); // Reload jobs list
  showToast('Job berhasil diperbarui!', 'success');
}

// ✅ Handle delete job
async function handleDelete(jobId, jobName) {
  if (!confirm(`Apakah Anda yakin ingin menghapus job "${jobName}"?`)) {
    return;
  }

  try {
    await jobService.deleteJob(jobId);
    showToast('Job berhasil dihapus!', 'success');
    fetchData(); // Reload jobs list
  } catch (error) {
    const errorMsg = error.response?.data?.error || error.message || 'Gagal menghapus job';
    showToast(errorMsg, 'error');
  }
}

// ✅ Handle trigger job
async function handleTrigger(jobId) {
  if (!confirm(`Apakah Anda yakin ingin menjalankan job manual ID ${jobId} sekarang?`)) {
    return;
  }
  
  try {
    await jobService.triggerManualJob(jobId);
    showToast(`Job ${jobId} berhasil di-trigger! Status akan diperbarui.`, 'success');
    fetchData();
  } catch (error) {
    const errorMsg = error.response?.data?.error || error.message || String(error);
    showToast(`Gagal men-trigger job: ${errorMsg}`, 'error');
  }
}

// ✅ Handle view script
async function handleViewScript(jobId) {
  try {
    isModalVisible.value = true;
    currentScript.value = 'Loading script...';
    currentJobId.value = jobId;
    
    const data = await jobService.getJobScript(jobId);
    currentScript.value = data.script_preview || data.script || 'No script available';
  } catch (error) {
    isModalVisible.value = false;
    const errorMsg = error.response?.data?.error || error.message;
    showToast(`Gagal memuat pratinjau script: ${errorMsg}`, 'error');
  }
}

function closeModal() {
  isModalVisible.value = false;
  currentScript.value = '';
  currentJobId.value = null;
}

// ✅ Show toast notification
function showToast(message, type = 'success') {
  toast.value = { show: true, message, type };
  setTimeout(() => {
    toast.value.show = false;
  }, 3000);
}

// ✅ Get toast CSS class
function getToastClass(type) {
  const classes = {
    success: 'toast-success',
    error: 'toast-error',
    info: 'toast-info'
  };
  return classes[type] || 'toast-info';
}
</script>

<style scoped>
.manual-jobs-view {
  padding: 2rem 1rem;
}

.header {
  margin-bottom: 2rem;
}

.header h1 {
  margin: 0 0 0.5rem 0;
  font-size: 1.875rem;
  font-weight: bold;
  color: #1a1a1a;
}

.header .subtitle {
  margin: 0;
  color: #666;
  font-size: 0.9375rem;
}

/* Status Messages */
.status-message {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 1rem;
  background: #f0f9ff;
  border: 1px solid #bfdbfe;
  border-radius: 6px;
  color: #1e40af;
  margin-bottom: 1rem;
  font-size: 0.9375rem;
}

.status-message.error {
  background: #fee2e2;
  border-color: #fecaca;
  color: #991b1b;
  justify-content: space-between;
}

.loading-dot {
  width: 8px;
  height: 8px;
  background: currentColor;
  border-radius: 50%;
  animation: pulse 1.5s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.retry-btn {
  background: currentColor;
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.875rem;
  font-weight: 500;
}

.retry-btn:hover {
  opacity: 0.9;
}

/* Empty State */
.empty-state {
  text-align: center;
  padding: 3rem 1rem;
}

.empty-state p {
  margin: 0 0 1rem 0;
  color: #666;
  font-size: 0.9375rem;
}

.btn-create {
  display: inline-block;
  background: #1a1a1a;
  color: white;
  padding: 0.625rem 1.5rem;
  border-radius: 6px;
  text-decoration: none;
  font-weight: 500;
  font-size: 0.9375rem;
  transition: all 0.2s;
}

.btn-create:hover {
  background: #333;
}

/* Table */
.table-container {
  overflow-x: auto;
}

.jobs-table {
  width: 100%;
  border-collapse: collapse;
  background: white;
  border: 1px solid #e5e5e5;
  border-radius: 6px;
  overflow: hidden;
}

.jobs-table thead {
  background: #fafafa;
  border-bottom: 1px solid #e5e5e5;
}

.jobs-table th {
  padding: 1rem;
  text-align: left;
  font-weight: bold;
  font-size: 0.875rem;
  color: #1a1a1a;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.jobs-table td {
  padding: 1rem;
  border-bottom: 1px solid #e5e5e5;
  font-size: 0.9375rem;
  color: #666;
}

.jobs-table tbody tr:last-child td {
  border-bottom: none;
}

.jobs-table tbody tr:hover {
  background: #f9f9f9;
}

/* Toast Notification */
.toast-notification {
  position: fixed;
  bottom: 2rem;
  right: 2rem;
  background: white;
  border-radius: 8px;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1);
  z-index: 1001;
  animation: slideIn 0.3s ease-out;
}

@keyframes slideIn {
  from {
    transform: translateX(400px);
    opacity: 0;
  }
  to {
    transform: translateX(0);
    opacity: 1;
  }
}

.toast-notification.toast-success {
  border-left: 4px solid #22c55e;
}

.toast-notification.toast-error {
  border-left: 4px solid #ef4444;
}

.toast-notification.toast-info {
  border-left: 4px solid #3b82f6;
}

.toast-content {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem 1.5rem;
}

.toast-icon {
  width: 20px;
  height: 20px;
  flex-shrink: 0;
}

.toast-notification.toast-success .toast-icon {
  color: #22c55e;
}

.toast-notification.toast-error .toast-icon {
  color: #ef4444;
}

.toast-notification.toast-info .toast-icon {
  color: #3b82f6;
}

.toast-message {
  margin: 0;
  color: #1a1a1a;
  font-weight: 500;
  font-size: 0.9375rem;
  white-space: nowrap;
}

.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}

.toast-enter-from {
  transform: translateX(400px);
  opacity: 0;
}

.toast-leave-to {
  transform: translateX(400px);
  opacity: 0;
}

/* Responsive */
@media (max-width: 768px) {
  .manual-jobs-view {
    padding: 1rem;
  }

  .header h1 {
    font-size: 1.5rem;
  }

  .jobs-table th,
  .jobs-table td {
    padding: 0.75rem 0.5rem;
    font-size: 0.8125rem;
  }

  .toast-notification {
    bottom: 1rem;
    right: 1rem;
    left: 1rem;
  }
}
</style>