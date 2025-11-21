<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="isVisible && log" class="modal-overlay" @click.self="handleClose">
        <div class="modal-container">
          <!-- Header -->
          <div class="modal-header">
            <div class="header-left">
              <h3>Log Details</h3>
              <span class="log-id">ID: {{ log.ID }}</span>
            </div>
            <button @click="handleClose" class="close-btn">×</button>
          </div>

          <!-- Content -->
          <div class="modal-body">
            <!-- Job Info -->
            <div class="info-section">
              <h4>Job Information</h4>
              <div class="info-grid">
                <div class="info-item">
                  <span class="label">Object</span>
                  <span class="value">{{ jobObject.sourcePath || 'Manual Job' }}</span>
                </div>
                <div class="info-item">
                  <span class="label">Job Name</span>
                  <span class="value">{{ getJobName(log) }}</span>
                </div>
                <div class="info-item">
                  <span class="label">Status</span>
                  <span class="value">
                    <span class="status-badge" :class="getStatusClass(log.Status)">
                      {{ log.Status }}
                    </span>
                  </span>
                </div>
                <div class="info-item">
                  <span class="label">Duration</span>
                  <span class="value">{{ log.DurationSec }} seconds</span>
                </div>
                <div class="info-item">
                  <span class="label">Transferred Byte</span>
                  <span class="value">{{ formatFileSize(log.TransferredByte) }}</span>
                </div>
                <div class="info-item">
                  <span class="label">Timestamp</span>
                  <span class="value">{{ formatFullTimestamp(log.Timestamp) }}</span>
                </div>
              </div>
            </div>

            <!-- Output Message -->
            <div class="info-section">
              <h4>Output Message</h4>
              <div class="output-box">
                <pre>{{ log.Message || 'No message available' }}</pre>
              </div>
            </div>

            <!-- Config Snapshot (if manual job) -->
            <div v-if="log.ConfigSnapshot" class="info-section">
              <h4>Configuration Snapshot</h4>
              <div class="config-box">
                <pre>{{ formatJSON(log.ConfigSnapshot) }}</pre>
              </div>
            </div>
          </div>

          <!-- Footer -->
          <div class="modal-footer">
            <button @click="handleClose" class="btn-close">Close</button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  isVisible: Boolean,
  log: Object
});

const emit = defineEmits(['close']);

// Computed property untuk JobObject yang berisi sourcePath dan transferredByte
const jobObject = computed(() => {
  if (!props.log) return {};
  
  return {
    sourcePath: props.log.ScheduledJob?.SourcePath || extractSourcePath(),
    transferredByte: props.log.TransferredByte || 0
  };
});

function extractSourcePath() {
  if (!props.log?.ConfigSnapshot) return null;
  try {
    const config = JSON.parse(props.log.ConfigSnapshot);
    return config.source_path || config.sourcePath || null;
  } catch (e) {
    return null;
  }
}

function handleClose() {
  emit('close');
}

function getJobName(log) {
  if (log.ScheduledJob?.JobName) return log.ScheduledJob.JobName;
  if (log.ConfigSnapshot) {
    try {
      const config = JSON.parse(log.ConfigSnapshot);
      return config.job_name || 'Manual Job';
    } catch (e) {
      return 'Manual Job';
    }
  }
  return 'Unknown Job';
}

function getStatusClass(status) {
  const s = status.toUpperCase();
  if (['SUCCESS', 'COMPLETED'].includes(s)) return 'status-success';
  if (s.includes('FAIL') || s === 'ERROR') return 'status-failed';
  if (s === 'RUNNING') return 'status-running';
  return 'status-default';
}

function formatFullTimestamp(timestamp) {
  try {
    return new Date(timestamp).toLocaleString('id-ID');
  } catch (e) {
    return timestamp;
  }
}

function formatFileSize(bytes) {
  if (!bytes || bytes === 0) return '0 B';
  
  const units = ['B', 'KB', 'MB', 'GB', 'TB'];
  const k = 1024;
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + units[i];
}

function formatJSON(jsonString) {
  try {
    return JSON.stringify(JSON.parse(jsonString), null, 2);
  } catch (e) {
    return jsonString;
  }
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-container {
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.2);
  max-width: 800px;
  width: 90%;
  max-height: 85vh;
  overflow-y: auto;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem;
  border-bottom: 1px solid #e5e5e5;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.header-left h3 {
  margin: 0;
  font-size: 1.25rem;
  color: #1a1a1a;
}

.log-id {
  font-size: 0.85rem;
  color: #999;
  font-weight: 500;
}

.close-btn {
  background: none;
  border: none;
  font-size: 1.5rem;
  color: #999;
  cursor: pointer;
  transition: color 0.2s;
}

.close-btn:hover {
  color: #1a1a1a;
}

.modal-body {
  padding: 1.5rem;
}

.info-section {
  margin-bottom: 2rem;
}

.info-section h4 {
  margin: 0 0 1rem 0;
  font-size: 1rem;
  color: #1a1a1a;
  font-weight: 600;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1rem;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.label {
  font-size: 0.75rem;
  color: #999;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.value {
  font-size: 0.9rem;
  color: #1a1a1a;
  font-weight: 500;
  word-break: break-word;
}

.status-badge {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  border-radius: 4px;
  font-size: 0.8rem;
  font-weight: 600;
  width: fit-content;
}

.status-success {
  background-color: #d4edda;
  color: #155724;
}

.status-failed {
  background-color: #f8d7da;
  color: #721c24;
}

.status-running {
  background-color: #d1ecf1;
  color: #0c5460;
}

.status-default {
  background-color: #e2e3e5;
  color: #383d41;
}

.output-box,
.config-box {
  background-color: #f8f9fa;
  border: 1px solid #e5e5e5;
  border-radius: 6px;
  padding: 1rem;
  overflow-x: auto;
}

.output-box pre,
.config-box pre {
  margin: 0;
  font-family: 'Courier New', monospace;
  font-size: 0.85rem;
  color: #1a1a1a;
  line-height: 1.5;
  white-space: pre-wrap;
  word-wrap: break-word;
}

.modal-footer {
  padding: 1rem 1.5rem;
  border-top: 1px solid #e5e5e5;
  display: flex;
  justify-content: flex-end;
}

.btn-close {
  padding: 0.5rem 1rem;
  background-color: #f5f5f5;
  border: 1px solid #e5e5e5;
  border-radius: 4px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.2s;
}

.btn-close:hover {
  background-color: #e5e5e5;
  border-color: #999;
}

/* Modal animations */
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

@media (max-width: 768px) {
  .modal-container {
    width: 95%;
  }

  .info-grid {
    grid-template-columns: 1fr;
  }

  .modal-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.5rem;
  }

  .close-btn {
    align-self: flex-end;
  }
}
</style>