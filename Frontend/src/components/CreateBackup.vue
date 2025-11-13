<template>
  <div class="config-view">
    <form @submit.prevent="handleBackupSubmit" class="config-form">
      <h2>📋 Backup Configuration</h2>
      <p>Buat template Job baru (Manual atau Terjadwal). Logika backup harus disediakan di Pre-Script.</p>

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

            <div class="cron-preview">
              <span class="preview-label">Cron Expression:</span>
              <code class="preview-code">{{ generatedCron || '-' }}</code>
              <span class="preview-description">{{ cronDescription }}</span>
            </div>
          </div>
        </transition>
      </div>

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
import { ref, computed, watch } from 'vue';
import jobService from '@/services/jobService';
import { useRouter } from 'vue-router';

const router = useRouter();
const isLoading = ref(false);
const errorMessage = ref(null);
const message = ref(null);

// Schedule state
const isScheduled = ref(false);
const scheduleType = ref('daily'); 

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

// Generate cron expression (Logic from original file)
const generatedCron = computed(() => {
  if (!isScheduled.value) return '';
  // ... (Cron generation logic remains the same)
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

// Cron description (Logic from original file)
const cronDescription = computed(() => {
  if (!generatedCron.value) return 'No schedule configured';
  // ... (Cron description logic remains the same)
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
</style>