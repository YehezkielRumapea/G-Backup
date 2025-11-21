<template>
    <tr> <div class="job-name">
        <td><strong>{{ job.job_name }}</strong></td>
        </div>
        
        <td>{{ job.source_path }}</td> 

        <td>{{ job.gdrive_target }}</td>

        <td>{{ job.last_run || 'N/A' }}</td>

        <td>
            <span class="status" :class="getStatusClass(job.status)">
                {{ job.status }}
            </span>
        </td>
        
        <td class="job-actions-col"> 
          <div class="actions">
            <!-- Edit -->
            <button 
              @click="$emit('edit', job.id)" 
              class="icon-btn edit" 
              title="Edit Job"
            >
              ✏️
            </button>

            <!-- Run -->
            <button 
              @click="$emit('trigger', job.id)" 
              class="icon-btn run" 
              title="Run Job"
            >
              ▶
            </button>

            <!-- View Script -->
            <button 
              @click="$emit('view-script', job.id)" 
              class="icon-btn view" 
              title="View Script"
            >
              (i)
            </button>

            <!-- Delete -->
            <button 
              @click="handleDelete" 
              class="icon-btn delete" 
              title="Delete Job"
            >
              🗑
            </button>

          </div>
        </td>
    </tr>
</template>

<script setup>
// ✅ Define props
const props = defineProps({
  job: {
    type: Object,
    required: true
  }
});

// ✅ Define emits (tambahkan 'edit')
const emit = defineEmits(['trigger', 'view-script', 'edit', 'delete']);

// ✅ Handle delete dengan emit ke parent
function handleDelete() {
  if (!confirm(`Apakah Anda yakin ingin menghapus job manual "${props.job.job_name}"?\nAksi ini tidak dapat dibatalkan.`)) {
    return;
  }
  
  // Emit event delete ke parent dengan jobId dan jobName
  emit('delete', props.job.id, props.job.job_name);
}

// ✅ Get status CSS class
function getStatusClass(status) {
  if (!status) return 'pending';
  const statusLower = status.toLowerCase();
  
  // Map status to CSS classes
  if (statusLower === 'completed' || statusLower === 'success') return 'completed';
  if (statusLower === 'running') return 'running';
  if (statusLower === 'pending') return 'pending';
  if (statusLower.includes('fail')) return 'failed';
  
  return statusLower;
}
</script>

<style scoped>/* =============================== */
/* == TABLE ROW HOVER EFFECT ==== */
/* =============================== */
tr:hover {
  background-color: #f8f9fa;
}

/* =============================== */
/* == STATUS BADGES (SAMA DGN REMOTE) == */
/* =============================== */
.status {
  display: inline-block;
  padding: 3px 8px;
  border-radius: 4px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  color: white;
}

.status.pending { 
  background-color: #7f8c8d; 
}

.status.running { 
  background-color: #3498db;
}

.status.completed,
.status.success { 
  background-color: #27ae60; 
}

.status.failed,
.status.fail_pre_script,
.status.fail_rclone,
.status.fail_post_script { 
  background-color: #e74c3c; 
}

/* =============================== */
/* == TABLE CELL STYLE ========== */
/* =============================== */
td {
  padding: 0.75rem;
  vertical-align: middle;
}

strong {
  color: #1a1a1a;
}

/* =============================== */
/* == ACTION COLUMN ============= */
/* =============================== */
.job-actions-col {
  width: 150px;
}

.job-name {
  font-size: 0.90rem; /* atau 0.65rem kalau mau lebih kecil lagi */
  margin-bottom: 0.25rem;
}

.actions {
  display: flex;
  justify-content: center;
  gap: 10px;
  align-items: center;
}

/* =============================== */
/* == ICON BUTTONS (SAMA DGN REMOTE) == */
/* =============================== */
.icon-btn {
  border: none;
  background: none;
  cursor: pointer;
  padding: 6px;
  font-size: 1.3rem;
  transition: transform 0.15s ease, opacity 0.2s ease;
}

.icon-btn:hover {
  opacity: 0.85;
  transform: scale(1.15);
}

/* Warna ikon — konsisten dengan yang sebelumnya */
.icon-btn.edit { 
  color: #f39c12; 
}

.icon-btn.run { 
  color: #f59e0b; 
}

.icon-btn.view { 
  color: #3498db;
  font-weight: bold;
}

.icon-btn.delete { 
  color: #e74c3c; 
}

/* Disabled */
.icon-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
  transform: none !important;
}

/* =============================== */
/* RESPONSIVE ==================== */
/* =============================== */
@media (max-width: 768px) {
  .icon-btn {
    font-size: 1.15rem;
  }

  .actions {
    gap: 6px;
  }

  .job-actions-col {
    width: 130px;
  }
}

</style>