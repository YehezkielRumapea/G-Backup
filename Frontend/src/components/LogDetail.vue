<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="isVisible && log" class="modal-overlay" @click.self="handleClose">
        <div class="modal-container">
          <!-- Header -->
          <div class="modal-header">
            <div class="header-left">
              <span class="icon">📝</span>
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
                  <span class="label">Job ID:</span>
                  <span class="value">{{ log.JobID || 'Manual Job' }}</span>
                </div>
                <div class="info-item">
                  <span class="label">Job Name:</span>
                  <span class="value">{{ getJobName(log) }}</span>
                </div>
                <div class="info-item">
                  <span class="label">Operation:</span>
                  <span class="value operation">{{ log.OperationType }}</span>
                </div>
                <div class="info-item">
                  <span class="label">Status:</span>
                  <span class="value">
                    <span class="status-badge" :class="getStatusClass(log.Status)">
                      {{ getStatusIcon(log.Status) }} {{ log.Status }}
                    </span>
                  </span>
                </div>
                <div class="info-item">
                  <span class="label">Duration:</span>
                  <span class="value">{{ log.DurationSec }} seconds</span>
                </div>
                <div class="info-item">
                  <span class="label">Timestamp:</span>
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
const props = defineProps({
  isVisible: Boolean,
  log: Object
});

const emit = defineEmits(['close']);

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

function getStatusIcon(status) {
  const s = status.toUpperCase();
  if (['SUCCESS', 'COMPLETED'].includes(s)) return '✅';
  if (s.includes('FAIL') || s === 'ERROR') return '❌';
  if (s === 'RUNNING') return '⏳';
  return '📝';
}

function formatFullTimestamp(timestamp) {
  try {
    return new Date(timestamp).toLocaleString('id-ID');
  } catch (e) {
    return timestamp;
  }
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
/* Modal styles - similar to ScriptPreview.vue */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.75);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  padding: 20px;
}

.modal-container {
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  width: 90%;
  max-width: 900px;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  animation: slideIn 0.3s ease-out;
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(-30px) scale(0.95);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

/* Header */
.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem 2rem;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-bottom: 2px solid rgba(255, 255, 255, 0.1);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.icon {
  font-size: 1.8rem;
}

.modal-header h3 {
  margin: 0;
  font-size: 1.5rem;
  font-weight: 600;
}

.log-id {
  background: rgba(255, 255, 255, 0.2);
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 0.85rem;
  font-weight: 500;
}

.close-btn {
  background: rgba(255, 255, 255, 0.15);
  border: none;
  color: white;
  width: 40px;
  height: 40px;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
  font-size: 2rem;
  line-height: 1;
}

.close-btn:hover {
  background: rgba(255, 255, 255, 0.25);
  transform: scale(1.1);
}

/* Body */
.modal-body {
  flex: 1;
  overflow-y: auto;
  padding: 2rem;
  background: #f8f9fa;
}

.modal-body::-webkit-scrollbar {
  width: 8px;
}

.modal-body::-webkit-scrollbar-track {
  background: #e9ecef;
}

.modal-body::-webkit-scrollbar-thumb {
  background: #cbd5e0;
  border-radius: 4px;
}

.modal-body::-webkit-scrollbar-thumb:hover {
  background: #a0aec0;
}

/* Info Sections */
.info-section {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  margin-bottom: 1.5rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.info-section:last-child {
  margin-bottom: 0;
}

.info-section h4 {
  margin: 0 0 1rem 0;
  font-size: 1.1rem;
  font-weight: 600;
  color: #2c3e50;
  padding-bottom: 0.75rem;
  border-bottom: 2px solid #e9ecef;
}

/* Info Grid */
.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.25rem;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.info-item .label {
  font-size: 0.85rem;
  color: #6c757d;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.info-item .value {
  font-size: 1rem;
  color: #2c3e50;
  font-weight: 600;
}

.info-item .value.operation {
  color: #667eea;
  text-transform: uppercase;
}

/* Status Badge */
.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 6px 14px;
  border-radius: 20px;
  font-weight: 600;
  font-size: 0.85rem;
  text-transform: uppercase;
}

.status-success {
  background: #d4edda;
  color: #155724;
}

.status-failed {
  background: #f8d7da;
  color: #721c24;
}

.status-running {
  background: #fff3cd;
  color: #856404;
}

.status-default {
  background: #e2e3e5;
  color: #383d41;
}

/* Output Box */
.output-box,
.config-box {
  background: #1e1e1e;
  border-radius: 8px;
  padding: 1.5rem;
  overflow-x: auto;
}

.output-box pre,
.config-box pre {
  margin: 0;
  color: #d4d4d4;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 0.9rem;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-word;
}

/* Footer */
.modal-footer {
  padding: 1.5rem 2rem;
  background: white;
  border-top: 1px solid #e9ecef;
  display: flex;
  justify-content: flex-end;
}

.btn-close {
  background: #667eea;
  color: white;
  border: none;
  padding: 0.75rem 2rem;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 600;
  font-size: 1rem;
  transition: all 0.2s;
}

.btn-close:hover {
  background: #5568d3;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

/* Transitions */
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-active .modal-container,
.modal-leave-active .modal-container {
  transition: transform 0.3s ease;
}

.modal-enter-from .modal-container,
.modal-leave-to .modal-container {
  transform: scale(0.9) translateY(-20px);
}

/* Responsive */
@media (max-width: 768px) {
  .modal-container {
    width: 95%;
    max-height: 85vh;
  }
  
  .modal-header {
    padding: 1rem 1.25rem;
  }
  
  .header-left {
    flex-wrap: wrap;
    gap: 0.75rem;
  }
  
  .modal-header h3 {
    font-size: 1.2rem;
  }
  
  .modal-body {
    padding: 1.25rem;
  }
  
  .info-grid {
    grid-template-columns: 1fr;
  }
  
  .info-section {
    padding: 1.25rem;
  }
}
</style>