<template>
  <div class="create-restore-view">
    <form @submit.prevent="handleSubmit" class="restore-form">
      <h2>🔄 Restore Configuration</h2>
      <p>Browse and select file from Google Drive to restore to local server.</p>
      
      <!-- Step 1: Select Remote -->
      <div class="form-group">
        <label for="remote">
          <span class="step-number">1</span>
          Select Remote (Source) *
        </label>
        <select
          id="remote"
          v-model="form.remote_name"
          @change="handleRemoteSelect"
          required
          class="remote-select"
        >
          <option value="">-- Select Remote --</option>
          <option v-for="remote in remotes" :key="remote.RemoteName" :value="remote.RemoteName">
            {{ remote.RemoteName }}
          </option>
        </select>
      </div>

      <!-- Step 2: Browse & Select File -->
      <div v-if="form.remote_name" class="form-group">
        <label>
          <span class="step-number">2</span>
          Browse & Select File from Google Drive *
        </label>
        
        <div class="file-browser">
          <!-- Browser Header -->
          <div class="browser-header">
            <div class="breadcrumb">
              <button type="button" @click="navigateTo('/')" class="breadcrumb-btn">
                🏠 Home
              </button>
              <span v-for="(part, index) in pathParts" :key="index">
                <span class="separator">/</span>
                <button type="button" @click="navigateTo(getPathUpTo(index))" class="breadcrumb-btn">
                  {{ part }}
                </button>
              </span>
            </div>
            <button 
              type="button"
              @click="loadFiles" 
              class="refresh-btn-small"
              :disabled="isLoadingFiles"
            >
              <span :class="{ spinning: isLoadingFiles }">🔄</span>
            </button>
          </div>

          <!-- Search Box -->
          <div class="search-box-mini">
            <span class="search-icon">🔍</span>
            <input 
              v-model="searchQuery"
              type="text"
              placeholder="Search files..."
              @input="handleSearch"
              class="search-input-mini"
            />
            <button 
              v-if="searchQuery" 
              type="button"
              @click="clearSearch" 
              class="clear-btn-mini"
            >
              ✕
            </button>
          </div>

          <!-- Loading -->
          <div v-if="isLoadingFiles" class="loading-mini">
            <div class="spinner-mini"></div>
            <p>Loading files...</p>
          </div>

          <!-- Error -->
          <div v-else-if="filesError" class="error-mini">
            <span>⚠️</span>
            <p>{{ filesError }}</p>
            <button type="button" @click="loadFiles" class="retry-btn-mini">Retry</button>
          </div>

          <!-- Empty -->
          <div v-else-if="files.length === 0" class="empty-mini">
            <span>📂</span>
            <p>{{ searchQuery ? 'No files found' : 'Folder is empty' }}</p>
          </div>

          <!-- Files List -->
          <div v-else class="files-list">
            <div 
              v-for="file in files" 
              :key="file.Path"
              @click="handleFileClick(file)"
              :class="['file-item', { 
                'is-dir': file.IsDir,
                'is-selected': isFileSelected(file)
              }]"
            >
              <span class="file-icon">{{ getFileIcon(file) }}</span>
              <div class="file-info">
                <div class="file-name">{{ file.Name }}</div>
                <div class="file-meta">
                  <span v-if="!file.IsDir" class="file-size">{{ formatFileSize(file.Size) }}</span>
                  <span class="file-time">{{ formatTime(file.ModTime) }}</span>
                </div>
              </div>
              <span v-if="isFileSelected(file)" class="selected-badge">✓</span>
            </div>
          </div>
        </div>

        <!-- Selected File Info -->
        <div v-if="form.source_path" class="selected-file-info">
          <div class="info-label">Selected File:</div>
          <div class="info-value">
            <code>{{ form.source_path }}</code>
            <button type="button" @click="clearSelection" class="clear-selection-btn">
              ✕ Clear
            </button>
          </div>
        </div>
      </div>

      <!-- Step 3: Destination Path -->
      <div v-if="form.source_path" class="form-group">
        <label for="dest">
          <span class="step-number">3</span>
          Destination Path (Local Server) *
        </label>
        <input 
          type="text" 
          id="dest" 
          v-model="form.destination_path" 
          required 
          placeholder="/home/user/restore/"
        />
        <small class="hint">Local path where the file will be restored</small>
      </div>
      
      <!-- Submit Button -->
      <div class="form-actions">
        <button 
          type="submit" 
          :disabled="isLoading || !form.source_path" 
          class="btn-submit btn-restore"
        >
          <span v-if="isLoading">⏳ Starting Restore...</span>
          <span v-else>🔄 Start Restore</span>
        </button>
      </div>
    </form>

    <!-- Messages -->
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
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import jobService from '@/services/jobService';
import driveService from '@/services/driveService';
import monitoringService from '@/services/monitoringService';
import { useRouter } from 'vue-router';

const router = useRouter();
const isLoading = ref(false);
const errorMessage = ref(null);
const message = ref(null);

const remotes = ref([]);
const currentPath = ref('/');
const searchQuery = ref('');
const files = ref([]);
const isLoadingFiles = ref(false);
const filesError = ref(null);

let searchTimeout = null;

const form = ref({
  remote_name: '',
  source_path: '',
  destination_path: ''
});

const pathParts = computed(() => {
  if (currentPath.value === '/') return [];
  return currentPath.value.split('/').filter(p => p);
});

onMounted(async () => {
  await fetchRemotes();
});

async function fetchRemotes() {
  try {
    const data = await monitoringService.getRemotes();
    remotes.value = Array.isArray(data) ? data : [];
  } catch (err) {
    console.error('Failed to fetch remotes:', err);
    errorMessage.value = 'Failed to load remotes list';
  }
}

function handleRemoteSelect() {
  currentPath.value = '/';
  searchQuery.value = '';
  form.value.source_path = '';
  loadFiles();
}

async function loadFiles() {
  if (!form.value.remote_name) return;
  
  isLoadingFiles.value = true;
  filesError.value = null;
  
  try {
    const data = await driveService.listFiles(
      form.value.remote_name,
      currentPath.value,
      searchQuery.value
    );
    
    files.value = data.files || [];
    files.value.sort((a, b) => {
      if (a.IsDir && !b.IsDir) return -1;
      if (!a.IsDir && b.IsDir) return 1;
      return a.Name.localeCompare(b.Name);
    });
    
  } catch (err) {
    console.error('Failed to load files:', err);
    filesError.value = err.response?.data?.error || 'Failed to load files from Google Drive';
  } finally {
    isLoadingFiles.value = false;
  }
}

function handleSearch() {
  clearTimeout(searchTimeout);
  searchTimeout = setTimeout(() => {
    loadFiles();
  }, 500);
}

function clearSearch() {
  searchQuery.value = '';
  loadFiles();
}

function navigateTo(path) {
  currentPath.value = path || '/';
  searchQuery.value = '';
  loadFiles();
}

function getPathUpTo(index) {
  return '/' + pathParts.value.slice(0, index + 1).join('/');
}

function handleFileClick(file) {
  if (file.IsDir) {
    const newPath = currentPath.value === '/' 
      ? '/' + file.Name 
      : currentPath.value + '/' + file.Name;
    navigateTo(newPath);
  } else {
    const fullPath = currentPath.value === '/' 
      ? '/' + file.Name 
      : currentPath.value + '/' + file.Name;
    form.value.source_path = fullPath;
  }
}

function isFileSelected(file) {
  if (file.IsDir) return false;
  const fullPath = currentPath.value === '/' 
    ? '/' + file.Name 
    : currentPath.value + '/' + file.Name;
  return form.value.source_path === fullPath;
}

function clearSelection() {
  form.value.source_path = '';
}

function getFileIcon(file) {
  if (file.IsDir) return '📁';
  const ext = file.Name.split('.').pop().toLowerCase();
  const iconMap = {
    'zip': '📦', 'tar': '📦', 'gz': '📦', 'rar': '📦',
    'sql': '🗄️', 'db': '🗄️',
    'pdf': '📄', 'doc': '📝', 'txt': '📝',
    'jpg': '🖼️', 'png': '🖼️', 'gif': '🖼️',
  };
  return iconMap[ext] || '📄';
}

function formatFileSize(bytes) {
  if (!bytes) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i];
}

function formatTime(timeString) {
  if (!timeString) return '-';
  try {
    const date = new Date(timeString);
    const now = new Date();
    const diff = now - date;
    
    if (diff < 86400000) {
      const hours = Math.floor(diff / 3600000);
      return hours === 0 ? 'Just now' : `${hours}h ago`;
    }
    
    return date.toLocaleDateString('id-ID', { 
      month: 'short',
      day: 'numeric',
      year: 'numeric'
    });
  } catch (e) {
    return timeString;
  }
}

async function handleSubmit() {
  isLoading.value = true;
  errorMessage.value = null;
  message.value = null;

  try {
    const response = await jobService.createRestoreJob(form.value);
    message.value = response.message || 'Restore job started successfully!';
    
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

<style scoped src="@/assets/styles/job-form.css"></style>
<style scoped src="@/assets/styles/file-browser.css"></style>