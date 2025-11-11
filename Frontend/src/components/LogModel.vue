<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="isVisible && log" class="modal-overlay" @click.self="close">
        <div class="modal-content">
          <!-- Header -->
          <div class="modal-header">
            <h3>📝 Log Detail</h3>
            <button @click="close" class="close-btn">×</button>
          </div>

          <!-- Body -->
          <div class="modal-body">
            <div class="detail-row">
              <span class="label">Job Name:</span>
              <span class="value">{{ getJobName(log) }}</span>
            </div>

            <div class="detail-row">
              <span class="label">Job ID:</span>
              <span class="value">{{ log.JobID || 'Manual' }}</span>
            </div>

            <div class="detail-row">
              <span class="label">Operation:</span>
              <span class="value">{{ log.OperationType || 'N/A' }}</span>
            </div>

            <div class="detail-row">
              <span class="label">Status:</span>
              <span class="value">
                <span class="status-badge" :class="getStatusClass(log.Status)">
                  {{ getStatusIcon(log.Status) }} {{ log.Status }}
                </span>
              </span>
            </div>

            <div class="detail-row">
              <span class="label">Duration:</span>
              <span class="value">{{ log.DurationSec }} seconds</span>
            </div>

            <div class="detail-row">
              <span class="label">Timestamp:</span>
              <span class="value">{{ formatFullDate(log.Timestamp) }}</span>
            </div>

            <div class="detail-section">
              <span class="label">Message:</span>
              <pre class="message-box">{{ log.Message || 'No message' }}</pre>
            </div>
          </div>

          <!-- Footer -->
          <div class="modal-footer">
            <button @click="close" class="btn-close">Close</button>
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

function close() {
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
  return 'Unknown';
}

function getStatusClass(status) {
  const s = status.toUpperCase();
  if (['SUCCESS', 'COMPLETED'].includes(s)) return 'success';
  if (s.includes('FAIL') || s === 'ERROR') return 'failed';
  if (s === 'RUNNING') return 'running';
  return 'pending';
}

function getStatusIcon(status) {
  const s = status.toUpperCase();
  if (['SUCCESS', 'COMPLETED'].includes(s)) return '✅';
  if (s.includes('FAIL') || s === 'ERROR') return '❌';
  if (s === 'RUNNING') return '⏳';
  return '⏱️';
}

function formatFullDate(timestamp) {
  try {
    return new Date(timestamp).toLocaleString('id-ID');
  } catch (e) {
    return timestamp;
  }
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  padding: 20px;
}

.modal-content {
  background: white;
  border-radius: 12px;
  width: 90%;
  max-width: 600px;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.modal-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 1.25rem 1.5rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-radius: 12px 12px 0 0;
}

.modal-header h3 {
  margin: 0;
  font-size: 1.25rem;
}

.close-btn {
  background: rgba(255, 255, 255, 0.2);
  border: none;
  color: white;
  width: 36px;
  height: 36px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 1.8rem;
  line-height: 1;
  transition: all 0.2s;
}

.close-btn:hover {
  background: rgba(255, 255, 255, 0.3);
}

.modal-body {
  flex: 1;
  overflow-y: auto;
  padding: 1.5rem;
}

.detail-row {
  display: flex;
  padding: 0.75rem 0;
  border-bottom: 1px solid #f0f0f0;
}

.detail-row:last-of-type {
  border-bottom: none;
}

.detail-section {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  padding: 0.75rem 0;
  border-top: 2px solid #e9ecef;
  margin-top: 0.75rem;
}

.label {
  font-weight: 600;
  color: #6c757d;
  min-width: 120px;
  font-size: 0.9rem;
}

.value {
  color: #2c3e50;
  font-size: 0.9rem;
}

.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 4px 10px;
  border-radius: 12px;
  font-weight: 600;
  font-size: 0.8rem;
}

.status-badge.success {
  background: #d4edda;
  color: #155724;
}

.status-badge.failed {
  background: #f8d7da;
  color: #721c24;
}

.status-badge.running {
  background: #fff3cd;
  color: #856404;
}

.status-badge.pending {
  background: #e2e3e5;
  color: #383d41;
}

.message-box {
  background: #1e1e1e;
  color: #d4d4d4;
  padding: 1rem;
  border-radius: 6px;
  font-family: 'Consolas', monospace;
  font-size: 0.85rem;
  line-height: 1.5;
  white-space: pre-wrap;
  word-break: break-word;
  margin: 0;
  max-height: 200px;
  overflow-y: auto;
}

.modal-footer {
  padding: 1rem 1.5rem;
  background: #f8f9fa;
  border-top: 1px solid #e9ecef;
  display: flex;
  justify-content: flex-end;
  border-radius: 0 0 12px 12px;
}

.btn-close {
  background: #667eea;
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 600;
  transition: all 0.2s;
}

.btn-close:hover {
  background: #5568d3;
}

.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.3s;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}
</style>