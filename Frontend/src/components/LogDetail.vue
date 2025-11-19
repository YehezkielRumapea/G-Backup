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
                  <span class="label">Job ID</span>
                  <span class="value">{{ log.JobID || 'Manual Job' }}</span>
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
                  <span class="label">TransferredByte</span>
                  <span class="value">{{ formatFullTimestamp(log.TransferredByte) }}</span>
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
