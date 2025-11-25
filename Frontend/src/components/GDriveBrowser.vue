<template>
  <div class="gdrive-browser">
    <!-- Header -->
    <div class="browser-header">
      <h3>📁 Browse {{ remoteName }}</h3>
      <div class="header-actions">
        <button v-if="currentPath !== '/'" @click="goBack" class="btn-back">
          ← Back
        </button>
      </div>
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

    <!-- Selected Item Info (Display Only) -->
    <div v-if="selectedItem" class="selected-banner">
      <div class="selected-icon">
        {{ selectedItem.is_dir ? '📁' : getFileIcon(selectedItem.name) }}
      </div>
      <div class="selected-content">
        <strong>✅ Selected:</strong> {{ selectedItem.name }}
        <br>
        <small>{{ selectedItem.is_dir ? 'Folder' : formatFileSize(selectedItem.size) }} • {{ currentPath }}</small>
      </div>
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
        :class="['file-item', { 
          'is-dir': file.is_dir,
          'is-selected': selectedItem?.path === file.path
        }]"
      >
        <!-- Folder: Double-click to enter, Click to select -->
        <template v-if="file.is_dir">
          <div 
            class="file-content"
            @click="selectItem(file)"
            @dblclick="navigateTo(file.path)"
          >
            <div class="file-icon">📁</div>
            <div class="file-info">
              <div class="file-name">{{ file.name }}</div>
              <div class="file-details">
                Folder • {{ formatDate(file.mod_time) }}
              </div>
            </div>
          </div>
          <div class="file-badge">🗂️ Folder</div>
        </template>

        <!-- File: Click to select -->
        <template v-else>
          <div 
            class="file-content"
            @click="selectItem(file)"
          >
            <div class="file-icon">{{ getFileIcon(file.name) }}</div>
            <div class="file-info">
              <div class="file-name">{{ file.name }}</div>
              <div class="file-details">
                {{ formatFileSize(file.size) }} • {{ formatDate(file.mod_time) }}
              </div>
            </div>
          </div>
          <div class="file-badge">📄 File</div>
        </template>
      </div>

      <!-- Total Size Info -->
      <div class="total-info">
        <p>📊 Total size: {{ formatFileSize(totalSize) }}</p>
        <p>📦 {{ files.length }} items</p>
      </div>
    </div>

    <!-- Instructions -->
    <div class="instructions">
      <small>
        💡 <strong>Folder:</strong> Click to select, Double-click to open<br>
        💡 <strong>File:</strong> Click to select<br>
        💡 Selection otomatis disimpan
      </small>
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
const selectedItem = ref(null);

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
  selectedItem.value = null;

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

// Select item dan langsung emit (tanpa perlu confirm button)
function selectItem(file) {
  selectedItem.value = file;
  console.log(`✅ Selected: ${file.name} (${file.is_dir ? 'Folder' : 'File'})`);
  
  // Auto-emit langsung
  emit('select-file', {
    name: file.name,
    path: file.path,
    size: file.size,
    is_dir: file.is_dir,
    remote: props.remoteName
  });
  
  console.log(`🎯 Selection emitted: ${file.name}`);
}

// Navigate ke path tertentu (untuk double-click folder)
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
    'zip': '🗜️', 'rar': '🗜️', '7z': '🗜️', 'tar': '🗜️', 'gz': '🗜️',
    'pdf': '📄', 'doc': '📝', 'docx': '📝', 'txt': '📄', 
    'xlsx': '📊', 'csv': '📊', 'xls': '📊', 'ppt': '🎞️', 'pptx': '🎞️',
    'jpg': '🖼️', 'jpeg': '🖼️', 'png': '🖼️', 'gif': '🖼️', 'svg': '🖼️', 'webp': '🖼️',
    'mp4': '🎬', 'avi': '🎬', 'mkv': '🎬', 'mov': '🎬', 'flv': '🎬',
    'mp3': '🎵', 'wav': '🎵', 'flac': '🎵', 'm4a': '🎵',
    'js': '⚙️', 'py': '⚙️', 'go': '⚙️', 'java': '⚙️', 'cpp': '⚙️', 
    'json': '⚙️', 'xml': '⚙️', 'html': '⚙️', 'css': '⚙️',
    'sql': '🗄️', 'db': '🗄️', 'sqlite': '🗄️',
  };
  
  return iconMap[ext] || '📄';
}

// Format file size ke human readable
function formatFileSize(bytes) {
  if (!bytes || bytes === 0) return '0 B';
  
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
  gap: 1rem;
}

/* Header */
.browser-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 1rem;
  border-bottom: 2px solid #f0f0f0;
}

.browser-header h3 {
  margin: 0;
  font-size: 1.125rem;
  font-weight: 600;
  color: #1a1a1a;
}

.header-actions {
  display: flex;
  gap: 0.5rem;
}

.btn-back {
  padding: 0.5rem 1rem;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.875rem;
  font-weight: 500;
  transition: all 0.2s;
  border: 1px solid #e5e5e5;
  background: #f5f5f5;
  color: #666;
}

.btn-back:hover {
  background: #e5e5e5;
  border-color: #d4d4d4;
}

/* Selected Banner */
.selected-banner {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem;
  background: #f0fdf4;
  border: 1px solid #86efac;
  border-radius: 6px;
}

.selected-icon {
  font-size: 1.5rem;
  flex-shrink: 0;
}

.selected-content {
  flex: 1;
  font-size: 0.875rem;
  color: #166534;
}

.selected-content strong {
  color: #15803d;
}

/* Breadcrumb */
.breadcrumb {
  display: flex;
  align-items: center;
  gap: 0.5rem;
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

/* Error */
.error-message {
  background: #fee2e2;
  border: 1px solid #fecaca;
  color: #991b1b;
  padding: 1rem;
  border-radius: 6px;
  text-align: center;
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
  margin-top: 0.5rem;
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

.files-container::-webkit-scrollbar-thumb {
  background: #d4d4d4;
  border-radius: 3px;
}

/* File Item */
.file-item {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 0.75rem;
  border: 2px solid #f0f0f0;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
  background: #fff;
}

.file-item:hover {
  background: #f9f9f9;
  border-color: #e5e5e5;
}

.file-item.is-selected {
  background: #f0fdf4;
  border-color: #86efac;
}

.file-item.is-dir:hover {
  border-color: #3b82f6;
}

.file-content {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 0.75rem;
  min-width: 0;
}

.file-icon {
  font-size: 1.5rem;
  flex-shrink: 0;
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

.file-badge {
  font-size: 0.75rem;
  padding: 0.25rem 0.5rem;
  background: #f0f0f0;
  border-radius: 4px;
  color: #666;
  flex-shrink: 0;
  white-space: nowrap;
}

/* Total Info */
.total-info {
  padding: 1rem 0;
  border-top: 1px solid #f0f0f0;
  font-size: 0.8125rem;
  color: #666;
}

.total-info p {
  margin: 0.25rem 0;
}

/* Instructions */
.instructions {
  background: #f0f8ff;
  border: 1px solid #b3d9ff;
  border-radius: 6px;
  padding: 0.75rem 1rem;
  font-size: 0.75rem;
  color: #0c4a6e;
}

.instructions small {
  line-height: 1.6;
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

  .header-actions {
    width: 100%;
  }

  .btn-back {
    width: 100%;
  }

  .file-item {
    gap: 0.5rem;
  }
}
</style>