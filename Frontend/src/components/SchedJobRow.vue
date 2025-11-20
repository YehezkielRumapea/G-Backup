<template>
  <tr>
    <td><strong>{{ job.job_name }}</strong></td>
    
    <td>{{ job.type }}</td>

    <td>{{ job.gdrive_target }}</td>

    <td>{{ job.last_run || 'N/A' }}</td>

    <td>
      <span class="status" :class="job.status.toLowerCase()">
        {{ job.status }}
      </span>
    </td>
    
    <td>{{ job.next_run || 'N/A' }}</td>
    
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

        <!-- View Script Button -->
        <button 
          @click="$emit('view-script', job.id)" 
          class="action-btn view" 
          title="View Script"
        >
          (i)
        </button>

        <!-- Delete Button -->
        <button 
          @click="confirmDelete" 
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
const props = defineProps({
  job: {
    type: Object,
    required: true
  }
});

// ✅ Tambah 'edit' ke emits
const emit = defineEmits(['trigger', 'view-script', 'edit', 'delete']); 

function confirmDelete() {
  if (confirm(`PERINGATAN! Menghapus Job Terjadwal '${props.job.job_name}' juga akan menghapus jadwal CRON.\nApakah Anda yakin ingin menghapus job ini?`)) {
    emit('delete', props.job.id, props.job.job_name);
  }
}
</script>

<style scoped>
/* --- STYLING BARIS JOB --- */
.status {
  padding: 4px 8px;
  border-radius: 12px;
  font-weight: bold;
  font-size: 0.8rem;
  color: white;
}

.status.pending { 
  background-color: #aaa; 
}

.status.running { 
  background-color: #3498db;
  font-size: 12px;
}

.status.completed, 
.status.success { 
  background-color: #2ecc71; 
}

.status.failed, 
.status.fail_source_check, 
.status.fail_pre_script, 
.status.fail_rclone,
.status.fail_post_script { 
  background-color: #e74c3c; 
}

/* Styling kolom aksi */
.job-actions-col {
  width: 150px;
}

.actions {
  display: flex;
  gap: 8px;
  justify-content: center;
}

.action-btn {
  border: none;
  background: none;
  cursor: pointer;
  font-size: 1.2rem;
  padding: 5px;
  transition: opacity 0.2s;
}

.action-btn:hover {
  opacity: 0.8;
}

/* Warna tombol aksi */
.action-btn.edit { 
  color: #f39c12; 
}

.action-btn.view { 
  color: #3499db; 
  font-weight: bold; 
}

.action-btn.delete { 
  color: #e74c3c; 
}

.action-btn.delete.disabled {
  color: #ccc;
  cursor: not-allowed;
}
</style>