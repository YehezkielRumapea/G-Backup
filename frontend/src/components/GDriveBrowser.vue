<template>
  <div class="gdrive-browser">
    <!-- Header -->
    <div class="browser-header">
      <h3>📁 Browse {{ remoteName }}</h3>
      <button v-if="currentPath !== '/'" @click="goBack" class="btn-back">
        ← Back
      </button>
    </div>

    <!-- Breadcrumb Navigation -->
    <div class="breadcrumb">
      <button 
        @click="navigateTo('/')" 
        :class="{ active: currentPath === '/' }"
        class="breadcrumb-item"
      >
        Root
      </button>
      <template v-for="(segment, index) in pathSegments" :key="index">
        <span class="breadcrumb-separator">/</span>
        <button 
          @click="navigateTo(buildPath(index))"
          :class="{ active: buildPath(index) === currentPath }"
          class="breadcrumb-item"
        >
          {{ segment }}
        </button>
      </template>
    </div>

    <!-- Loading State -->
    <div v-if="isLoading" class="loading">
      <div class="spinner"></div>
      <p>Loading files...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="error-message">
      <p>❌ {{ error }}</p>
      <button @click="loadFiles" class="btn-retry">Try Again</button>
    </div>

    <!-- Empty State -->
    <div v-else-if="files.length === 0" class="empty-state">
      <p>📭 No files or folders found</p>
    </div>

    <!-- Files List -->
    <div v-else class="files-container">
      <!-- Folder/File Items -->
      <div 
        v-for="file in files" 
        :key="file.path"
        @click="handleFileClick(file)"
        :class="['file-item', { 'is-dir': file.is_dir }]"
      >
        <!-- Icon -->
        <div class="file-icon">
          <span v-if="file.is_dir" class="icon">📁</span>
          <span v-else :class="getFileIcon(file.name)" class="icon">{{ getFileIcon(file.name) }}</span>
        </div>

        <!-- File Info -->
        <div class="file-info">
          <div class="file-name">{{ file.name }}</div>
          <div v-if="!file.is_dir" class="file-details">
            {{ formatFileSize(file.size) }} • {{ formatDate(file.mod_time) }}
          </div>
          <div v-else class="file-details">
            Folder • {{ formatDate(file.mod_time) }}
          </div>
        </div>

        <!-- Arrow for folders -->
        <div v-if="file.is_dir" class="file-arrow">
          →
        </div>
      </div>

      <!-- Total Size Info -->
      <div class="total-info">
        <p>📊 Total size: {{ formatFileSize(totalSize) }}</p>
        <p>📦 {{ files.length }} items</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import driveService from '@/services/driveService';

// Props
const props = defineProps({
  remoteName: {
    type: String,
    default: 'Gdrive1'
  },
  initialPath: {
    type: String,
    default: '/'
  }
});

// Emits
const emit = defineEmits(['select-file', 'navigate']);

// State
const currentPath = ref('/');
const files = ref([]);
const isLoading = ref(false);
const error = ref(null);
const totalSize = ref(0);

// Computed
const pathSegments = computed(() => {
  if (currentPath.value === '/') return [];
  return currentPath.value
    .split('/')
    .filter(s => s !== '')
    .map(s => s);
});

// ============================================
// Methods
// ============================================

// Load files dari API
async function loadFiles() {
  isLoading.value = true;
  error.value = null;
  files.value = [];
  totalSize.value = 0;

  try {
    console.log(`🔄 Loading files from ${props.remoteName}:${currentPath.value}`);
    
    // Call service untuk browse files
    const response = await driveService.browseFiles(props.remoteName, currentPath.value);
    
    files.value = response.files || [];
    totalSize.value = response.total_size || 0;

    console.log(`✅ Loaded ${files.value.length} items`);
    
    // Emit navigation event
    emit('navigate', {
      remote: props.remoteName,
      path: currentPath.value,
      filesCount: files.value.length
    });

  } catch (err) {
    error.value = err.response?.data?.error || err.message || 'Failed to load files';
    console.error('❌ Browser error:', error.value);
  } finally {
    isLoading.value = false;
  }
}

// Handle file/folder click
function handleFileClick(file) {
  if (file.is_dir) {
    // Navigate ke folder
    navigateTo(file.path);
  } else {
    // Emit select-file event untuk file
    emit('select-file', {
      name: file.name,
      path: file.path,
      size: file.size,
      remote: props.remoteName
    });
  }
}

// Navigate ke path tertentu
function navigateTo(path) {
  currentPath.value = path;
  loadFiles();
}

// Go back ke parent folder
function goBack() {
  const parts = currentPath.value
    .split('/')
    .filter(s => s !== '')
    .slice(0, -1);
  
  const parentPath = parts.length === 0 ? '/' : '/' + parts.join('/');
  navigateTo(parentPath);
}

// Build path dari breadcrumb index
function buildPath(index) {
  const parts = pathSegments.value.slice(0, index + 1);
  return '/' + parts.join('/');
}

// ============================================
// Utility Functions
// ============================================

// Get file icon berdasarkan extension
function getFileIcon(fileName) {
  const ext = fileName.split('.').pop().toLowerCase();
  
  const iconMap = {
    // Archives
    'zip': '🗜️',
    'rar': '🗜️',
    '7z': '🗜️',
    'tar': '🗜️',
    'gz': '🗜️',
    
    // Documents
    'pdf': '📄',
    'doc': '📝',
    'docx': '📝',
    'txt': '📄',
    'xlsx': '📊',
    'csv': '📊',
    'xls': '📊',
    'ppt': '🎞️',
    'pptx': '🎞️',
    
    // Images
    'jpg': '🖼️',
    'jpeg': '🖼️',
    'png': '🖼️',
    'gif': '🖼️',
    'svg': '🖼️',
    'webp': '🖼️',
    
    // Videos
    'mp4': '🎬',
    'avi': '🎬',
    'mkv': '🎬',
    'mov': '🎬',
    'flv': '🎬',
    
    // Audio
    'mp3': '🎵',
    'wav': '🎵',
    'flac': '🎵',
    'm4a': '🎵',
    
    // Code
    'js': '⚙️',
    'py': '⚙️',
    'go': '⚙️',
    'java': '⚙️',
    'cpp': '⚙️',
    'json': '⚙️',
    'xml': '⚙️',
    'html': '⚙️',
    'css': '⚙️',
    
    // Database
    'sql': '🗄️',
    'db': '🗄️',
    'sqlite': '🗄️',
  };
  
  return iconMap[ext] || '📄';
}

// Format file size ke human readable
function formatFileSize(bytes) {
  if (bytes === 0) return '0 B';
  
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

// Format date
function formatDate(dateString) {
  if (!dateString) return 'N/A';
  
  try {
    const date = new Date(dateString);
    return date.toLocaleDateString('id-ID', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  } catch (e) {
    return dateString;
  }
}

// ============================================
// Lifecycle
// ============================================

onMounted(() => {
  currentPath.value = props.initialPath;
  loadFiles();
});
</script>

<style scoped>
.gdrive-browser {
  background: #fff;
  border: 1px solid #e5e5e5;
  border-radius: 8px;
  padding: 1.5rem;
  max-height: 600px;
  display: flex;
  flex-direction: column;
}

/* Header */
.browser-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
  padding-bottom: 1rem;
  border-bottom: 2px solid #f0f0f0;
}

.browser-header h3 {
  margin: 0;
  font-size: 1.125rem;
  font-weight: 600;
  color: #1a1a1a;
}

.btn-back {
  background: #f5f5f5;
  border: 1px solid #e5e5e5;
  padding: 0.5rem 1rem;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.875rem;
  font-weight: 500;
  transition: all 0.2s;
}

.btn-back:hover {
  background: #e5e5e5;
  border-color: #d4d4d4;
}

/* Breadcrumb */
.breadcrumb {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 1rem;
  overflow-x: auto;
  padding: 0.5rem 0;
  font-size: 0.875rem;
}

.breadcrumb-item {
  background: transparent;
  border: none;
  color: #3b82f6;
  cursor: pointer;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  transition: all 0.2s;
  white-space: nowrap;
}

.breadcrumb-item:hover {
  background: #f0f0f0;
}

.breadcrumb-item.active {
  color: #1a1a1a;
  font-weight: 600;
  background: #f0f0f0;
}

.breadcrumb-separator {
  color: #999;
  margin: 0 0.25rem;
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
  width: 40px;
  height: 40px;
  border: 3px solid #e5e5e5;
  border-top-color: #3b82f6;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.loading p {
  margin: 0;
  color: #666;
  font-size: 0.9375rem;
}

/* Error */
.error-message {
  background: #fee2e2;
  border: 1px solid #fecaca;
  color: #991b1b;
  padding: 1rem;
  border-radius: 6px;
  text-align: center;
}

.error-message p {
  margin: 0 0 0.5rem 0;
  font-weight: 500;
}

.btn-retry {
  background: #991b1b;
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.875rem;
  font-weight: 500;
}

.btn-retry:hover {
  background: #7f1d1d;
}

/* Empty State */
.empty-state {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem;
  color: #999;
  font-size: 0.9375rem;
}

/* Files Container */
.files-container {
  flex: 1;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.files-container::-webkit-scrollbar {
  width: 6px;
}

.files-container::-webkit-scrollbar-track {
  background: transparent;
}

.files-container::-webkit-scrollbar-thumb {
  background: #d4d4d4;
  border-radius: 3px;
}

.files-container::-webkit-scrollbar-thumb:hover {
  background: #999;
}

/* File Item */
.file-item {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 0.75rem;
  border: 1px solid #f0f0f0;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
}

.file-item:hover {
  background: #f9f9f9;
  border-color: #e5e5e5;
}

.file-item.is-dir {
  background: #f9fafb;
}

.file-item.is-dir:hover {
  background: #f0f0f0;
  border-color: #3b82f6;
}

.file-icon {
  font-size: 1.5rem;
  flex-shrink: 0;
}

.icon {
  display: block;
}

.file-info {
  flex: 1;
  min-width: 0;
}

.file-name {
  font-weight: 500;
  color: #1a1a1a;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  font-size: 0.9375rem;
}

.file-details {
  font-size: 0.8125rem;
  color: #999;
  margin-top: 0.25rem;
}

.file-arrow {
  color: #3b82f6;
  font-size: 1.25rem;
  font-weight: bold;
  flex-shrink: 0;
}

/* Total Info */
.total-info {
  padding: 1rem 0;
  border-top: 1px solid #f0f0f0;
  margin-top: auto;
  font-size: 0.8125rem;
  color: #666;
}

.total-info p {
  margin: 0.25rem 0;
}

/* Responsive */
@media (max-width: 768px) {
  .gdrive-browser {
    max-height: 400px;
    padding: 1rem;
  }

  .browser-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.5rem;
  }

  .breadcrumb {
    font-size: 0.8125rem;
  }

  .file-item {
    gap: 0.75rem;
  }
}
</style>