<template>
  <transition name="fade">
    <div v-if="props.isVisible" class="modal-overlay" @click.self="close">
      <div class="modal-content">
        <!-- Header -->
        <div class="modal-header">
          <h2>Restore Configuration</h2>
          <button type="button" class="close-btn" @click="close">×</button>
        </div>

        <!-- Form & Messages -->
        <div class="restore-view">
          <form @submit.prevent="handleRestoreSubmit" class="config-form">
            <p class="form-description">Download file dari cloud storage dan restore ke path lokal.</p>

            <div class="form-group">
              <label for="restore-remote">Remote Name (Source) *</label>
              <input 
                type="text" 
                id="restore-remote" 
                v-model="restoreForm.remote_name" 
                required 
                placeholder="Gdrive1"
              />
            </div>

            <div class="form-group">
              <label for="restore-source">Source Path (Cloud) *</label>
              <input 
                type="text" 
                id="restore-source" 
                v-model="restoreForm.source_path" 
                required 
                placeholder="/backups/database/backup_20240110.sql.gz"
              />
            </div>

            <div class="form-group">
              <label for="restore-dest">Destination Path (Local Server) *</label>
              <input 
                type="text" 
                id="restore-dest" 
                v-model="restoreForm.destination_path" 
                required 
                placeholder="/home/user/restore/"
              />
            </div>

            <div class="form-actions">
              <button type="submit" :disabled="isLoading" class="btn-submit">
                <span v-if="isLoading">Starting Restore...</span>
                <span v-else>Start Restore</span>
              </button>
            </div>
          </form>

          <transition name="fade">
            <div v-if="message" class="message success">
              {{ message }}
            </div>
          </transition>

          <transition name="fade">
            <div v-if="errorMessage" class="message error">
              {{ errorMessage }}
            </div>
          </transition>
        </div>
      </div>
    </div>
  </transition>
</template>

<script setup>
import { ref } from 'vue';
import jobService from '@/services/jobService';
import { useRouter } from 'vue-router';

// Props & Emits
const props = defineProps({
  isVisible: { type: Boolean, default: false }
});
const emit = defineEmits(['close', 'success']);

const router = useRouter();
const isLoading = ref(false);
const errorMessage = ref(null);
const message = ref(null);

function close() {
  emit('close');
}

const restoreForm = ref({
  remote_name: '',
  source_path: '',
  destination_path: ''
});

async function handleRestoreSubmit() {
  isLoading.value = true;
  message.value = null;
  errorMessage.value = null;

  try {
    const response = await jobService.createRestoreJob(restoreForm.value);
    message.value = response.message || 'Restore job started successfully!';
    emit('success');

    setTimeout(() => close(), 1500);
    setTimeout(() => router.push('/logs'), 1500);
  } catch (error) {
    console.error('Create restore job error:', error);
    errorMessage.value = error.response?.data?.error || 'Failed to start restore.';
  } finally {
    isLoading.value = false;
  }
}
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
  max-width: 600px;
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
  max-width: 100%;
}

/* Form */
.config-form {
  padding: 1.5rem;
}

.form-description {
  color: #666;
  margin: 0 0 1.5rem 0;
  font-size: 0.9375rem;
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
  border-color: #1a1a1a;
}

/* Form Actions */
.form-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 1.5rem;
  padding-top: 1.25rem;
  border-top: 1px solid #e5e5e5;
}

.btn-submit {
  background: #1a1a1a;
  color: white;
  padding: 0.625rem 1.5rem;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 500;
  font-size: 0.9375rem;
  transition: all 0.2s;
}

.btn-submit:hover:not(:disabled) {
  background: #333;
}

.btn-submit:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* Messages */
.message {
  padding: 0.875rem 1rem;
  border-radius: 6px;
  margin: 0 1.5rem 1.5rem;
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
  }
  
  .config-form {
    padding: 1rem;
  }
  
  .form-actions {
    flex-direction: column;
  }
  
  .btn-submit {
    width: 100%;
  }
}
</style>