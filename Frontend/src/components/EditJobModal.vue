<template>
  <div class="edit-job-view">
    <div v-if="isOpen" class="modal-overlay" @click.self="handleClose">
      <div class="modal-content">
        <div class="modal-header">
          <h2>Edit Job: {{ jobData?.job_name || `#${jobId}` }}</h2>
          <button class="close-btn" @click="handleClose">×</button>
        </div>  
        
        <!-- Loading State -->
        <div v-if="isLoading" class="loading-container">
          <div class="spinner"></div>
          <p>Loading job data...</p>
        </div>

        <!-- Error State -->
        <div v-else-if="loadError" class="error-banner">
          <p>{{ loadError }}</p>
          <button @click="loadJobData" class="retry-btn">Try Again</button>
        </div>

        <!-- Form -->
        <form v-else @submit.prevent="handleSubmit" class="config-form">
          <p class="form-description">Edit job configuration. Operation mode cannot be changed.</p>


          <!-- Job Name -->
          <div class="form-group">
            <label for="edit-jobName">Job Name *</label>
            <input 
              type="text" 
              id="edit-jobName" 
              v-model="formData.job_name" 
              required 
              placeholder="e.g., Daily Database Backup"
            />
          </div>

          <!-- Source Path -->
          <div class="form-group">
            <label for="edit-source">Source Path (Lokal) *</label>
            <input 
              type="text" 
              id="edit-source" 
              v-model="formData.source_path" 
              required 
              placeholder="/tmp/backup_file.zip"
            />
            <small class="hint">Path file hasil dari pre-script yang akan di-upload</small>
          </div>

          <!-- Target Storage -->
          <div>
            <label for="gdrive_target">Pilih GDrive Target:</label>
            <select v-model="formData.gdrive_target" required>
              <option value="" disabled>Pilih Storage</option>
              <option 
                v-for="remote in availableRemotes" 
                :key="remote.remote_name" 
                :value="remote.remote_name"
              >
                {{ remote.remote_name }}
              </option>
            </select>
          </div>


          <!-- Destination Path -->
          <div class="form-group">
            <label for="edit-dest">Destination Path (Cloud) *</label>
            <input 
              type="text" 
              id="edit-dest" 
              v-model="formData.destination_path" 
              required 
              placeholder="/backups/database/"
            />
            <small class="hint">Folder tujuan di cloud storage</small>
          </div>

          <!-- Schedule Section -->
          <div class="schedule-section">
            <div class="section-header">
              <h3>Schedule Configuration</h3>
              <label class="toggle-switch">
                <input 
                  type="checkbox" 
                  v-model="isScheduled"
                  @change="handleScheduleToggle"
                />
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
            <label for="edit-pre">Pre-Script (Executed BEFORE Rclone)</label>
            <textarea 
              id="edit-pre" 
              v-model="formData.pre_script" 
              rows="1" 
              placeholder="#!/bin/bash
# Example: Database dump
mysqldump -u user -p password database > /tmp/backup.sql
gzip /tmp/backup.sql"
            ></textarea>
            <small class="hint">Script untuk generate file backup (e.g., mysqldump, tar, zip)</small>
          </div>

          <!-- Post-Script -->
          <div class="form-group">
            <label for="edit-post">Post-Script (Executed AFTER successful upload)</label>
            <textarea 
              id="edit-post" 
              v-model="formData.post_script" 
              rows="1" 
              placeholder="#!/bin/bash
# Example: Cleanup
rm /tmp/backup.sql.gz"
            ></textarea>
            <small class="hint">Script untuk cleanup atau notifikasi</small>
          </div>

          <!-- Changes Summary -->
          <div v-if="changedFieldsCount > 0" class="changes-summary">
            <p class="summary-title">📝 Changes detected ({{ changedFieldsCount }} fields):</p>
            <ul class="changes-list">
              <li v-for="field in changedFieldsList" :key="field">
                <span class="field-name">{{ formatFieldName(field) }}</span>
                <span class="field-change">{{ originalData[field] }} → {{ formData[field] }}</span>
              </li>
            </ul>
          </div>

          <!-- Error Message -->
          <div v-if="submitError" class="error-banner">
            <p>{{ submitError }}</p>
          </div>

          <!-- Buttons -->
          <div class="form-actions">
            <button type="button" @click="handleClose" :disabled="isSaving" class="btn-secondary">
              Cancel
            </button>
            <button 
              type="submit" 
              :disabled="isSaving || changedFieldsCount === 0" 
              class="btn-submit"
            >
              <span v-if="isSaving">Saving...</span>
              <span v-else>Save Changes ({{ changedFieldsCount }})</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue';
import jobService from '@/services/jobService';

const props = defineProps({
  jobId: {
    type: Number,
    required: true,
  },
  isOpen: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(['close', 'success']);

// State
const isLoading = ref(false);
const isSaving = ref(false);
const loadError = ref(null);
const submitError = ref(null);
const jobData = ref(null);
const originalData = ref(null);

const isScheduled = ref(false);
const scheduleType = ref('daily');
const scheduleConfig = ref({
  hours: 1,
  time: '00:00',
  weekdays: [],
  dayOfMonth: 1,
  customCron: ''
});

const scheduleTypes = [
  { value: 'hourly', label: 'Hourly' },
  { value: 'daily', label: 'Daily' },
  { value: 'weekly', label: 'Weekly' },
  { value: 'monthly', label: 'Monthly' },
  { value: 'custom', label: 'Custom' }
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

const formData = ref({
  job_name: '',
  source_path: '',
  remote_name: '',
  destination_path: '',
  schedule_cron: '',
  pre_script: '',
  post_script: ''
});

// Watch modal open
watch(() => props.isOpen, (newValue) => {
  if (newValue && props.jobId) {
    loadJobData();
  }
});

// Load job data
const loadJobData = async () => {
  try {
    isLoading.value = true;
    loadError.value = null;

    const data = await jobService.getJobById(props.jobId);
    jobData.value = data;

    // Store original for comparison
originalData.value = {
  job_name: data.JobName || '',
  source_path: data.SourcePath || '',
  remote_name: data.RemoteName || '',
  destination_path: data.DestinationPath || '',
  schedule_cron: data.ScheduleCron || '',
  pre_script: data.PreScript || '',
  post_script: data.PostScript || ''
};

// Set formData
formData.value = { ...originalData.value };

// Parse schedule
parseScheduleFromCron(data.ScheduleCron || '');

    isLoading.value = false;
  } catch (err) {
    loadError.value = err.response?.data?.error || err.message || 'Failed to load job data';
    console.error('Load job error:', err);
    isLoading.value = false;
  }
};

// Parse cron expression back to UI
const parseScheduleFromCron = (cron) => {
  if (!cron || cron.trim() === '') {
    isScheduled.value = false;
    return;
  }

  isScheduled.value = true;
  const parts = cron.trim().split(' ');

  if (parts.length !== 5) {
    scheduleType.value = 'custom';
    scheduleConfig.value.customCron = cron;
    return;
  }

  const [minute, hour, dayOfMonth, month, dayOfWeek] = parts;

  // Try to detect schedule type
  if (hour.startsWith('*/')) {
    scheduleType.value = 'hourly';
    scheduleConfig.value.hours = parseInt(hour.substring(2)) || 1;
  } else if (dayOfMonth === '*' && month === '*' && dayOfWeek === '*') {
    scheduleType.value = 'daily';
    scheduleConfig.value.time = `${hour.padStart(2, '0')}:${minute.padStart(2, '0')}`;
  } else if (dayOfMonth === '*' && month === '*' && dayOfWeek !== '*') {
    scheduleType.value = 'weekly';
    scheduleConfig.value.weekdays = dayOfWeek.split(',').map(d => parseInt(d));
    scheduleConfig.value.time = `${hour.padStart(2, '0')}:${minute.padStart(2, '0')}`;
  } else if (month === '*' && dayOfWeek === '*') {
    scheduleType.value = 'monthly';
    scheduleConfig.value.dayOfMonth = parseInt(dayOfMonth) || 1;
    scheduleConfig.value.time = `${hour.padStart(2, '0')}:${minute.padStart(2, '0')}`;
  } else {
    scheduleType.value = 'custom';
    scheduleConfig.value.customCron = cron;
  }
};

// Generate cron from UI
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

// Watch for cron changes
watch(generatedCron, (newCron) => {
  formData.value.schedule_cron = newCron;
});

// Handle schedule toggle
function handleScheduleToggle() {
  if (!isScheduled.value) {
    formData.value.schedule_cron = '';
  }
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

// Get changed fields
const getChangedFields = () => {
  if (!originalData.value) return {};

  const changes = {};
  
  Object.keys(formData.value).forEach(key => {
    if (formData.value[key] !== originalData.value[key]) {
      changes[key] = formData.value[key];
    }
  });

  return changes;
};

// Computed changed fields
const changedFieldsCount = computed(() => {
  return Object.keys(getChangedFields()).length;
});

const changedFieldsList = computed(() => {
  return Object.keys(getChangedFields());
});

// Format field name
const formatFieldName = (fieldName) => {
  return fieldName.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase());
};

// Handle submit
const handleSubmit = async () => {
  const changes = getChangedFields();

  if (Object.keys(changes).length === 0) {
    submitError.value = 'No changes detected';
    return;
  }

  try {
    isSaving.value = true;
    submitError.value = null;

    await jobService.updateJob(props.jobId, changes);

    emit('success');
    emit('close');
    
  } catch (err) {
    submitError.value = err.response?.data?.error || err.message || 'Failed to update job';
    console.error('Update job error:', err);
  } finally {
    isSaving.value = false;
  }
};

// Handle close
const handleClose = () => {
  if (!isSaving.value) {
    emit('close');
  }
};

function formatGB(value) {
  if (!value) return "0.00";
  return Number(value).toFixed(2);
}

function calcPercentage(used, total) {
  if (!total || total === 0) return 0;
  return (used / total) * 100;
}

</script>

<style scoped>
.edit-job-view {
  max-width: 900px;
  margin: 0 auto;
}

/* Modal */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}

.modal-content {
  background: #fff;
  border-radius: 8px;
  border: 1px solid #e5e5e5;
  width: 90%;
  max-width: 800px;
  max-height: 90vh;
  overflow-y: auto;
  position: relative;
}

.modal-content::-webkit-scrollbar {
  width: 8px;
}

.modal-content::-webkit-scrollbar-track {
  background: #f0f0f0;
}

.modal-content::-webkit-scrollbar-thumb {
  background: #d4d4d4;
  border-radius: 4px;
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

/* Loading & Error */
.loading-container {
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
  border: 3px solid #e5e5e5;
  border-top-color: #1a1a1a;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.error-banner {
  background: #fee2e2;
  border: 1px solid #fecaca;
  color: #991b1b;
  padding: 1rem;
  margin: 1rem;
  border-radius: 6px;
}

.error-banner p {
  margin: 0 0 0.5rem 0;
  font-weight: 500;
}

.retry-btn {
  background: #991b1b;
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.875rem;
  font-weight: 500;
}

.retry-btn:hover {
  background: #7f1d1d;
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

.form-group input[type="text"],
.form-group input[type="number"],
.form-group input[type="time"],
.form-group select,
.form-group textarea {
  width: 100%;
  padding: 0.625rem 0.875rem;
  border: 1px solid #e5e5e5;
  border-radius: 6px;
  font-size: 0.9375rem;
  transition: all 0.2s;
  box-sizing: border-box;
}

.form-group input:focus,
.form-group select:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #1a1a1a;
}

.form-group textarea {
  font-family: 'Consolas', 'Monaco', monospace;
  font-size: 0.875rem;
  line-height: 1.6;
  resize: vertical;
}

.hint {
  display: block;
  margin-top: 0.375rem;
  font-size: 0.8125rem;
  color: #666;
}

.hint a {
  color: #1a1a1a;
  text-decoration: underline;
}

/* Info Display (Read-only fields) */
.info-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
  color: #1a1a1a;
  font-size: 0.875rem;
}

.info-display {
  background: #f5f5f5;
  padding: 0.625rem 0.875rem;
  border-radius: 6px;
  border: 1px solid #e5e5e5;
  color: #666;
  font-size: 0.9375rem;
}

.badge {
  display: inline-block;
  padding: 0.375rem 0.75rem;
  border-radius: 4px;
  font-size: 0.8125rem;
  font-weight: 600;
  text-transform: uppercase;
}

.badge.backup {
  background: #dbeafe;
  color: #1e40af;
}

.badge.restore {
  background: #dcfce7;
  color: #166534;
}

/* Schedule Section */
.schedule-section {
  background: #fafafa;
  padding: 1rem;
  border-radius: 6px;
  margin-bottom: 1rem;
  border: 1px solid #e5e5e5;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.25rem;
}

.section-header h3 {
  margin: 0;
  font-size: 0.9375rem;
  font-weight: 600;
  color: #1a1a1a;
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
  width: 44px;
  height: 24px;
  appearance: none;
  background: #d4d4d4;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.toggle-switch input[type="checkbox"]:checked {
  background: #1a1a1a;
}

.toggle-switch input[type="checkbox"]::before {
  content: '';
  position: absolute;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: white;
  top: 2px;
  left: 2px;
  transition: all 0.2s;
}

.toggle-switch input[type="checkbox"]:checked::before {
  left: 22px;
}

.toggle-label {
  font-weight: 500;
  color: #1a1a1a;
  font-size: 0.875rem;
}

/* Schedule Options */
.schedule-options {
  margin-top: 1.25rem;
}

.schedule-type-selector {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(80px, 1fr));
  gap: 0.5rem;
  margin-bottom: 1.25rem;
}

.type-btn {
  padding: 0.625rem;
  background: white;
  border: 1px solid #e5e5e5;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
  font-weight: 500;
  font-size: 0.875rem;
  color: #666;
}

.type-btn:hover {
  border-color: #1a1a1a;
  color: #1a1a1a;
}

.type-btn.active {
  background: #1a1a1a;
  border-color: #1a1a1a;
  color: white;
}

/* Schedule Config */
.schedule-config {
  background: white;
  padding: 1.25rem;
  border-radius: 6px;
  margin-bottom: 1rem;
  border: 1px solid #e5e5e5;
}

.schedule-config label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
  color: #1a1a1a;
  font-size: 0.875rem;
}

.input-group {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-bottom: 1rem;
}

.time-input {
  flex: 0 0 auto;
  padding: 0.625rem 0.875rem;
  border: 1px solid #e5e5e5;
  border-radius: 6px;
  font-size: 0.9375rem;
  font-weight: 500;
  min-width: 100px;
}

.time-input:focus {
  outline: none;
  border-color: #1a1a1a;
}

.input-suffix {
  color: #666;
  font-weight: 500;
  font-size: 0.875rem;
}

/* Weekdays Selector */
.weekdays-selector {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 0.5rem;
  margin-bottom: 1rem;
}

.weekday-btn {
  padding: 0.625rem 0.375rem;
  background: white;
  border: 1px solid #e5e5e5;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 500;
  font-size: 0.875rem;
  color: #666;
  transition: all 0.2s;
}

.weekday-btn:hover {
  border-color: #1a1a1a;
  color: #1a1a1a;
}

.weekday-btn.active {
  background: #1a1a1a;
  border-color: #1a1a1a;
  color: white;
}

/* Cron Input */
.cron-input {
  width: 100%;
  padding: 0.625rem 0.875rem;
  border: 1px solid #e5e5e5;
  border-radius: 6px;
  font-family: 'Consolas', monospace;
  font-size: 0.875rem;
}

.cron-input:focus {
  outline: none;
  border-color: #1a1a1a;
}

/* Cron Preview */
.cron-preview {
  background: #1a1a1a;
  padding: 1rem;
  border-radius: 6px;
  display: flex;
  align-items: center;
  gap: 0.75rem;
  flex-wrap: wrap;
}

.preview-label {
  color: #999;
  font-weight: 500;
  font-size: 0.8125rem;
}

.preview-code {
  background: #2a2a2a;
  color: #22c55e;
  padding: 0.375rem 0.75rem;
  border-radius: 4px;
  font-family: 'Consolas', monospace;
  font-size: 0.875rem;
  font-weight: 500;
}

.preview-description {
  color: #d4d4d4;
  font-size: 0.875rem;
}

/* Changes Summary */
.changes-summary {
  background: #fef3c7;
  border: 1px solid #fcd34d;
  padding: 1rem;
  border-radius: 6px;
  margin-bottom: 1rem;
}

.summary-title {
  margin: 0 0 0.75rem 0;
  font-weight: 600;
  color: #92400e;
  font-size: 0.875rem;
}

.changes-list {
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.changes-list li {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.8125rem;
  color: #92400e;
}

.field-name {
  font-weight: 600;
  min-width: 120px;
}

.field-change {
  font-family: 'Consolas', monospace;
  font-size: 0.75rem;
  background: rgba(0, 0, 0, 0.05);
  padding: 0.25rem 0.5rem;
  border-radius: 3px;
}

/* Form Actions */
.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  margin-top: 1.5rem;
  padding-top: 1.25rem;
  border-top: 1px solid #e5e5e5;
}

.btn-secondary {
  background: white;
  color: #666;
  border: 1px solid #e5e5e5;
  padding: 0.625rem 1.125rem;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 500;
  font-size: 0.9375rem;
  transition: all 0.2s;
}

.btn-secondary:hover:not(:disabled) {
  border-color: #1a1a1a;
  color: #1a1a1a;
  background: #f5f5f5;
}

.btn-secondary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
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
  
  .schedule-type-selector {
    grid-template-columns: repeat(auto-fit, minmax(70px, 1fr));
  }
  
  .weekdays-selector {
    gap: 0.375rem;
  }
  
  .weekday-btn {
    padding: 0.5rem 0.25rem;
    font-size: 0.8125rem;
  }
  
  .form-actions {
    flex-direction: column;
  }
  
  .btn-submit,
  .btn-secondary {
    width: 100%;
  }

  .changes-list li {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.25rem;
  }
}
</style>