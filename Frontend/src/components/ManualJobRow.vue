<template>
    <tr>
        <td><strong>{{ job.job_name }}</strong></td>
        
        <td>{{ job.type }}</td> 

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

<style scoped>
/* Row hover effect */
/* =============================== */
/* == TABLE ROW HOVER EFFECT ==== */
/* =============================== */
tr:hover {
  background-color: #f8f9fa;
}

/* =============================== */
/* == STATUS BADGES ============= */
/* =============================== */
.status {
  padding: 4px 10px;
  border-radius: 12px;
  font-weight: 600;
  font-size: 0.8rem;
  color: white;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  display: inline-block;
}

.status.pending { 
  background-color: #94a3b8; 
}

.status.running { 
  background-color: #3b82f6;
  animation: pulse-blue 2s ease-in-out infinite;
}

.status.completed,
.status.success { 
  background-color: #10b981; 
}

.status.failed,
.status.fail,
.status.fail_pre_script,
.status.fail_rclone,
.status.fail_post_script { 
  background-color: #ef4444; 
}

@keyframes pulse-blue {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.7; }
}

/* =============================== */
/* == ACTIONS COLUMN ============ */
/* =============================== */
.job-actions-col {
  width: 160px;
}

.actions {
  display: flex;
  gap: 10px;
  justify-content: center;
  align-items: center;
}

/* =============================== */
/* == ICON-ONLY BUTTON STYLE ==== */
/* =============================== */
.icon-btn {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 1.4rem;
  padding: 4px;
  transition: transform 0.15s ease, color 0.15s ease;
  color: #6b7280; /* default gray */
}

.icon-btn:hover {
  transform: scale(1.2);
  color: #000;
}

.icon-btn:active {
  transform: scale(1.05);
}

/* =============================== */
/* == SPECIFIC ICON COLORS ====== */
/* =============================== */
.icon-btn.edit { 
  color: #3b82f6; 
}

.icon-btn.run { 
  color: #f59e0b; 
}

.icon-btn.view { 
  color: #3499db; 
  font-weight: bold;
}

.icon-btn.delete { 
  color: #ef4444; 
}

.icon-btn.edit:hover { 
  color: #1d4ed8; 
}

.icon-btn.run:hover { 
  color: #d97706; 
}

.icon-btn.view:hover { 
  color: #6d28d9; 
}

.icon-btn.delete:hover { 
  color: #b91c1c; 
}

/* =============================== */
/* == DISABLED STATE ============ */
/* =============================== */
.icon-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
  transform: none !important;
}

/* =============================== */
/* == RESPONSIVE ================ */
/* =============================== */
@media (max-width: 768px) {
  .icon-btn {
    font-size: 1.2rem;
  }
  
  .actions {
    gap: 6px;
  }
  
  .job-actions-col {
    width: 130px;
  }
}


</style>