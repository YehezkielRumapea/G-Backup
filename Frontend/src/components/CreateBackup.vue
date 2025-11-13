<template>
  <div class="config-view">
      <div v-if="isVisible" class="modal-overlay" @click.self="close">
      <div class="modal-content">
      <div class="modal-header">
          <h2>📋 Backup Configuration</h2>
          <button class="close-btn" @click="close">✖</button>
        </div>  
      <form @submit.prevent="handleBackupSubmit" class="config-form">
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

      <!-- Schedule Section -->
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

        <!-- Schedule Options -->
        <transition name="slide-fade">
          <div v-if="isScheduled" class="schedule-options">
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

            <!-- HOURLY -->
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

            <!-- DAILY -->
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

            <!-- WEEKLY -->
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

            <!-- MONTHLY -->
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

            <!-- CUSTOM -->
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

            <!-- CRON PREVIEW -->
            <div class="cron-preview">
              <span class="preview-label">Cron Expression:</span>
              <code class="preview-code">{{ generatedCron || '-' }}</code>
              <span class="preview-description">{{ cronDescription }}</span>
            </div>
          </div>
        </transition>
      </div>

      <!-- Pre-Script -->
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

      <!-- Post-Script -->
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

      <!-- Buttons -->
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
  </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue';
import { useRouter } from 'vue-router';
import jobService from '@/services/jobService';

const props = defineProps({
  isVisible: Boolean
})
const emit = defineEmits(['close', 'success'])

// Router
const router = useRouter();

// Reactive states
const isLoading = ref(false);
const message = ref(null);
const errorMessage = ref(null);

// Schedule
const isScheduled = ref(false);
const scheduleType = ref('daily');
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

const weekdays = [
  { value: 0, short: 'Sun', full: 'Sunday' },
  { value: 1, short: 'Mon', full: 'Monday' },
  { value: 2, short: 'Tue', full: 'Tuesday' },
  { value: 3, short: 'Wed', full: 'Wednesday' },
  { value: 4, short: 'Thu', full: 'Thursday' },
  { value: 5, short: 'Fri', full: 'Friday' },
  { value: 6, short: 'Sat', full: 'Saturday' }
];

// Form data
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

// Computed: Cron Expression
const generatedCron = computed(() => {
  if (!isScheduled.value) return '';
  const cfg = scheduleConfig.value;

  switch (scheduleType.value) {
    case 'hourly':
      return `0 */${cfg.hours} * * *`;
    case 'daily': {
      const [h, m] = cfg.time.split(':');
      return `${m} ${h} * * *`;
    }
    case 'weekly': {
      const [h, m] = cfg.time.split(':');
      const days = cfg.weekdays.sort().join(',');
      return days ? `${m} ${h} * * ${days}` : '';
    }
    case 'monthly': {
      const [h, m] = cfg.time.split(':');
      return `${m} ${h} ${cfg.dayOfMonth} * *`;
    }
    case 'custom':
      return cfg.customCron;
    default:
      return '';
  }
});

// Computed: Cron Description
const cronDescription = computed(() => {
  if (!generatedCron.value) return 'No schedule configured';
  const cfg = scheduleConfig.value;

  switch (scheduleType.value) {
    case 'hourly':
      return `Every ${cfg.hours} hour${cfg.hours > 1 ? 's' : ''}`;
    case 'daily':
      return `Every day at ${cfg.time}`;
    case 'weekly': {
      const days = cfg.weekdays
        .map(d => weekdays.find(w => w.value === d)?.full)
        .join(', ');
      return `Every ${days || 'no days selected'} at ${cfg.time}`;
    }
    case 'monthly':
      return `On day ${cfg.dayOfMonth} of every month at ${cfg.time}`;
    case 'custom':
      return `Custom cron: ${cfg.customCron}`;
    default:
      return 'No schedule configured';
  }
});

// Watch cron to update form
watch(generatedCron, (newCron) => {
  backupForm.value.schedule_cron = newCron;
});

// Functions
function handleScheduleToggle() {
  if (!isScheduled.value) backupForm.value.schedule_cron = '';
}

function selectScheduleType(type) {
  scheduleType.value = type;
}

function toggleWeekday(day) {
  const days = scheduleConfig.value.weekdays;
  const idx = days.indexOf(day);
  if (idx > -1) days.splice(idx, 1);
  else days.push(day);
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
  isLoading.value = true
  message.value = null
  errorMessage.value = null
  try {
    const res = await jobService.createBackupJob(backupForm.value)
    message.value = res.message || 'Job created successfully!'
    emit('success')
    setTimeout(() => emit('close'), 1500)
  } catch (err) {
    errorMessage.value = err.response?.data?.error || 'Failed to create backup job.'
  } finally {
    isLoading.value = false
  }
}

function close() {
  emit('close')
}
</script>

<style scoped>
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

.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

/* Modal box */
.modal-content {
  background: #fff;
  border-radius: 12px;
  width: 90%;
  max-width: 800px;
  max-height: 90vh;
  overflow-y: auto;
  padding: 2rem;
  position: relative;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.2);
}

/* Header */
.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.modal-header h2 {
  margin: 0;
  font-size: 1.4rem;
  font-weight: 700;
  color: #2c3e50;
}

.close-btn {
  background: transparent;
  border: none;
  font-size: 1.3rem;
  cursor: pointer;
  color: #888;
}

.close-btn:hover {
  color: #000;
}

/* Transisi */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>