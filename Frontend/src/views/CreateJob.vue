<template>
  <div class="config-view">
    <form @submit.prevent="handleBackupSubmit" class="config-form">
      <h2>📋 Backup Configuration</h2>
      <p>Buat template Job baru (Manual atau Terjadwal). Logika backup harus disediakan di Pre-Script.</p>

      <!-- Job Name -->
      <div class="form-group">
        <label for="backup-jobName">Job Name (Nama Pekerjaan) *</label>
        <input 
          type="text" 
          id="backup-jobName" 
          v-model="backupForm.job_name" 
          required 
          placeholder="e.g., Daily Database Backup"
        />
      </div>

      <!-- Source Path -->
      <div class="form-group">
        <label for="backup-source">Source Path (Lokal) *</label>
        <input 
          type="text" 
          id="backup-source" 
          v-model="backupForm.source_path" 
          required 
          placeholder="/tmp/backup_file.zip"
        />
        <small class="hint">Path file hasil dari pre-script yang akan di-upload</small>
      </div>

      <!-- Remote Name -->
      <div class="form-group">
        <label for="backup-remote">Remote Name *</label>
        <input 
          type="text" 
          id="backup-remote" 
          v-model="backupForm.remote_name" 
          required 
          placeholder="Gdrive1"
        />
        <small class="hint">Nama remote yang sudah dikonfigurasi di rclone</small>
      </div>

      <!-- Destination Path -->
      <div class="form-group">
        <label for="backup-dest">Destination Path (Cloud) *</label>
        <input 
          type="text" 
          id="backup-dest" 
          v-model="backupForm.destination_path" 
          required 
          placeholder="/backups/database/"
        />
        <small class="hint">Folder tujuan di cloud storage</small>
      </div>

      <!-- ⭐ Schedule Section (Enhanced) -->
      <div class="schedule-section">
        <div class="section-header">
          <h3>⏰ Schedule Configuration</h3>
          <label class="toggle-switch">
            <input 
              type="checkbox" 
              v-model="isScheduled"
              @change="handleScheduleToggle"
            />
            <span class="slider"></span>
            <span class="toggle-label">
              {{ isScheduled ? 'Scheduled Job' : 'Manual Job' }}
            </span>
          </label>
        </div>

        <!-- Schedule Options (Show only if scheduled) -->
        <transition name="slide-fade">
          <div v-if="isScheduled" class="schedule-options">
            <!-- Schedule Type Selector -->
            <div class="schedule-type-selector">
              <button 
                type="button"
                v-for="type in scheduleTypes" 
                :key="type.value"
                @click="selectScheduleType(type.value)"
                :class="['type-btn', { active: scheduleType === type.value }]"
              >
                <span class="type-icon">{{ type.icon }}</span>
                <span class="type-label">{{ type.label }}</span>
              </button>
            </div>

            <!-- Hourly Configuration -->
            <div v-if="scheduleType === 'hourly'" class="schedule-config">
              <label>Every</label>
              <div class="input-group">
                <input 
                  type="number" 
                  v-model.number="scheduleConfig.hours"
                  min="1"
                  max="23"
                  class="time-input"
                />
                <span class="input-suffix">hour(s)</span>
              </div>
            </div>

            <!-- Daily Configuration -->
            <div v-if="scheduleType === 'daily'" class="schedule-config">
              <label>Every day at</label>
              <div class="input-group">
                <input 
                  type="time" 
                  v-model="scheduleConfig.time"
                  class="time-input"
                />
              </div>
            </div>

            <!-- Weekly Configuration -->
            <div v-if="scheduleType === 'weekly'" class="schedule-config">
              <label>Every</label>
              <div class="weekdays-selector">
                <button 
                  type="button"
                  v-for="day in weekdays" 
                  :key="day.value"
                  @click="toggleWeekday(day.value)"
                  :class="['weekday-btn', { active: scheduleConfig.weekdays.includes(day.value) }]"
                >
                  {{ day.short }}
                </button>
              </div>
              <label>at</label>
              <div class="input-group">
                <input 
                  type="time" 
                  v-model="scheduleConfig.time"
                  class="time-input"
                />
              </div>
            </div>

            <!-- Monthly Configuration -->
            <div v-if="scheduleType === 'monthly'" class="schedule-config">
              <label>On day</label>
              <div class="input-group">
                <input 
                  type="number" 
                  v-model.number="scheduleConfig.dayOfMonth"
                  min="1"
                  max="31"
                  class="time-input"
                />
                <span class="input-suffix">of every month</span>
              </div>
              <label>at</label>
              <div class="input-group">
                <input 
                  type="time" 
                  v-model="scheduleConfig.time"
                  class="time-input"
                />
              </div>
            </div>

            <!-- Custom Cron -->
            <div v-if="scheduleType === 'custom'" class="schedule-config">
              <label>Custom Cron Expression</label>
              <input 
                type="text" 
                v-model="scheduleConfig.customCron"
                placeholder="*/5 * * * *"
                class="cron-input"
              />
              <small class="hint">
                Format: minute hour day month weekday
                <a href="https://crontab.guru" target="_blank" rel="noopener">
                  Need help? →
                </a>
              </small>
            </div>

            <!-- Generated Cron Preview -->
            <div class="cron-preview">
              <span class="preview-label">Cron Expression:</span>
              <code class="preview-code">{{ generatedCron || '-' }}</code>
              <span class="preview-description">{{ cronDescription }}</span>
            </div>
          </div>
        </transition>
      </div>

      <!-- Pre Script -->
      <div class="form-group">
        <label for="backup-pre">Pre-Script (Executed BEFORE Rclone)</label>
        <textarea 
          id="backup-pre" 
          v-model="backupForm.pre_script" 
          rows="6" 
          placeholder="#!/bin/bash
# Example: Database dump
mysqldump -u user -p password database > /tmp/backup.sql
gzip /tmp/backup.sql"
        ></textarea>
        <small class="hint">Script untuk generate file backup (e.g., mysqldump, tar, zip)</small>
      </div>

      <!-- Post Script -->
      <div class="form-group">
        <label for="backup-post">Post-Script (Executed AFTER successful upload)</label>
        <textarea 
          id="backup-post" 
          v-model="backupForm.post_script" 
          rows="4" 
          placeholder="#!/bin/bash
# Example: Cleanup
rm /tmp/backup.sql.gz"
        ></textarea>
        <small class="hint">Script untuk cleanup atau notifikasi</small>
      </div>

      <div class="form-actions">
        <button type="button" @click="resetForm" class="btn-secondary">
          Reset
        </button>
        <button type="submit" :disabled="isLoading" class="btn-submit">
          <span v-if="isLoading">⏳ Processing...</span>
          <span v-else>
            {{ isScheduled ? '✅ Create Scheduled Job' : '▶️ Create & Run Manual Job' }}
          </span>
        </button>
      </div>
    </form>

    <hr class="divider" />

    <!-- Restore Form (Simplified) -->
<form @submit.prevent="handleRestoreSubmit" class="config-form restore-form">
      <h2>🔄 Restore Configuration</h2>
      <p>Browse and select file from Google Drive to restore to local server.</p>
      
      <!-- Step 1: Select Remote -->
      <div class="form-group">
        <label for="restore-remote">
          <span class="step-number">1</span>
          Select Remote (Source) *
        </label>
        <select
          id="restore-remote"
          v-model="restoreForm.remote_name"
          @change="handleRemoteSelect"
          required
          class="remote-select"
        >
          <option value="">-- Select Remote --</option>
          <option v-for="remote in remotes" :key="remote.remote_name" :value="remote.remote_name">
            {{ remote.remote_name }}
          </option>
        </select>
      </div>

      <!-- Step 2: Browse & Select File -->
      <div v-if="restoreForm.remote_name" class="form-group">
        <label>
          <span class="step-number">2</span>
          Browse & Select File from Google Drive *
        </label>
        
        <!-- File Browser Section -->
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

          <!-- Loading State -->
          <div v-if="isLoadingFiles" class="loading-mini">
            <div class="spinner-mini"></div>
            <p>Loading files...</p>
          </div>

          <!-- Error State -->
          <div v-else-if="filesError" class="error-mini">
            <span>⚠️</span>
            <p>{{ filesError }}</p>
            <button type="button" @click="loadFiles" class="retry-btn-mini">Retry</button>
          </div>

          <!-- Empty State -->
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
        <div v-if="restoreForm.source_path" class="selected-file-info">
          <div class="info-label">Selected File:</div>
          <div class="info-value">
            <code>{{ restoreForm.source_path }}</code>
            <button type="button" @click="clearSelection" class="clear-selection-btn">
              ✕ Clear
            </button>
          </div>
        </div>
      </div>

      <!-- Step 3: Destination Path -->
      <div v-if="restoreForm.source_path" class="form-group">
        <label for="restore-dest">
          <span class="step-number">3</span>
          Destination Path (Local Server) *
        </label>
        <input 
          type="text" 
          id="restore-dest" 
          v-model="restoreForm.destination_path" 
          required 
          placeholder="/home/user/restore/"
        />
        <small class="hint">Local path where the file will be restored</small>
      </div>
      
      <!-- Submit Button -->
      <div class="form-actions">
        <button 
          type="submit" 
          :disabled="isLoading || !restoreForm.source_path" 
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
import { ref, computed, watch, onMounted } from 'vue';
import jobService from '@/services/jobService';
import driveService from '@/services/driveService';
import monitoringService from '@/services/monitoringService';
import { useRouter } from 'vue-router';

const router = useRouter();
const isLoading = ref(false);
const errorMessage = ref(null);
const message = ref(null);

// Schedule state
const isScheduled = ref(false);
const scheduleType = ref('daily'); // hourly, daily, weekly, monthly, custom


// Restore 
// ⭐ New: File Browser State
const remotes = ref([]);
const currentPath = ref('/');
const searchQuery = ref('');
const files = ref([]);
const isLoadingFiles = ref(false);
const filesError = ref(null);

let searchTimeout = null;


// Computed
const pathParts = computed(() => {
  if (currentPath.value === '/') return [];
  return currentPath.value.split('/').filter(p => p);
});

// Lifecycle
onMounted(async () => {
  await fetchRemotes();
});

// ⭐ Fetch Remotes
async function fetchRemotes() {
  try {
    const data = await monitoringService.getRemoteStatus();
    remotes.value = Array.isArray(data) ? data : [];
  } catch (err) {
    console.error('Failed to fetch remotes:', err);
    errorMessage.value = 'Failed to load remotes list';
  }
}

// ⭐ Handle Remote Selection
function handleRemoteSelect() {
  currentPath.value = '/';
  searchQuery.value = '';
  restoreForm.value.source_path = '';
  if (form.value.remote_name) {
  loadFiles();
 } else {
      // Kosongkan daftar file jika remote dibatalkan/belum dipilih
      files.value = [];
  }
}

// ⭐ Load Files from Google Drive
async function loadFiles() {
  if (!restoreForm.value.remote_name) return;
  
  isLoadingFiles.value = true;
  filesError.value = null;
  
  try {
    const data = await driveService.listFiles(
      restoreForm.value.remote_name,
      currentPath.value,
      searchQuery.value
    );
    
    files.value = data.files || [];
    
    // Sort: directories first, then files
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

// ⭐ Handle Search
function handleSearch() {
  clearTimeout(searchTimeout);
  searchTimeout = setTimeout(() => {
    loadFiles();
  }, 500);
}

// ⭐ Clear Search
function clearSearch() {
  searchQuery.value = '';
  loadFiles();
}

// ⭐ Navigate to Directory
function navigateTo(path) {
  currentPath.value = path || '/';
  searchQuery.value = '';
  loadFiles();
}

// ⭐ Get Path Up To Index
function getPathUpTo(index) {
  return '/' + pathParts.value.slice(0, index + 1).join('/');
}

// ⭐ Handle File Click
function handleFileClick(file) {
  if (file.IsDir) {
    // Navigate into directory
    const newPath = currentPath.value === '/' 
      ? '/' + file.Name 
      : currentPath.value + '/' + file.Name;
    navigateTo(newPath);
  } else {
    // Select file
    const fullPath = currentPath.value === '/' 
      ? '/' + file.Name 
      : currentPath.value + '/' + file.Name;
    restoreForm.value.source_path = fullPath;
  }
}

// ⭐ Check if File is Selected
function isFileSelected(file) {
  if (file.IsDir) return false;
  const fullPath = currentPath.value === '/' 
    ? '/' + file.Name 
    : currentPath.value + '/' + file.Name;
  return restoreForm.value.source_path === fullPath;
}

// ⭐ Clear Selection
function clearSelection() {
  restoreForm.value.source_path = '';
}

// ⭐ Helper Functions
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
    
    if (diff < 86400000) { // Less than 24 hours
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

// Restore Submit
// async function handleRestoreSubmit() {
//   isLoading.value = true;
//   errorMessage.value = null;
//   message.value = null;

//   try {
//     const response = await jobService.createRestoreJob(restoreForm.value);
//     message.value = response.message || 'Restore job started successfully!';
    
//     setTimeout(() => {
//       router.push('/logs');
//     }, 1500);
    
//   } catch (error) {
//     console.error('Create restore job error:', error);
//     errorMessage.value = error.response?.data?.error || 'Failed to start restore.';
//   } finally {
//     isLoading.value = false;
//   }
// }
// Schedule configuration
const scheduleConfig = ref({
  hours: 1,
  time: '00:00',
  weekdays: [],
  dayOfMonth: 1,
  customCron: ''
});

// Schedule types
const scheduleTypes = [
  { value: 'hourly', label: 'Hourly', icon: '🕐' },
  { value: 'daily', label: 'Daily', icon: '📅' },
  { value: 'weekly', label: 'Weekly', icon: '📆' },
  { value: 'monthly', label: 'Monthly', icon: '🗓️' },
  { value: 'custom', label: 'Custom', icon: '⚙️' }
];

// Weekdays
const weekdays = [
  { value: 0, short: 'Sun', full: 'Sunday' },
  { value: 1, short: 'Mon', full: 'Monday' },
  { value: 2, short: 'Tue', full: 'Tuesday' },
  { value: 3, short: 'Wed', full: 'Wednesday' },
  { value: 4, short: 'Thu', full: 'Thursday' },
  { value: 5, short: 'Fri', full: 'Friday' },
  { value: 6, short: 'Sat', full: 'Saturday' }
];

// Form state
const backupForm = ref({
  job_name: '',
  rclone_mode: 'COPY',
  source_path: '',
  remote_name: '',
  destination_path: '',
  schedule_cron: '',
  pre_script: '',
  post_script: ''
});

const restoreForm = ref({
  remote_name: '',
  source_path: '',
  destination_path: ''
});

// Generate cron expression
const generatedCron = computed(() => {
  if (!isScheduled.value) return '';
  
  if (scheduleType.value === 'custom') {
    return scheduleConfig.value.customCron;
  }
  
  if (scheduleType.value === 'hourly') {
    const hours = scheduleConfig.value.hours;
    return `0 */${hours} * * *`;
  }
  
  if (scheduleType.value === 'daily') {
    const [hour, minute] = scheduleConfig.value.time.split(':');
    return `${minute} ${hour} * * *`;
  }
  
  if (scheduleType.value === 'weekly') {
    const [hour, minute] = scheduleConfig.value.time.split(':');
    const days = scheduleConfig.value.weekdays.sort().join(',');
    return days ? `${minute} ${hour} * * ${days}` : '';
  }
  
  if (scheduleType.value === 'monthly') {
    const [hour, minute] = scheduleConfig.value.time.split(':');
    const day = scheduleConfig.value.dayOfMonth;
    return `${minute} ${hour} ${day} * *`;
  }
  
  return '';
});

// Cron description
const cronDescription = computed(() => {
  if (!generatedCron.value) return 'No schedule configured';
  
  if (scheduleType.value === 'hourly') {
    const h = scheduleConfig.value.hours;
    return `Every ${h} hour${h > 1 ? 's' : ''}`;
  }
  
  if (scheduleType.value === 'daily') {
    return `Every day at ${scheduleConfig.value.time}`;
  }
  
  if (scheduleType.value === 'weekly') {
    const days = scheduleConfig.value.weekdays
      .map(d => weekdays.find(w => w.value === d)?.full)
      .join(', ');
    return `Every ${days || 'no days selected'} at ${scheduleConfig.value.time}`;
  }
  
  if (scheduleType.value === 'monthly') {
    const day = scheduleConfig.value.dayOfMonth;
    const suffix = day === 1 ? 'st' : day === 2 ? 'nd' : day === 3 ? 'rd' : 'th';
    return `On the ${day}${suffix} of every month at ${scheduleConfig.value.time}`;
  }
  
  return 'Custom schedule';
});

// Watch generated cron and update form
watch(generatedCron, (newCron) => {
  backupForm.value.schedule_cron = newCron;
});

// Functions
function handleScheduleToggle() {
  if (!isScheduled.value) {
    backupForm.value.schedule_cron = '';
  }
}

function selectScheduleType(type) {
  scheduleType.value = type;
}

function toggleWeekday(day) {
  const index = scheduleConfig.value.weekdays.indexOf(day);
  if (index > -1) {
    scheduleConfig.value.weekdays.splice(index, 1);
  } else {
    scheduleConfig.value.weekdays.push(day);
  }
}

function resetForm() {
  backupForm.value = {
    job_name: '',
    rclone_mode: 'COPY',
    source_path: '',
    remote_name: '',
    destination_path: '',
    schedule_cron: '',
    pre_script: '',
    post_script: ''
  };
  isScheduled.value = false;
  scheduleConfig.value = {
    hours: 1,
    time: '00:00',
    weekdays: [],
    dayOfMonth: 1,
    customCron: ''
  };
}

async function handleBackupSubmit() {
  isLoading.value = true;
  errorMessage.value = null;
  message.value = null;
  
  try {
    const response = await jobService.createBackupJob(backupForm.value);
    message.value = response.message || 'Job created successfully!';
    
    setTimeout(() => {
      if (isScheduled.value) {
        router.push('/scheduled');
      } else {
        router.push('/manual');
      }
    }, 1500);
    
  } catch (error) {
    console.error('Create backup job error:', error);
    errorMessage.value = error.response?.data?.error || 'Failed to create backup job. Please check your input.';
  } finally {
    isLoading.value = false;
  }
}

async function handleRestoreSubmit() {
  isLoading.value = true;
  errorMessage.value = null;
  message.value = null;

  try {
    const response = await jobService.createRestoreJob(restoreForm.value);
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

<style scoped>
/* ⭐ Step Numbers */
.step-number {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  background: linear-gradient(135deg, #3498db 0%, #2980b9 100%);
  color: white;
  border-radius: 50%;
  font-size: 0.85rem;
  font-weight: 700;
  margin-right: 0.5rem;
}

/* ⭐ Remote Select */
.remote-select {
  width: 100%;
  padding: 0.75rem;
  border: 2px solid #e9ecef;
  border-radius: 8px;
  font-size: 1rem;
  cursor: pointer;
  transition: all 0.2s;
  background: white;
}

.remote-select:focus {
  outline: none;
  border-color: #3498db;
  box-shadow: 0 0 0 3px rgba(52, 152, 219, 0.1);
}

/* ⭐ File Browser */
.file-browser {
  border: 2px solid #e9ecef;
  border-radius: 8px;
  overflow: hidden;
  background: #f8f9fa;
  margin-top: 0.5rem;
}

/* Browser Header */
.browser-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem 1rem;
  background: white;
  border-bottom: 1px solid #e9ecef;
}

.breadcrumb {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 0.25rem;
  flex: 1;
}

.breadcrumb-btn {
  background: none;
  border: none;
  color: #3498db;
  cursor: pointer;
  font-weight: 600;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  transition: all 0.2s;
  font-size: 0.9rem;
}

.breadcrumb-btn:hover {
  background: rgba(52, 152, 219, 0.1);
}

.separator {
  color: #6c757d;
  margin: 0 0.25rem;
}

.refresh-btn-small {
  background: #e9ecef;
  border: none;
  width: 32px;
  height: 32px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 1rem;
  transition: all 0.2s;
}

.refresh-btn-small:hover:not(:disabled) {
  background: #dee2e6;
}

.refresh-btn-small:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* Search Box Mini */
.search-box-mini {
  position: relative;
  padding: 0.75rem 1rem;
  background: white;
  border-bottom: 1px solid #e9ecef;
}

.search-icon {
  position: absolute;
  left: 1.75rem;
  top: 50%;
  transform: translateY(-50%);
  font-size: 1rem;
  color: #6c757d;
}

.search-input-mini {
  width: 100%;
  padding: 0.5rem 0.75rem 0.5rem 2.5rem;
  border: 1px solid #e9ecef;
  border-radius: 6px;
  font-size: 0.9rem;
}

.search-input-mini:focus {
  outline: none;
  border-color: #3498db;
}

.clear-btn-mini {
  position: absolute;
  right: 1.75rem;
  top: 50%;
  transform: translateY(-50%);
  background: #e9ecef;
  border: none;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  cursor: pointer;
  font-size: 0.85rem;
}

.clear-btn-mini:hover {
  background: #dee2e6;
}

/* Loading/Error/Empty States */
.loading-mini,
.error-mini,
.empty-mini {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 2rem 1rem;
  gap: 0.75rem;
  background: white;
}

.spinner-mini {
  width: 32px;
  height: 32px;
  border: 3px solid #f3f3f3;
  border-top: 3px solid #3498db;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

.loading-mini p,
.error-mini p,
.empty-mini p {
  margin: 0;
  color: #6c757d;
  font-size: 0.9rem;
}

.error-mini span,
.empty-mini span {
  font-size: 2rem;
  opacity: 0.3;
}

.retry-btn-mini {
  background: #dc3545;
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.85rem;
}

.retry-btn-mini:hover {
  background: #c82333;
}

/* Files List */
.files-list {
  max-height: 400px;
  overflow-y: auto;
  background: white;
}

.files-list::-webkit-scrollbar {
  width: 8px;
}

.files-list::-webkit-scrollbar-track {
  background: #f1f1f1;
}

.files-list::-webkit-scrollbar-thumb {
  background: #cbd5e0;
  border-radius: 4px;
}

/* File Item */
.file-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem 1rem;
  border-bottom: 1px solid #f0f0f0;
  cursor: pointer;
  transition: all 0.2s;
  position: relative;
}

.file-item:hover {
  background: #f8f9fa;
}

.file-item.is-dir {
  background: rgba(52, 152, 219, 0.02);
}

.file-item.is-dir:hover {
  background: rgba(52, 152, 219, 0.05);
}

.file-item.is-selected {
  background: rgba(52, 152, 219, 0.1);
  border-left: 3px solid #3498db;
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
  font-weight: 600;
  color: #2c3e50;
  font-size: 0.9rem;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.file-meta {
  display: flex;
  gap: 1rem;
  font-size: 0.8rem;
  color: #6c757d;
  margin-top: 0.25rem;
}

.selected-badge {
  background: #3498db;
  color: white;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.85rem;
  flex-shrink: 0;
}

/* Selected File Info */
.selected-file-info {
  background: #d1ecf1;
  border: 1px solid #bee5eb;
  border-radius: 8px;
  padding: 1rem;
  margin-top: 1rem;
}

.info-label {
  font-weight: 600;
  color: #0c5460;
  margin-bottom: 0.5rem;
  font-size: 0.9rem;
}

.info-value {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.info-value code {
  flex: 1;
  background: white;
  padding: 0.5rem 0.75rem;
  border-radius: 6px;
  font-family: 'Consolas', monospace;
  font-size: 0.85rem;
  color: #2c3e50;
  word-break: break-all;
}

.clear-selection-btn {
  background: #dc3545;
  color: white;
  border: none;
  padding: 0.5rem 0.75rem;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.8rem;
  font-weight: 600;
  white-space: nowrap;
  transition: all 0.2s;
}

.clear-selection-btn:hover {
  background: #c82333;
}

/* Responsive */
@media (max-width: 768px) {
  .browser-header {
    flex-direction: column;
    gap: 0.5rem;
    align-items: stretch;
  }
  
  .file-meta {
    flex-direction: column;
    gap: 0.25rem;
  }
  
  .info-value {
    flex-direction: column;
    align-items: stretch;
  }
  
  .clear-selection-btn {
    width: 100%;
  }
}

.config-view {
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

.divider {
  margin: 3rem 0;
  border: 0;
  border-top: 2px solid #e9ecef;
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

.form-group textarea {
  font-family: 'Consolas', 'Monaco', monospace;
  font-size: 0.9rem;
  line-height: 1.6;
}

.hint {
  display: block;
  margin-top: 0.5rem;
  font-size: 0.85rem;
  color: #6c757d;
}

.hint a {
  color: #667eea;
  text-decoration: none;
}

.hint a:hover {
  text-decoration: underline;
}

/* Schedule Section */
.schedule-section {
  background: #f8f9fa;
  padding: 1.5rem;
  border-radius: 12px;
  margin-bottom: 1.5rem;
  border: 2px solid #e9ecef;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.section-header h3 {
  margin: 0;
  font-size: 1.1rem;
  font-weight: 600;
  color: #2c3e50;
}

/* Toggle Switch */
.toggle-switch {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  cursor: pointer;
  user-select: none;
}

.toggle-switch input[type="checkbox"] {
  position: relative;
  width: 52px;
  height: 28px;
  appearance: none;
  background: #cbd5e0;
  border-radius: 14px;
  cursor: pointer;
  transition: all 0.3s;
}

.toggle-switch input[type="checkbox"]:checked {
  background: #667eea;
}

.toggle-switch input[type="checkbox"]::before {
  content: '';
  position: absolute;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: white;
  top: 2px;
  left: 2px;
  transition: all 0.3s;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}

.toggle-switch input[type="checkbox"]:checked::before {
  left: 26px;
}

.toggle-label {
  font-weight: 600;
  color: #2c3e50;
  font-size: 0.95rem;
}

/* Schedule Options */
.schedule-options {
  margin-top: 1.5rem;
}

.schedule-type-selector {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(100px, 1fr));
  gap: 0.75rem;
  margin-bottom: 1.5rem;
}

.type-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
  padding: 1rem 0.75rem;
  background: white;
  border: 2px solid #e9ecef;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.type-btn:hover {
  border-color: #667eea;
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(102, 126, 234, 0.15);
}

.type-btn.active {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-color: #667eea;
  color: white;
}

.type-icon {
  font-size: 1.5rem;
}

.type-label {
  font-weight: 600;
  font-size: 0.85rem;
}

/* Schedule Config */
.schedule-config {
  background: white;
  padding: 1.5rem;
  border-radius: 8px;
  margin-bottom: 1rem;
}

.schedule-config label {
  display: block;
  margin-bottom: 0.75rem;
  font-weight: 600;
  color: #2c3e50;
}

.input-group {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-bottom: 1rem;
}

.time-input {
  flex: 0 0 auto;
  padding: 0.75rem;
  border: 2px solid #e9ecef;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 600;
  min-width: 100px;
}

.time-input:focus {
  outline: none;
  border-color: #667eea;
}

.input-suffix {
  color: #6c757d;
  font-weight: 500;
}

/* Weekdays Selector */
.weekdays-selector {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 0.5rem;
  margin-bottom: 1rem;
}

.weekday-btn {
  padding: 0.75rem 0.5rem;
  background: white;
  border: 2px solid #e9ecef;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 600;
  transition: all 0.2s;
}

.weekday-btn:hover {
  border-color: #667eea;
}

.weekday-btn.active {
  background: #667eea;
  border-color: #667eea;
  color: white;
}

/* Cron Input */
.cron-input {
  width: 100%;
  padding: 0.75rem;
  border: 2px solid #e9ecef;
  border-radius: 8px;
  font-family: 'Consolas', monospace;
  font-size: 0.95rem;
}

.cron-input:focus {
  outline: none;
  border-color: #667eea;
}

/* Cron Preview */
.cron-preview {
  background: #2c3e50;
  padding: 1rem 1.25rem;
  border-radius: 8px;
  display: flex;
  align-items: center;
  gap: 1rem;
  flex-wrap: wrap;
}

.preview-label {
  color: #a0aec0;
  font-weight: 600;
  font-size: 0.85rem;
}

.preview-code {
  background: #1a202c;
  color: #48bb78;
  padding: 0.5rem 1rem;
  border-radius: 6px;
  font-family: 'Consolas', monospace;
  font-size: 0.95rem;
  font-weight: 600;
}

.preview-description {
  color: #cbd5e0;
  font-size: 0.9rem;
  font-style: italic;
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

.btn-secondary {
  background: white;
  color: #6c757d;
  border: 2px solid #e9ecef;
  padding: 0.75rem 1.5rem;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 600;
  transition: all 0.2s;
}

.btn-secondary:hover {
  border-color: #cbd5e0;
  background: #f8f9fa;
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

.btn-submit:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.4);
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
.slide-fade-enter-active {
  transition: all 0.3s ease-out;
}

.slide-fade-leave-active {
  transition: all 0.2s cubic-bezier(1, 0.5, 0.8, 1);
}

.slide-fade-enter-from,
.slide-fade-leave-to {
  transform: translateY(-10px);
  opacity: 0;
}

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
  .config-view {
    padding: 1rem;
  }
  
  .config-form {
    padding: 1.5rem;
  }
  
  .schedule-type-selector {
    grid-template-columns: repeat(auto-fit, minmax(80px, 1fr));
  }
  
  .weekdays-selector {
    gap: 0.35rem;
  }
  
  .weekday-btn {
    padding: 0.5rem 0.25rem;
    font-size: 0.85rem;
  }
  
  .form-actions {
    flex-direction: column;
  }
  
  .btn-submit,
  .btn-secondary {
    width: 100%;
  }
}
</style>