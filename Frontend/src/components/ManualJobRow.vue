<template>
    <tr>
        <td><strong>{{ job.id }}</strong></td> 
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
            <!-- ✅ Edit Button -->
            <button 
              @click="$emit('edit', job.id)" 
              class="action-btn edit" 
              title="Edit Job"
            >
              ✏️
            </button>

            <!-- Run Button -->
            <button 
              @click="$emit('trigger', job.id)" 
              class="action-btn play" 
              title="Run Now"
            >
              ▶
            </button>

            <!-- View Script Button -->
            <button 
              @click="$emit('view-script', job.id)" 
              class="action-btn view" 
              title="View Script"
            >
              👁
            </button>

            <!-- Delete Button -->
            <button 
              @click="handleDelete" 
              class="action-btn delete" 
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
tr:hover {
  background-color: #f8f9fa;
}

/* Styling status bar */
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

/* Pulse animation for running status */
@keyframes pulse-blue {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.7; }
}

/* Styling kolom aksi */
.job-actions-col {
  width: 180px; /* ✅ Tambah lebar untuk 4 tombol */
}

.actions {
  display: flex;
  gap: 6px;
  justify-content: center;
  align-items: center;
}

.action-btn {
  border: none;
  background: rgba(255, 255, 255, 0.9);
  cursor: pointer;
  font-size: 1.1rem;
  padding: 8px 10px;
  border-radius: 6px;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  min-width: 36px;
  min-height: 36px;
}

.action-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
}

.action-btn:active {
  transform: translateY(0);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

/* ✅ Edit button (blue) */
.action-btn.edit {
  background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
}

.action-btn.edit:hover {
  background: linear-gradient(135deg, #2563eb 0%, #1d4ed8 100%);
}

/* Play button (orange/yellow) */
.action-btn.play {
  background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
}

.action-btn.play:hover {
  background: linear-gradient(135deg, #d97706 0%, #b45309 100%);
}

/* View button (purple) */
.action-btn.view {
  background: linear-gradient(135deg, #8b5cf6 0%, #7c3aed 100%);
}

.action-btn.view:hover {
  background: linear-gradient(135deg, #7c3aed 0%, #6d28d9 100%);
}

/* Delete button (red) */
.action-btn.delete {
  background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
}

.action-btn.delete:hover {
  background: linear-gradient(135deg, #dc2626 0%, #b91c1c 100%);
}

/* Disabled state */
.action-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none !important;
}

/* Responsive */
@media (max-width: 768px) {
  .action-btn {
    font-size: 1rem;
    padding: 6px 8px;
    min-width: 32px;
    min-height: 32px;
  }
  
  .actions {
    gap: 4px;
  }
  
  .job-actions-col {
    width: 160px;
  }
}
</style>