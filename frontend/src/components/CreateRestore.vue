<template>
  <transition name="fade">
    <div v-if="props.isVisible" class="modal-overlay" @click.self="close">
      <div class="modal-content">
        <!-- Header -->
        <div class="modal-header">
          <h2>Restore Configuration</h2>
          <button type="button" class="close-btn" @click="close">×</button>
        </div>

        <!-- Content -->
        <div class="restore-view">
          <!-- Step Indicator -->
          <div class="step-indicator">
            <div :class="['step', { active: currentStep >= 1, completed: currentStep > 1 }]">
              <div class="step-number">1</div>
              <div class="step-label">Select Remote</div>
            </div>
            <div class="step-line" :class="{ completed: currentStep > 1 }"></div>
            
            <div :class="['step', { active: currentStep >= 2, completed: currentStep > 2 }]">
              <div class="step-number">2</div>
              <div class="step-label">Select File</div>
            </div>
            <div class="step-line" :class="{ completed: currentStep > 2 }"></div>
            
            <div :class="['step', { active: currentStep >= 3 }]">
              <div class="step-number">3</div>
              <div class="step-label">Destination</div>
            </div>
          </div>

          <!-- Step 1: Select Remote Drive -->
          <div v-show="currentStep === 1" class="step-content">
            <h3>☁️ Pilih Cloud Storage</h3>
            <p class="form-description">Pilih drive cloud mana yang berisi file backup</p>
            
            <div v-if="loadingRemotes" class="loading">
              <div class="spinner"></div>
              <p>Loading remotes...</p>
            </div>

            <div v-else-if="remotesList.length === 0" class="empty-state">
              <p>❌ Tidak ada remote ditemukan</p>
            </div>

            <div v-else class="remotes-grid">
              <button 
                v-for="remote in remotesList"
                :key="remote.name"
                @click="selectRemote(remote.name)"
                :class="['remote-card', { selected: selectedRemote === remote.name }]"
                type="button"
              >
                <div class="remote-icon">☁️</div>
                <div class="remote-name">{{ remote.name }}</div>
              </button>
            </div>

            <div class="step-actions">
              <button type="button" @click="close" class="btn-secondary">Batal</button>
              <button 
                type="button"
                @click="goToStep(2)" 
                class="btn-submit"
                :disabled="!selectedRemote"
              >
                Lanjut →
              </button>
            </div>
          </div>

          <!-- Step 2: Select File/Folder -->
          <div v-show="currentStep === 2" class="step-content">
            <h3>📦 Pilih File atau Folder</h3>
            <p class="form-description">Pilih file/folder dari {{ selectedRemote }} yang ingin di-restore</p>
            
            <div class="file-browser-wrapper">
              <GDriveBrowser
                v-if="selectedRemote"
                :remote-name="selectedRemote"
                :initial-path="'/'"
                @select-file="handleSelectFile"
              />
            </div>

            <div v-if="selectedFile" class="selected-info">
              <div class="info-item">
                <strong>✅ Selected:</strong> {{ selectedFile.name }}
              </div>
              <div class="info-item">
                <strong>Path:</strong> {{ selectedFile.path }}
              </div>
              <div class="info-item">
                <strong>Type:</strong> {{ selectedFile.is_dir ? '📁 Folder' : '📄 File' }}
              </div>
            </div>

            <div class="step-actions">
              <button type="button" @click="goToStep(1)" class="btn-secondary">← Kembali</button>
              <button 
                type="button"
                @click="goToStep(3)" 
                class="btn-submit"
                :disabled="!selectedFile"
              >
                Lanjut →
              </button>
            </div>
          </div>

          <!-- Step 3: Set Destination Path -->
          <div v-show="currentStep === 3" class="step-content">
            <h3>🎯 Lokasi Penyimpanan</h3>
            <p class="form-description">Tentukan lokasi di server lokal untuk menyimpan file yang di-restore</p>

            <!-- Review -->
            <div class="review-box">
              <div class="review-item">
                <strong>📍 Remote:</strong>
                <span>{{ selectedRemote }}</span>
              </div>
              <div class="review-item">
                <strong>📦 File:</strong>
                <span>{{ selectedFile?.name }}</span>
              </div>
              <div class="review-item">
                <strong>Path:</strong>
                <span>{{ selectedFile?.path }}</span>
              </div>
            </div>

            <form @submit.prevent="handleRestoreSubmit" class="config-form">
              <!-- Destination Path Input -->
              <div class="form-group">
                <label for="restore-dest">Lokasi Penyimpanan di Server *</label>
                <input 
                  type="text" 
                  id="restore-dest" 
                  v-model="destinationPath" 
                  required 
                  placeholder="/home/user/restore/"
                />
                <small class="hint">Contoh: /home/user/restore atau /opt/backups</small>
              </div>

              <!-- Error Message -->
              <transition name="fade">
                <div v-if="errorMessage" class="message error">
                  {{ errorMessage }}
                </div>
              </transition>

              <!-- Success Message -->
              <transition name="fade">
                <div v-if="successMessage" class="message success">
                  {{ successMessage }}
                </div>
              </transition>

              <!-- Form Actions -->
              <div class="step-actions">
                <button type="button" @click="goToStep(2)" class="btn-secondary" :disabled="isLoading">
                  ← Kembali
                </button>
                <button 
                  type="submit" 
                  :disabled="isLoading || !destinationPath" 
                  class="btn-restore"
                >
                  <span v-if="isLoading">Memulai Restore...</span>
                  <span v-else>🚀 Mulai Restore</span>
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
  </transition>
</template>

<script setup>
import { ref, watch } from 'vue'
import GDriveBrowser from '@/components/GDriveBrowser.vue'
import jobService from '@/services/jobService'
import { useRouter } from 'vue-router'
import driveService from '../services/driveService'

// Props & Emits
const props = defineProps({
  isVisible: { type: Boolean, default: false }
})
const emit = defineEmits(['close', 'success'])

const router = useRouter()

// Step management
const currentStep = ref(1)

// Step 1: Remote
const loadingRemotes = ref(false)
const remotesList = ref([])
const selectedRemote = ref('')

// Step 2: File
const selectedFile = ref(null)

// Step 3: Destination
const destinationPath = ref('')

// API state
const isLoading = ref(false)
const errorMessage = ref(null)
const successMessage = ref(null)

// ============================================
// Methods
// ============================================

// Load remotes
// CreateRestore.vue (Fungsi loadRemotes yang Dikoreksi)

async function loadRemotes() {
 loadingRemotes.value = true
 try {
  const response = await driveService.listRemotes()
    
    // Asumsi: Kita menggunakan solusi koreksi sebelumnya di mana response adalah array atau object
    // Jika backend mengembalikan array langsung:
    remotesList.value = Array.isArray(response) ? response : (response.remotes || []) 

 } catch (error) { 
  console.error('Failed to load remotes:', error)
    
    // Hapus baris yang salah. Cukup set ke array kosong.
  remotesList.value = [] 
    
    // Opsional: Jika Anda ingin menampilkan error ke UI (remotesList.value.length === 0)
    // Anda bisa mengatur variabel error baru di state.
    
 } finally {
  loadingRemotes.value = false
 }
}

// Select remote
function selectRemote(remoteName) {
  selectedRemote.value = remoteName
}

// Handle file selection from browser
function handleSelectFile(file) {
  selectedFile.value = file
}

// Navigate to step
function goToStep(step) {
  if (Math.abs(step - currentStep.value) === 1 || step < currentStep.value) {
    currentStep.value = step
  }
}

// Submit restore
async function handleRestoreSubmit() {
  if (!selectedRemote.value || !selectedFile.value || !destinationPath.value) {
    errorMessage.value = 'Semua field harus diisi'
    return
  }

  isLoading.value = true
  errorMessage.value = null
  successMessage.value = null

  try {
    const response = await jobService.createRestoreJob({
      source_remote: selectedRemote.value,
      source_path: selectedFile.value.path,
      destination_path: destinationPath.value,
      job_name: `Restore ${selectedFile.value.name}`
    })

    successMessage.value = response.message || 'Restore job berhasil dimulai!'
    emit('success')

    // Close modal dan redirect
    setTimeout(() => close(), 1500)
    setTimeout(() => router.push('/logs'), 1500)
  } catch (error) {
    console.error('Trigger restore error:', error)
    errorMessage.value = error.response?.data?.error || 'Gagal memulai restore'
  } finally {
    isLoading.value = false
  }
}

// Close modal
function close() {
  emit('close')
  resetForm()
}

// Reset form
function resetForm() {
  currentStep.value = 1
  selectedRemote.value = ''
  selectedFile.value = null
  destinationPath.value = ''
  errorMessage.value = null
  successMessage.value = null
}

// Load remotes when modal opens
watch(() => props.isVisible, (newVal) => {
  if (newVal) {
    loadRemotes()
  }
})
</script>

<style scoped>
/* Modal */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 999;
  padding: 20px;
}

.modal-content {
  background: #fff;
  border-radius: 8px;
  border: 1px solid #e5e5e5;
  width: 90%;
  max-width: 700px;
  max-height: 90vh;
  overflow-y: auto;
  position: relative;
}

/* Header */
.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.25rem 1.5rem;
  background: #fafafa;
  border-bottom: 1px solid #e5e5e5;
  position: sticky;
  top: 0;
  z-index: 10;
}

.modal-header h2 {
  margin: 0;
  font-size: 1.125rem;
  font-weight: 600;
  color: #1a1a1a;
}

.close-btn {
  background: transparent;
  border: 1px solid #e5e5e5;
  color: #666;
  width: 32px;
  height: 32px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 1.5rem;
  line-height: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.close-btn:hover {
  background: #f5f5f5;
  border-color: #1a1a1a;
  color: #1a1a1a;
}

/* Restore View */
.restore-view {
  padding: 1.5rem;
}

/* Step Content */
.step-content {
  animation: fadeIn 0.3s ease-in;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.step-content h3 {
  margin: 0 0 0.5rem 0;
  font-size: 1.25rem;
  color: #1a1a1a;
}

.form-description {
  color: #666;
  margin: 0 0 1.5rem 0;
  font-size: 0.9375rem;
}

/* Step Indicator */
.step-indicator {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 2rem;
  gap: 0.5rem;
}

.step {
  display: flex;
  flex-direction: column;
  align-items: center;
  flex: 1;
}

.step-number {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: #f0f0f0;
  border: 2px solid #e5e5e5;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  color: #666;
  font-size: 0.875rem;
  transition: all 0.3s;
}

.step.active .step-number {
  background: #3b82f6;
  border-color: #3b82f6;
  color: white;
}

.step.completed .step-number {
  background: #10b981;
  border-color: #10b981;
  color: white;
}

.step-label {
  font-size: 0.75rem;
  font-weight: 500;
  color: #666;
  margin-top: 0.375rem;
  text-align: center;
}

.step-line {
  flex: 1;
  height: 2px;
  background: #e5e5e5;
  transition: all 0.3s;
}

.step-line.completed {
  background: #10b981;
}

/* Remotes Grid */
.remotes-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.remote-card {
  background: #000000;
  border: 2px solid #ffffff;
  border-radius: 8px;
  padding: 1.25rem;
  cursor: pointer;
  transition: all 0.2s;
  text-align: center;
}

.remote-card:hover {
  border-color: #3b82f6;
  background: #f0f8ff;
}

.remote-card.selected {
  background: #3b82f6;
  border-color: #3b82f6;
  color: rgb(0, 0, 0);
}

.remote-icon {
  font-size: 2rem;
  margin-bottom: 0.5rem;
}

.remote-name {
  font-weight: 600;
  font-size: 0.9375rem;
}

/* File Browser */
.file-browser-wrapper {
  margin-bottom: 1.5rem;
  border: 1px solid #e5e5e5;
  border-radius: 8px;
  overflow: hidden;
}

/* Selected Info */
.selected-info {
  background: #f0fdf4;
  border: 1px solid #86efac;
  border-radius: 6px;
  padding: 1rem;
  margin-bottom: 1.5rem;
}

.info-item {
  font-size: 0.875rem;
  padding: 0.25rem 0;
  color: #166534;
}

.info-item strong {
  font-weight: 600;
}

/* Review Box */
.review-box {
  background: #f9f9f9;
  border: 1px solid #e5e5e5;
  border-radius: 8px;
  padding: 1rem;
  margin-bottom: 1.5rem;
}

.review-item {
  display: flex;
  justify-content: space-between;
  padding: 0.5rem 0;
  border-bottom: 1px solid #e5e5e5;
  font-size: 0.875rem;
}

.review-item:last-child {
  border-bottom: none;
}

.review-item strong {
  color: #1a1a1a;
  min-width: 100px;
}

.review-item span {
  color: #3b82f6;
  font-weight: 500;
  text-align: right;
}

/* Form */
.config-form {
  width: 100%;
}

/* Form Group */
.form-group {
  margin-bottom: 1.25rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
  color: #1a1a1a;
  font-size: 0.875rem;
}

.form-group input[type="text"] {
  width: 100%;
  padding: 0.625rem 0.875rem;
  border: 1px solid #e5e5e5;
  border-radius: 6px;
  font-size: 0.9375rem;
  transition: all 0.2s;
  box-sizing: border-box;
}

.form-group input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.hint {
  display: block;
  margin-top: 0.375rem;
  font-size: 0.8125rem;
  color: #666;
}

/* Messages */
.message {
  padding: 0.875rem 1rem;
  border-radius: 6px;
  margin-bottom: 1rem;
  font-weight: 500;
  font-size: 0.9375rem;
}

.message.success {
  background: #d1fae5;
  color: #065f46;
  border-left: 3px solid #22c55e;
}

.message.error {
  background: #fee2e2;
  color: #991b1b;
  border-left: 3px solid #ef4444;
}

/* Loading */
.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 2rem;
  gap: 1rem;
}

.spinner {
  width: 36px;
  height: 36px;
  border: 3px solid #e5e5e5;
  border-top-color: #3b82f6;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.empty-state {
  text-align: center;
  padding: 2rem;
  color: #999;
  font-size: 0.9375rem;
}

/* Step Actions */
.step-actions {
  display: flex;
  justify-content: space-between;
  gap: 1rem;
  margin-top: 1.5rem;
  padding-top: 1.25rem;
  border-top: 1px solid #e5e5e5;
}

.btn-secondary,
.btn-submit,
.btn-restore {
  padding: 0.625rem 1.5rem;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 500;
  font-size: 0.9375rem;
  transition: all 0.2s;
}

.btn-secondary {
  background: #f5f5f5;
  color: #666;
  border: 1px solid #e5e5e5;
}

.btn-secondary:hover:not(:disabled) {
  background: #e5e5e5;
}

.btn-submit {
  background: #1a1a1a;
  color: white;
  flex: 1;
}

.btn-submit:hover:not(:disabled) {
  background: #333;
}

.btn-restore {
  background: #ef4444;
  color: white;
  flex: 1;
}

.btn-restore:hover:not(:disabled) {
  background: #dc2626;
}

.btn-secondary:disabled,
.btn-submit:disabled,
.btn-restore:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* Transitions */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* Responsive */
@media (max-width: 768px) {
  .modal-content {
    width: 95%;
    max-height: 95vh;
  }
  
  .restore-view {
    padding: 1rem;
  }

  .step-indicator {
    flex-direction: column;
    gap: 0;
    margin-bottom: 1.5rem;
  }

  .step-line {
    width: 2px;
    height: 20px;
  }

  .remotes-grid {
    grid-template-columns: 1fr;
  }

  .step-actions {
    flex-direction: column;
  }

  .btn-secondary,
  .btn-submit,
  .btn-restore {
    width: 100%;
  }
}
</style>