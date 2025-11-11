<template>
  <div class="manual-jobs-view">
    <div class="page-header">
      <div>
        <h1>🔧 Manual Jobs (Templates)</h1>
        <p>Menjalankan Job (Backup/Restore) yang telah disimpan tanpa jadwal.</p>
      </div>
    </div>
    
    <!-- ⭐ Script Preview Modal dengan Syntax Highlighting -->
    <ScriptPreview
      :is-visible="isModalVisible"
      :job-id="currentJobId"
      :script-content="currentScript"
      @close="closeModal"
    />

    <div v-if="isLoading" class="loading">
      <div class="spinner"></div>
      <p>Memuat data manual jobs...</p>
    </div>
    
    <div v-if="errorMessage" class="error-box">
      <span class="error-icon">⚠️</span>
      <div>
        <strong>Error:</strong> {{ errorMessage }}
        <button @click="fetchData" class="retry-btn">🔄 Coba Lagi</button>
      </div>
    </div>

    <div v-if="!isLoading && jobs.length === 0 && !errorMessage" class="empty-state">
      <div class="empty-icon">📋</div>
      <h3>Tidak Ada Job Manual</h3>
      <p>Tidak ada template job manual yang ditemukan.</p>
      <router-link to="/create" class="btn-create">
        ➕ Buat Job Baru
      </router-link>
    </div>

    <div v-if="!isLoading && jobs.length > 0" class="table-container">
      <table class="jobs-table">
        <thead>
          <tr>
            <th>Nama Job</th>
            <th>Object</th>
            <th>Target GDrive</th>
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
          />
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import jobService from '@/services/jobService';
import JobRow from '@/components/ManualJobRow.vue';
import ScriptPreview from '@/components/ScriptPreview.vue'; // ⭐ Import

const jobs = ref([]);
const isLoading = ref(true);
const errorMessage = ref(null);

// ⭐ State untuk Script Preview Modal
const isModalVisible = ref(false);
const currentScript = ref('');
const currentJobId = ref(null);

// Fetch data
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

// Trigger job
async function handleTrigger(jobId) {
  if (!confirm(`Apakah Anda yakin ingin menjalankan job manual ID ${jobId} sekarang?`)) {
    return;
  }
  
  try {
    await jobService.triggerManualJob(jobId);
    
    // Success notification (bisa diganti dengan toast)
    alert(`✅ Job ${jobId} berhasil di-trigger! Status akan diperbarui.`);
    
    // Reload data
    fetchData();
  } catch (error) {
    const errorMsg = error.response?.data?.error || error.message || String(error);
    alert(`❌ Gagal men-trigger job: ${errorMsg}`);
  }
}

// ⭐ View script dengan modal
async function handleViewScript(jobId) {
  try {
    // Show loading state
    isModalVisible.value = true;
    currentScript.value = 'Loading script...';
    currentJobId.value = jobId;
    
    // Fetch script
    const data = await jobService.getJobScript(jobId);
    
    // Update content
    currentScript.value = data.script_preview || data.script || 'No script available';
  } catch (error) {
    // Close modal and show error
    isModalVisible.value = false;
    
    const errorMsg = error.response?.data?.error || error.message;
    alert(`❌ Gagal memuat pratinjau script: ${errorMsg}`);
  }
}

// Close modal
function closeModal() {
  isModalVisible.value = false;
  currentScript.value = '';
  currentJobId.value = null;
}
</script>

<style scoped>
.manual-jobs-view {
  padding: 2rem;
  max-width: 1400px;
  margin: 0 auto;
}

/* Page Header */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
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

/* Loading State */
.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 3rem;
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

/* Error Box */
.error-box {
  background: #fee;
  border: 1px solid #fcc;
  border-radius: 8px;
  padding: 1rem 1.5rem;
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.error-icon {
  font-size: 1.5rem;
}

.error-box strong {
  color: #c33;
}

.retry-btn {
  background: #dc3545;
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.9rem;
  margin-left: 1rem;
  transition: all 0.2s;
}

.retry-btn:hover {
  background: #c82333;
  transform: translateY(-1px);
}

/* Empty State */
.empty-state {
  text-align: center;
  padding: 4rem 2rem;
}

.empty-icon {
  font-size: 4rem;
  margin-bottom: 1rem;
  opacity: 0.5;
}

.empty-state h3 {
  font-size: 1.5rem;
  color: #2c3e50;
  margin-bottom: 0.5rem;
}

.empty-state p {
  color: #6c757d;
  margin-bottom: 2rem;
}

.btn-create {
  display: inline-block;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 0.75rem 1.5rem;
  border-radius: 8px;
  text-decoration: none;
  font-weight: 600;
  transition: all 0.3s;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.btn-create:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.4);
}

/* Table Container */
.table-container {
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  overflow: hidden;
}

.jobs-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.9rem;
}

.jobs-table thead {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.jobs-table th {
  padding: 1rem 1.25rem;
  text-align: left;
  font-weight: 600;
  letter-spacing: 0.5px;
  text-transform: uppercase;
  font-size: 0.8rem;
}

.jobs-table td {
  padding: 1rem 1.25rem;
  border-bottom: 1px solid #f0f0f0;
}

.jobs-table tbody tr {
  transition: background-color 0.2s;
}

.jobs-table tbody tr:hover {
  background-color: #f8f9fa;
}

.jobs-table tbody tr:last-child td {
  border-bottom: none;
}

/* Responsive */
@media (max-width: 768px) {
  .manual-jobs-view {
    padding: 1rem;
  }
  
  .page-header h1 {
    font-size: 1.5rem;
  }
  
  .table-container {
    overflow-x: auto;
  }
  
  .jobs-table {
    font-size: 0.85rem;
  }
  
  .jobs-table th,
  .jobs-table td {
    padding: 0.75rem;
  }
}
</style>