<template>
  <tr>
    <td><strong>{{ job.job_name }}</strong></td>
    <td>{{ job.type }}</td>
    <td>{{ job.gdrive_target }}</td>
    <td>{{ job.mode }}</td>
    <td>{{ job.last_run || 'N/A' }}</td>
    <td>
      <span class="status" :class="job.status.toLowerCase()">
        {{ job.status }}
      </span>
    </td>
    <td>{{ job.next_run || 'N/A' }}</td>
    
    <td>
      <div class="actions">
        <button @click="$emit('trigger', job.id)" class="action-btn play" title="Run Now">
          ▶
        </button>
        <button @click="$emit('view-script', job.id)" class="action-btn view" title="View Script">
          (i)
        </button>
        <button @click="handleDelete" class="action-btn delete" title="Delete">
          🗑
        </button>
      </div>
    </td>
  </tr>
</template>

<script setup>
// Terima data 'job' dari parent (ScheduledJobs.vue)
const props = defineProps({
  job: {
    type: Object,
    required: true
  }
})

// Definisikan event yang akan dikirim ke parent
const emit = defineEmits(['trigger', 'view-script'])

// (WIP) Logika untuk Hapus Job
function handleDelete() {
  if (confirm(`Hapus job "${props.job.job_name}"?`)) {
    // Panggil jobService.deleteJob(props.job.id)
    alert(`(WIP) Menghapus job ${props.job.id}`)
  }
}
</script>

<style scoped>
.status {
  padding: 4px 8px;
  border-radius: 12px;
  font-weight: bold;
  font-size: 0.8rem;
  color: white;
}
.status.pending { background-color: #aaa; }
.status.running { background-color: #3498db; }
.status.completed, .status.success { background-color: #2ecc71; }
.status.failed, .status.fail, .status.fail_pre_script, .status.fail_rclone { background-color: #e74c3c; }

.actions {
  display: flex;
  gap: 5px;
}
.action-btn {
  border: none;
  background: none;
  cursor: pointer;
  font-size: 1.2rem;
  padding: 5px;
}
.action-btn.play { color: #f39c12; }
.action-btn.delete { color: #e74c3c; }
.action-btn.view { color: #3498db; font-weight: bold; }
</style>