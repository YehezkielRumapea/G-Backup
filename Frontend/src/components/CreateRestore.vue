<template>
  <div class="restore-view">
    <form @submit.prevent="handleRestoreSubmit" class="config-form restore-form">
      <h2>🔄 Restore Configuration</h2>
      <p>Download file dari cloud storage dan restore ke path lokal.</p>
      
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
        <button type="submit" :disabled="isLoading" class="btn-submit btn-restore">
          <span v-if="isLoading">⏳ Starting Restore...</span>
          <span v-else>🔄 Start Restore</span>
        </button>
      </div>
    </form>

    <transition name="fade">
      <div v-if="message" class="message success">
        <span class="message-icon">✅</span>
        {{ message }}
      </div>
    </transition>

    <transition name="fade">
      <div v-if="errorMessage" class="message error">
        <span class="message-icon">❌</span>
        {{ errorMessage }}
      </div>
    </transition>
  </div>
  <div v-if="props.isVisible" class="modal-overlay" @click.self="emit('close')">
    <div class="modal-content">
      <button type="button" @click="emit('close')" class="modal-close-btn">✕</button>
      
      <div class="restore-view">
        </div>
      
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import jobService from '@/services/jobService'; // Pastikan path ini benar
import { useRouter } from 'vue-router';

const router = useRouter();
const isLoading = ref(false);
const errorMessage = ref(null);
const message = ref(null);

// Di dalam <script setup> CreateRestore.vue
const props = defineProps({
    isVisible: {
        type: Boolean,
        default: false,
        required: true // Pastikan ini ada
    },
    // ... props lain jika ada
});

// Dan definisikan emits
const emit = defineEmits(['close', 'success']);

const restoreForm = ref({
  remote_name: '',
  source_path: '',
  destination_path: ''
});

async function handleRestoreSubmit() {
  isLoading.value = true;
  errorMessage.value = null;
  message.value = null;

  try {
    const response = await jobService.createRestoreJob(restoreForm.value);
    message.value = response.message || 'Restore job started successfully!';
    
    // Redirect setelah 1.5 detik
    setTimeout(() => {
      router.push('/logs');
    }, 1500);
    
  } catch (error) {
    console.error('Create restore job error:', error);
    errorMessage.value = error.response?.data?.error || 'Failed to start restore.';
  } finally {
    isLoading.value = false;
  }
}
</script>

<style scoped>
/* Anda mungkin perlu menyesuaikan `max-width` jika ini adalah komponen utama di halaman Anda */
.restore-view {
  max-width: 900px;
  margin: 0 auto;
  padding: 2rem;
}

/* Form Styles */
.config-form {
  background: #fff;
  padding: 2rem;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  margin-bottom: 2rem;
}

.config-form h2 {
  margin: 0 0 0.5rem 0;
  font-size: 1.5rem;
  font-weight: 700;
  color: #2c3e50;
}

.config-form > p {
  color: #6c757d;
  margin: 0 0 2rem 0;
  font-size: 0.95rem;
}

.restore-form {
  border-top: 4px solid #3498db;
}

/* Form Group */
.form-group {
  margin-bottom: 1.5rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 600;
  color: #2c3e50;
  font-size: 0.95rem;
}

.form-group input[type="text"],
.form-group input[type="number"],
.form-group input[type="time"],
.form-group select,
.form-group textarea {
  width: 100%;
  padding: 0.75rem;
  border: 2px solid #e9ecef;
  border-radius: 8px;
  font-size: 0.95rem;
  transition: all 0.2s;
  box-sizing: border-box;
}

.form-group input:focus,
.form-group select:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.hint {
  display: block;
  margin-top: 0.5rem;
  font-size: 0.85rem;
  color: #6c757d;
}

/* Form Actions */
.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  margin-top: 2rem;
  padding-top: 1.5rem;
  border-top: 2px solid #e9ecef;
}

.btn-submit {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 0.75rem 2rem;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 600;
  font-size: 1rem;
  transition: all 0.2s;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.btn-submit:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.btn-submit.btn-restore {
  background: linear-gradient(135deg, #3498db 0%, #2980b9 100%);
  box-shadow: 0 4px 12px rgba(52, 152, 219, 0.3);
}

.btn-submit.btn-restore:hover:not(:disabled) {
  box-shadow: 0 6px 20px rgba(52, 152, 219, 0.4);
}

/* Messages */
.message {
  padding: 1rem 1.5rem;
  border-radius: 8px;
  margin-top: 1.5rem;
  display: flex;
  align-items: center;
  gap: 1rem;
  font-weight: 500;
}

.message.success {
  background: #d4edda;
  color: #155724;
  border-left: 4px solid #28a745;
}

.message.error {
  background: #f8d7da;
  color: #721c24;
  border-left: 4px solid #dc3545;
}

.message-icon {
  font-size: 1.5rem;
}

/* Transitions */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* Responsive */
@media (max-width: 768px) {
  .restore-view {
    padding: 1rem;
  }
  
  .config-form {
    padding: 1.5rem;
  }
  
  .form-actions {
    flex-direction: column;
  }
  
  .btn-submit {
    width: 100%;
  }
}
</style>